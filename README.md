## statSAMP ##
Service created for collection statistics San Andreas Multiplayer (sa-mp.com) servers. List of servers receive from the official master-server (master.sa-mp.com).


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
