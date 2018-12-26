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
	"flag"
	"github.com/BurntSushi/toml"
	"log"
	"os"
)

var Config struct {
	Server     string
	Username   string
	Password   string
	Security   string
	Recipients []string
	Format     string
	Match      []map[string]string
}

var configFile = flag.String("config", "", "configuration file")

var defaultConfigFiles = []string{
	"config.toml",
	"/etc/herald/config.toml",
	"/etc/herald.toml",
}

func LoadConfig() {
	// Set default values
	Config.Security = "starttls"
	Config.Format = "{{._COMM}}: {{.MESSAGE}}"
	// Locate configuration file
	if *configFile == "" {
		for _, filename := range defaultConfigFiles {
			if _, err := os.Stat(filename); err == nil {
				*configFile = filename
				break
			}
		}
		if *configFile == "" {
			log.Fatal("Configuration file not found")
			return
		}
	}
	// Load configuration
	if _, err := toml.DecodeFile(*configFile, &Config); err != nil {
		log.Print(*configFile, ": ", err)
	}
}
