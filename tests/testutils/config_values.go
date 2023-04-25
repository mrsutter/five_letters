package testutils

import (
  "github.com/beego/beego/v2/server/web"
)

func GetConfigStringValue(key string) string {
  str, _ := web.AppConfig.String(key)
  return str
}

func GetConfigIntValue(key string) int {
  intValue, _ := web.AppConfig.Int(key)
  return intValue
}

func GetConfigBoolValue(key string) bool {
  boolValue, _ := web.AppConfig.Bool(key)
  return boolValue
}
