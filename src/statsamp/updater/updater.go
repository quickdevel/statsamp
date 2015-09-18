package main

import (
  "fmt"
  "log"
  "sync"
  "time"
  // DB libs
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
  // StatSAMP libs
  "statsamp/config"
  "statsamp/query"
  "statsamp/serverlist"
)

const (
  maxQueryThreads = 25
  maxQueryAttempts = 3
  queryTimeout = 1
)

type HistoryResult struct {
  ServersTotal  int
  ServersOnline int
  SlotsTotal    int
  SlotsUsed     int
}

func main() {
  var historyResult HistoryResult
  var waitGroup sync.WaitGroup
  startTime := time.Now()
  // Load config
  cfg, err := config.GetConfigData()
  if err != nil {
    log.Fatal(err)
  }
  // Connect to DB
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.Database))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
  // Request server-list from lists.sa-mp.com
  log.Println("Download server-list from lists.sa-mp.com.")
  var serversAddr []string
  var currentAddr int
  request := serverlist.NewRequest(cfg.Updater.Version)
  if err := request.Exec(); err != nil {
    log.Fatal("Error (request.Exec()):", err)
  }
  serversAddr, err = request.ReadAll()
  if err != nil {
    log.Fatal("Error (request.ReadAll()):", err)
  }
  // Clear database table
  if err := clearServersInDB(db); err != nil {
    log.Fatal("Error (clearServersInDB()):", err)
  }
  // Create threads
  waitGroup.Add(maxQueryThreads)
  for i := 0; i < maxQueryThreads; i++ {
    go processLine(&serversAddr, &currentAddr, &historyResult, db, &waitGroup)
  }
  // Wait threads
  waitGroup.Wait()
  // Write history to database
  writeHistoryToDB(db, &historyResult)
  // Print results
  log.Println("Execution time:", time.Since(startTime))
  log.Println("Servers offline:", historyResult.ServersOnline, "of", historyResult.ServersTotal)
  log.Println("Slots used:", historyResult.SlotsUsed, "of", historyResult.SlotsTotal)
}

func writeHistoryToDB(db *sql.DB, result *HistoryResult) {
  statement, err := db.Prepare("INSERT INTO `history` (history_servers_total, history_servers_online, history_slots_total, history_slots_used, history_date) VALUES(?, ?, ?, ?, NOW())")
  if err != nil {
	   log.Fatal(err)
  }
  _, err = statement.Exec(result.ServersTotal, result.ServersOnline, result.SlotsTotal, result.SlotsUsed)
  if err != nil {
	   log.Fatal(err)
  }
}

func clearServersInDB(db *sql.DB) error {
  _, err := db.Exec("TRUNCATE `servers`")
  return err
}

func writeServerToDB(db *sql.DB, addr string, server *query.ServerInfo) error {
  statement, err := db.Prepare("INSERT INTO `servers` (server_addr, server_name, server_players, server_maxplayers) VALUES(?, ?, ?, ?)")
  if err != nil {
     return err
  }
  _, err = statement.Exec(addr, server.Hostname, server.Players, server.MaxPlayers)
  if err != nil {
    return err
  }
  return nil
}

func processLine(serversAddr *[]string, currentAddr *int, result *HistoryResult, db *sql.DB, waitGroup *sync.WaitGroup) {
  *currentAddr++
  if len(*serversAddr) <= *currentAddr {
    waitGroup.Done()
    return
  }
  processServerInfo((*serversAddr)[*currentAddr], result, db)
  processLine(serversAddr, currentAddr, result, db, waitGroup)
}

func processServerInfo(addr string, result *HistoryResult, db *sql.DB) {
  result.ServersTotal++
  var srvInfo query.ServerInfo
  for i := 1; i <= maxQueryAttempts; i++ {
    info, err := query.GetServerInfo(addr, queryTimeout)
    if err != nil {
      if i < maxQueryAttempts {
        log.Println(addr, " - error. Retrying...")
        continue
      } else {
        log.Println(addr, "- error:", err)
        return
      }
    }
    srvInfo = info
    break
  }
  if err := writeServerToDB(db, addr, &srvInfo); err != nil {
    log.Fatal("Error (writeServerToDB()):", err)
  }
  result.ServersOnline++
  result.SlotsTotal += srvInfo.MaxPlayers
  result.SlotsUsed += srvInfo.Players
  log.Println(addr, "- success!")
}
