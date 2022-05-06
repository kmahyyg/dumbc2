# dumbyc2

Yet Another Simple C2(Command and Control) Tool

## DISCLAIMER

USE IT IN A LEGAL WAY. WE WON'T TAKE ANY RESPONSIBILITY FOR YOUR BEHAVIOUR.

## NOTE

We will not offer agent binary, since it needs your own CA cert to authenticaiton. 
You need to download Golang SDK and build it yourself.

Build Steps:

- Download `certgen` from release
- Run `certgen` and copy `client*.pem` `cacert.pem` and `serverpin.txt` from output directory, default is `~/.dumbyc2`
- Put the `client*.pem` `cacert.pem` and `serverpin.txt` to `buildtime/certs`
- Run: `make agent`

## Build

`make dumbc2` for Controller.

`make certgen` for Certificate Generator.

`make prune` for Clearing all output binaries.

## Current Status

v1.1.0-git

## Feature

- Agent Generator (Reverse Supported)
- Shell Command
- File Transfer Operation
- ALL Traffic Encrypted using TLS with Force SSL Pinning

## License

 dumbYC2
 Copyright (C) 2020  kmahyyg
 
 This program is free software: you can redistribute it and/or modify
 it under the terms of the GNU Affero General Public License as published by
 the Free Software Foundation, either version 3 of the License, or
 (at your option) any later version.
 
 This program is distributed in the hope that it will be useful,
 but WITHOUT ANY WARRANTY; without even the implied warranty of
 MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 GNU Affero General Public License for more details.
 
 You should have received a copy of the GNU Affero General Public License
 along with this program.  If not, see <http://www.gnu.org/licenses/>.

## Acknowledgement

- https://github.com/tiagorlampert/CHAOS
- https://github.com/brimstone/go-shellcode
- https://github.com/lesnuages/hershell
- https://github.com/thesecondsun/gosh
