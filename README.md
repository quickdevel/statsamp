## statSAMP ##
[![Build Status](https://travis-ci.org/quickdevel/statsamp.svg)](https://travis-ci.org/quickdevel/statsamp)

Service created for collection statistics San Andreas Multiplayer (sa-mp.com) servers. List of servers receive from the official master-server.


----------


## Getting started ##
 1. Go to the desired folder:
  `cd /your_path/`
 2. Set GOPAPTH:
 `export GOPATH=/your_path/statsamp`
 3. Download source code:
 `git clone https://github.com/quickdevel/statsamp.git`
 4. Download dependencies: 
 `go get -d ./...`
 5. Create folder for binaries: 
 `mkdir $GOPATH/bin && cd $GOPATH/bin`
 6. Build *server* and *updater* packages: 
 `go build statsamp/server statsamp/updater`
 7. Return to desired folder: 
 `cd /your_path`
 8. Import database structure from *mysql_structure.sql*.
 9. Configure database access in *statsamp.cfg*.
 10. Configure *crontab* for updating stats (use `crontab -e`).
Sample *(update every hour in 55 minutes)*:
`55 * * * * cd /your_path/statsamp && bin/updater`
 11. Start web-service:
`bin/server`


----------


## Master-servers versions ##

 - 0.3e - [0.3.4](http://lists.sa-mp.com/0.3.4/servers)
 - 0.3x - [0.3.5](http://lists.sa-mp.com/0.3.5/servers) / [0.3.5b](http://lists.sa-mp.com/0.3.5b/servers)
 - 0.3z - [0.3.6](http://lists.sa-mp.com/0.3.6/servers)
 - 0.3.7 - [0.3.7](http://lists.sa-mp.com/0.3.7/servers)

You can set required version in statsamp.cfg.
