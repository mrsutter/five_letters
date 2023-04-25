package main

import (
  _ "five_letters/routers"
  "five_letters/tasks"
  "five_letters/utils"
  "github.com/beego/beego/v2/server/web"
  "github.com/beego/beego/v2/task"
  "github.com/beego/i18n"
)

func initORM() {
  dbConnString, _ := web.AppConfig.String("dbConnectionString")
  dbDriver, _ := web.AppConfig.String("dbDriver")
  dbAlias, _ := web.AppConfig.String("dbAlias")
  dbMaxIdleConns, _ := web.AppConfig.Int("dbMaxIdleConnections")
  dbMaxConns, _ := web.AppConfig.Int("dbMaxConnections")
  dbLogQueries, _ := web.AppConfig.Bool("dbLogQueries")

  utils.InitBeegoOrm(
    dbConnString,
    dbDriver,
    dbAlias,
    dbMaxIdleConns,
    dbMaxConns,
    dbLogQueries,
  )
}

func initLogs() {
  logSettingsString, _ := web.AppConfig.String("logSettingsString")
  utils.InitBeegoLogs(logSettingsString)
}

func initLocales() {
  _ = i18n.SetMessage("en", "conf/locale_en.ini")
}

func init() {
  initORM()
  initLogs()
  initLocales()
}

func main() {
  for _, t := range tasks.InitTasks() {
    newTask := task.NewTask(t.Id, t.RunAt, t.Func)
    task.AddTask(t.Id, newTask)
  }
  task.StartTask()
  defer task.StopTask()

  web.Run()
}
