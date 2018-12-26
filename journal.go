/* Herald Journal-to-XMPP bridge
 * Copyright (C) 2018  Vianney le Clément de Saint-Marcq
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package main

import (
	"bytes"
	"github.com/coreos/go-systemd/sdjournal"
	"log"
	"regexp"
	"strconv"
	"text/template"
)

const TAIL_ENTRIES = 10

type criterion interface {
	Matches(e *sdjournal.JournalEntry) bool
}

// Exact match ("=VALUE")
type exactMatch struct {
	key   string
	value string
}

func (c *exactMatch) Matches(e *sdjournal.JournalEntry) bool {
	return e.Fields[c.key] == c.value
}

// Inverse match ("!VALUE")
type inverseMatch struct {
	key   string
	value string
}

func (c *inverseMatch) Matches(e *sdjournal.JournalEntry) bool {
	return e.Fields[c.key] != c.value
}

// Regular expression ("~PATTERN")
type regexpMatch struct {
	key     string
	pattern *regexp.Regexp
}

func (c *regexpMatch) Matches(e *sdjournal.JournalEntry) bool {
	return c.pattern.MatchString(e.Fields[c.key])
}

// Inequality ("<INTEGER")
type lessMatch struct {
	key   string
	value int
}

func (c *lessMatch) Matches(e *sdjournal.JournalEntry) bool {
	v, err := strconv.Atoi(e.Fields[c.key])
	return err == nil && v < c.value
}

// Inequality (">INTEGER")
type greaterMatch struct {
	key   string
	value int
}

func (c *greaterMatch) Matches(e *sdjournal.JournalEntry) bool {
	v, err := strconv.Atoi(e.Fields[c.key])
	return err == nil && v > c.value
}

// Set of criteria
type rule []criterion

func (r rule) Matches(e *sdjournal.JournalEntry) bool {
	for _, c := range r {
		if !c.Matches(e) {
			return false
		}
	}
	return true
}

// Parse the rules from the configuration
func parseRules() ([]rule, error) {
	rules := make([]rule, 0)
	for _, m := range Config.Match {
		r := make(rule, 0)
		for key, value := range m {
			switch value[0] {
			case '=':
				r = append(r, &exactMatch{key, value[1:]})
			case '!':
				r = append(r, &inverseMatch{key, value[1:]})
			case '~':
				re, err := regexp.Compile(value[1:])
				if err != nil {
					return nil, err
				}
				r = append(r, &regexpMatch{key, re})
			case '<':
				v, err := strconv.Atoi(value[1:])
				if err != nil {
					return nil, err
				}
				r = append(r, &lessMatch{key, v})
			case '>':
				v, err := strconv.Atoi(value[1:])
				if err != nil {
					return nil, err
				}
				r = append(r, &greaterMatch{key, v})
			default:
				r = append(r, &exactMatch{key, value})
			}
		}
		rules = append(rules, r)
	}
	return rules, nil
}

// ReadJournal reads and filters the systemd journal and sends messages of
// interest through the given channel.
func ReadJournal(messages chan<- string) {
	// Parse configuration
	rules, err := parseRules()
	if err != nil {
		log.Fatal(err)
	}
	tmpl, err := template.New("output").Parse(Config.Format)
	if err != nil {
		log.Fatal(err)
	}

	// Open journal
	journal, err := sdjournal.NewJournal()
	if err != nil {
		log.Fatal(err)
	}
	defer journal.Close()

	// Skip to the end of the journal (minus TAIL_ENTRIES entries)
	if err = journal.SeekTail(); err != nil {
		log.Fatal(err)
	}
	skip, err := journal.PreviousSkip(TAIL_ENTRIES + 1)
	if err != nil {
		log.Fatal(err)
	}
	if skip != TAIL_ENTRIES+1 {
		if err = journal.SeekHead(); err != nil {
			log.Fatal(err)
		}
	}

	// Read all entries
	for {
		for {
			n, err := journal.Next()
			if err != nil {
				log.Fatal(err)
			}
			if n == 0 { // we have reached the tail
				break
			}
			entry, err := journal.GetEntry()
			if err != nil {
				messages <- err.Error()
				continue
			}
			for _, r := range rules {
				if r.Matches(entry) {
					var buf bytes.Buffer
					err := tmpl.Execute(&buf, entry.Fields)
					if err == nil {
						messages <- buf.String()
					} else {
						messages <- err.Error()
					}
					break
				}
			}
		}
		journal.Wait(sdjournal.IndefiniteWait)
	}
}
