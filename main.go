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
	"log"
	"os"
)

func main() {
	flag.Parse()
	if os.Getenv("JOURNAL_STREAM") != "" {
		// stdout is connected to the systemd journal: remove timestamps
		// from log messages as they are taken care of by the journal
		log.SetFlags(0)
	}
	LoadConfig()
	msgs := make(chan string, 10)
	go SendOverXmpp(msgs)
	ReadJournal(msgs)
}
