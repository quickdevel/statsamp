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
  // Соединение с БД
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.Database))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
  // Загрузка списка серверов с list.sa-mp.com
  log.Println("Запрос списка серверов с lists.sa-mp.com.")
  var serversAddr []string
  var currentAddr int
  request := serverlist.NewRequest(cfg.Updater.Version)
  if err := request.Exec(); err != nil {
    log.Fatal("Ошибка (request.Exec()):", err)
  }
  serversAddr, err = request.ReadAll()
  if err != nil {
    log.Fatal("Ошибка (request.ReadAll()):", err)
  }
  // Очистка таблицы в БД
  if err := clearServersInDB(db); err != nil {
    log.Fatal("Ошибка (clearServersInDB()):", err)
  }
  // Создание потоков для сбора статистики
  waitGroup.Add(maxQueryThreads)
  for i := 0; i < maxQueryThreads; i++ {
    go processLine(&serversAddr, &currentAddr, &historyResult, db, &waitGroup)
  }
  // Ожидание завершения потоков
  waitGroup.Wait()
  // Запись результатов в БД
  writeHistoryToDB(db, &historyResult)
  // Вывод результатов
  log.Println("Время выполнения:", time.Since(startTime))
  log.Println("Серверов онлайн:", historyResult.ServersOnline, "из", historyResult.ServersTotal)
  log.Println("Слотов занято:", historyResult.SlotsUsed, "из", historyResult.SlotsTotal)
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
        log.Println(addr, " - ошибка. Пробуем еще раз...")
        continue
      } else {
        log.Println(addr, "- ошибка:", err)
        return
      }
    }
    srvInfo = info
    break
  }
  if err := writeServerToDB(db, addr, &srvInfo); err != nil {
    log.Fatal("Ошибка (writeServerToDB()):", err)
  }
  result.ServersOnline++
  result.SlotsTotal += srvInfo.MaxPlayers
  result.SlotsUsed += srvInfo.Players
  log.Println(addr, "- успешно!")
}
