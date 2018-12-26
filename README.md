Herald Journal-to-XMPP bridge
==============================

Herald is a simple daemon watching the systemd journal and sending selected
messages over XMPP.  It can be used as a simple server monitoring solution.


Quick start
------------

1. Install the herald package for your distribution (preferred), or install
   from source ([Go 1.11][golang] or later is required):
   ``` shell
   make
   sudo make install
   ```
2. Edit `/etc/herald/config.toml`.  You should at least set up your XMPP server
   and credentials, and the list of recipients.
3. `sudo systemctl enable --now herald.service`

[golang]: https://golang.org/dl/


Configuration
--------------

Herald will look for its configuration file in several locations in the
following order:

1. The path given as `-config PATH` argument in the command line,
2. `config.toml` in the current working directory,
3. `/etc/herald/config.toml`,
4. `/etc/herald.toml`.

See [`etc/config.toml`](etc/config.toml) for a commented example.


License
--------

Copyright (C) 2018  Vianney le Clément de Saint-Marcq

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.

This project makes use of various open-source libraries.
See [`go.mod`](go.mod) for the complete list of dependencies.
