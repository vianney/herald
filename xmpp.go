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
	"github.com/jpillora/backoff"
	"github.com/mattn/go-xmpp"
	"log"
	"time"
)

// SendOverXmpp sends all messages from the given channel to all recipients
// of the configuration over XMPP.
func SendOverXmpp(messages <-chan string) {
	// Parse configuration
	options := xmpp.Options{
		Host:     Config.Server,
		User:     Config.Username,
		Password: Config.Password,
	}
	switch Config.Security {
	case "none":
		options.NoTLS = true
		options.StartTLS = false
	case "tls":
		options.NoTLS = false
		options.StartTLS = false
	case "starttls":
		options.NoTLS = true
		options.StartTLS = true
	default:
		log.Fatal("Unknown security option: ", Config.Security)
	}

	// Connect to server with exponential backoff
	bo := &backoff.Backoff{
		Min:    100 * time.Millisecond,
		Max:    10 * time.Second,
		Factor: 2,
		Jitter: false,
	}
	for {
		done := make(chan struct{})
		talk, err := options.NewClient()
		if err != nil {
			log.Print(err)
			goto retry
		}
		defer talk.Close()
		log.Print("Connected as ", Config.Username)
		bo.Reset()

		// Receive handler: grant presence subscription to configured recipients
		go func() {
			for {
				chat, err := talk.Recv()
				if err != nil {
					log.Print("XMPP receive error: ", err)
					close(done)
					break
				}
				switch v := chat.(type) {
				case xmpp.Presence:
					if v.Type == "subscribe" {
						for _, recipient := range Config.Recipients {
							if v.From == recipient {
								log.Print("Granting presence subsciption to ", v.From)
								talk.ApproveSubscription(v.From)
								break
							}
						}
					}
				}
			}
		}()

		// Send all messages
		for {
			select {
			case msg := <-messages:
				for _, recipient := range Config.Recipients {
					_, err := talk.Send(xmpp.Chat{
						Remote: recipient,
						Type:   "chat",
						Text:   msg,
					})
					if err != nil {
						log.Print("XMPP send error: ", err)
						goto retry
					}
				}
			case <-done:
				goto retry
			}
		}

	retry:
		if talk != nil {
			talk.Close()
		}
		d := bo.Duration()
		log.Print("Reconnecting in ", d)
		time.Sleep(d)
	}
}
