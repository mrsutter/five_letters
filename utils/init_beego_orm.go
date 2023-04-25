package utils

import (
  "github.com/beego/beego/v2/client/orm"
  _ "github.com/lib/pq"
  "time"
)

func InitBeegoOrm(dbConnString string, dbDriver string, dbAlias string, dbMaxIdleConns int, dbMaxConns int, dbLogQueries bool) {
  orm.RegisterDriver(dbDriver, orm.DRPostgres)
  orm.RegisterDataBase(dbAlias, dbDriver, dbConnString)
  orm.SetMaxIdleConns(dbAlias, dbMaxIdleConns)
  orm.SetMaxOpenConns(dbAlias, dbMaxConns)
  orm.DefaultTimeLoc = time.UTC
  orm.Debug = dbLogQueries
}
