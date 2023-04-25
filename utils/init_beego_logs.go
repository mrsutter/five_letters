package utils

import (
  "github.com/beego/beego/v2/core/logs"
  "github.com/beego/beego/v2/server/web"
)

func InitBeegoLogs(logSettingsString string) {
  web.BConfig.Log.AccessLogs = true
  logs.SetLogger(logs.AdapterFile, logSettingsString)
}
