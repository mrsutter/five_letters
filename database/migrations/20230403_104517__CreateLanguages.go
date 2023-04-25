package main

import (
  "github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type CreateLanguages_20230403_104517 struct {
  migration.Migration
}

// DO NOT MODIFY
func init() {
  m := &CreateLanguages_20230403_104517{}
  m.Created = "20230403_104517"

  migration.Register("CreateLanguages_20230403_104517", m)
}

// Run the migrations
func (m *CreateLanguages_20230403_104517) Up() {
  m.SQL(`
    CREATE TABLE languages (
      "id" serial NOT NULL PRIMARY KEY,
      "slug" varchar(255) NOT NULL UNIQUE,
      "name" varchar(255) NOT NULL,
      "letters_list" varchar(255) NOT NULL,
      "available" BOOLEAN NOT NULL DEFAULT TRUE,
      "created_at" timestamp with time zone NOT NULL,
      "updated_at" timestamp with time zone NOT NULL
    );
  `)

}

// Reverse the migrations
func (m *CreateLanguages_20230403_104517) Down() {
  m.SQL(`
    DROP TABLE languages;
    DELETE from migrations WHERE name = 'CreateLanguages_20230403_104517';
  `)
}
