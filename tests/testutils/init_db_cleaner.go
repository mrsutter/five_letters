package testutils

import (
  "gopkg.in/khaiql/dbcleaner.v2"
  "gopkg.in/khaiql/dbcleaner.v2/engine"
)

func InitDbCleaner(cleaner dbcleaner.DbCleaner) {
  dbConnString := GetConfigStringValue("dbConnectionString")
  engine := engine.NewPostgresEngine(dbConnString)
  cleaner.SetEngine(engine)
}
