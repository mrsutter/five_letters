package main

import (
  "github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type CreateWordLists_20230403_114901 struct {
  migration.Migration
}

// DO NOT MODIFY
func init() {
  m := &CreateWordLists_20230403_114901{}
  m.Created = "20230403_114901"

  migration.Register("CreateWordLists_20230403_114901", m)
}

// Run the migrations
func (m *CreateWordLists_20230403_114901) Up() {
  m.SQL(`
    CREATE TABLE word_lists (
      "id" serial NOT NULL PRIMARY KEY,
      "language_id" integer NOT NULL REFERENCES languages (id),
      "version" integer NOT NULL,
      "created_at" timestamp with time zone NOT NULL,
      "updated_at" timestamp with time zone NOT NULL
    );
  `)
}

// Reverse the migrations
func (m *CreateWordLists_20230403_114901) Down() {
  m.SQL(`
    DROP TABLE word_lists;
    DELETE from migrations WHERE name = 'CreateWordLists_20230403_114901';
  `)
}
