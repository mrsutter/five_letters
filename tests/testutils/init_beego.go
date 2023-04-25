package testutils

import (
  "five_letters/utils"
  "github.com/beego/beego/v2/server/web"
  "github.com/beego/i18n"
  "path/filepath"
  "runtime"
)

func InitBeego() {
  rootPath := getRootPath()
  web.TestBeegoInit(rootPath)
  initORM()
  initLogs()
  initLocales(rootPath)
  initDocs()
}

func getRootPath() string {
  _, file, _, _ := runtime.Caller(0)
  rootPath := filepath.Join(filepath.Dir(file), "../..")
  return rootPath
}

func initORM() {
  dbConnString := GetConfigStringValue("dbConnectionString")
  dbDriver := GetConfigStringValue("dbDriver")
  dbAlias := GetConfigStringValue("dbAlias")
  dbMaxIdleConns := GetConfigIntValue("dbMaxIdleConnections")
  dbMaxConns := GetConfigIntValue("dbMaxConnections")
  dbLogQueries := GetConfigBoolValue("dbLogQueries")

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
  logSettingsString := GetConfigStringValue("logSettingsString")
  utils.InitBeegoLogs(logSettingsString)
}

func initLocales(rootPath string) {
  _ = i18n.SetMessage("en", filepath.Join(rootPath, "conf/locale_en.ini"))
}

func initDocs() {
  web.BConfig.WebConfig.DirectoryIndex = true
  web.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
}
