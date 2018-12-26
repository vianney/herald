# Herald Journal-to-XMPP bridge
# Copyright (C) 2018  Vianney le Clément de Saint-Marcq
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program.  If not, see <https://www.gnu.org/licenses/>.

.PHONY: all clean install

DESTDIR ?=

all:
	go build -v

clean:
	-rm herald

install: all
	install -Dm755 herald $(DESTDIR)/usr/bin/herald
	install -Dm644 etc/config.toml $(DESTDIR)/etc/herald/config.toml
	install -Dm644 etc/herald.service $(DESTDIR)/usr/lib/systemd/system/herald.service
	install -Dm644 etc/sysusers.conf $(DESTDIR)/usr/lib/sysusers.d/herald.conf
	install -Dm644 etc/tmpfiles.conf $(DESTDIR)/usr/lib/tmpfiles.d/herald.conf
