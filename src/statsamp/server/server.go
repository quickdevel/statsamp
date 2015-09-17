package main

import (
  "fmt"
  "log"
  "net/http"
  // DB libs
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
  // Gorilla libs
  "github.com/gorilla/mux"
  // StatSAMP libs
  "statsamp/config"
  "statsamp/response"
)

type GeneralData struct {
  ServersTotal    int     `json:"servers_total"`
  ServersOnline   int     `json:"servers_online"`
  ServersOffline  int     `json:"servers_offline"`
  SlotsTotal      int     `json:"slots_total"`
  SlotsUsed       int     `json:"slots_used"`
  SlotsUnused     int     `json:"slots_unused"`
  Date            string  `json:"date"`
}

type HistoryServers struct {
  Total  int     `json:"total"`
  Date   string  `json:"date"`
}

type HistoryPlayers struct {
  Total  int     `json:"total"`
  Date   string  `json:"date"`
}

func main() {
  // Load config
  cfg, err := config.GetConfigData()
  if err != nil {
    log.Fatal(err)
  }
  // Initialize database
  db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.Database))
  if err != nil {
    log.Fatal(err)
  }
  defer db.Close()

  // Initialize HTTP-server
  router := mux.NewRouter()
  router.HandleFunc("/api/general", func(w http.ResponseWriter, r *http.Request) {
    data, err := GetGeneralData(db)
    if err != nil {
      response.JSONResponseError(&w, http.StatusInternalServerError, err.Error())
      return
    }
    response.JSONResponse(&w, http.StatusOK, data)
  })

  router.HandleFunc("/api/history/servers/{range}", func(w http.ResponseWriter, r *http.Request) {
    var days int
    vars := mux.Vars(r)
    requestRange := vars["range"]
    switch requestRange {
      case "day":
        days = 1
      case "week":
        days = 7
      case "month":
        days = 30
      case "year":
        days = 365
      default:
        response.JSONResponseError(&w, http.StatusInternalServerError, "Incorrect range")
        return
    }
    data, err := GetHistoryServers(db, days)
    if err != nil {
      response.JSONResponseError(&w, http.StatusInternalServerError, err.Error())
      return
    }
    response.JSONResponse(&w, http.StatusOK, data)
  })

  router.HandleFunc("/api/history/players/{range}", func(w http.ResponseWriter, r *http.Request) {
    var days int
    vars := mux.Vars(r)
    requestRange := vars["range"]
    switch requestRange {
      case "day":
        days = 1
      case "week":
        days = 7
      case "month":
        days = 30
      case "year":
        days = 365
      default:
        response.JSONResponseError(&w, http.StatusInternalServerError, "Incorrect range")
        return
    }
    data, err := GetHistoryPlayers(db, days)
    if err != nil {
      response.JSONResponseError(&w, http.StatusInternalServerError, err.Error())
      return
    }
    response.JSONResponse(&w, http.StatusOK, data)
  })

  router.PathPrefix("/").Handler(http.FileServer(http.Dir("./public")))
  http.Handle("/", router)
  log.Println(fmt.Sprintf("Listening %s...", cfg.Main.Listen))
  if err := http.ListenAndServe(cfg.Main.Listen, nil); err != nil {
    log.Fatal("ListenAndServe: ", err)
  }
}

func GetGeneralData(db *sql.DB) (GeneralData, error) {
  var data GeneralData
  statement := db.QueryRow("SELECT history_servers_total, history_servers_online, history_slots_total, history_slots_used, history_date FROM history ORDER BY history_id DESC LIMIT 1")
  if err := statement.Scan(&data.ServersTotal, &data.ServersOnline, &data.SlotsTotal, &data.SlotsUsed, &data.Date); err != nil {
    return data, err
  }
  data.ServersOffline = data.ServersTotal - data.ServersOnline
  data.SlotsUnused = data.SlotsTotal - data.SlotsUsed
  return data, nil
}

func GetHistoryServers(db *sql.DB, days int) ([]HistoryServers, error) {
  var data []HistoryServers
  var query string
  switch days {
    case 1:
      query = "SELECT history_servers_online, history_date FROM history WHERE history_date > NOW() - INTERVAL 1 DAY"
    case 7:
      query = "SELECT a.history_servers_online, a.history_date FROM (SELECT history_servers_online, history_date, DATE_FORMAT(history_date, '%p%d%m%y') AS g FROM history WHERE history_date > NOW() - INTERVAL 1 WEEK GROUP BY g) AS a"
    case 30:
      query = "SELECT a.history_servers_online, a.history_date FROM (SELECT history_servers_online, history_date, DATE_FORMAT(history_date, '%d%m%y') AS g FROM history WHERE history_date > NOW() - INTERVAL 1 MONTH GROUP BY g) AS a"
    case 365:
      query = "SELECT a.history_servers_online, a.history_date FROM (SELECT history_servers_online, history_date, DATE_FORMAT(history_date, '%u%y') AS g FROM history WHERE history_date > NOW() - INTERVAL 1 YEAR GROUP BY g) AS a"
  }
  rows, err := db.Query(query)
  if err != nil {
	  return data, err
  }
  defer rows.Close()
  for rows.Next() {
    var tmp HistoryServers
	  if err := rows.Scan(&tmp.Total, &tmp.Date); err != nil {
      return data, err
    }
    data = append(data, tmp)
  }
  return data, nil
}

func GetHistoryPlayers(db *sql.DB, days int) ([]HistoryPlayers, error) {
  var data []HistoryPlayers
  var query string
  switch days {
    case 1:
      query = "SELECT history_slots_used, history_date FROM history WHERE history_date > NOW() - INTERVAL 1 DAY"
    case 7:
      query = "SELECT a.history_slots_used, a.history_date FROM (SELECT history_slots_used, history_date, DATE_FORMAT(history_date, '%p%d%m%y') AS g FROM history WHERE history_date > NOW() - INTERVAL 1 WEEK GROUP BY g) AS a"
    case 30:
      query = "SELECT a.history_slots_used, a.history_date FROM (SELECT history_slots_used, history_date, DATE_FORMAT(history_date, '%d%m%y') AS g FROM history WHERE history_date > NOW() - INTERVAL 1 MONTH GROUP BY g) AS a"
    case 365:
      query = "SELECT a.history_slots_used, a.history_date FROM (SELECT history_slots_used, history_date, DATE_FORMAT(history_date, '%u%y') AS g FROM history WHERE history_date > NOW() - INTERVAL 1 YEAR GROUP BY g) AS a"
  }
  rows, err := db.Query(query)
  if err != nil {
    return data, err
  }
  defer rows.Close()
  for rows.Next() {
    var tmp HistoryPlayers
    if err := rows.Scan(&tmp.Total, &tmp.Date); err != nil {
      return data, err
    }
    data = append(data, tmp)
  }
  return data, nil
}
