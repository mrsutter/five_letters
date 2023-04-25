package main

import (
  "github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type CreateWords_20230403_120311 struct {
  migration.Migration
}

// DO NOT MODIFY
func init() {
  m := &CreateWords_20230403_120311{}
  m.Created = "20230403_120311"

  migration.Register("CreateWords_20230403_120311", m)
}

// Run the migrations
func (m *CreateWords_20230403_120311) Up() {
  m.SQL(`
    CREATE TABLE words (
      "id" serial NOT NULL PRIMARY KEY,
      "name" varchar(255) NOT NULL UNIQUE,
      "archived" BOOLEAN NOT NULL DEFAULT FALSE,
      "created_at" timestamp with time zone NOT NULL,
      "updated_at" timestamp with time zone NOT NULL,
      "word_list_id" integer NOT NULL REFERENCES word_lists (id),
      CHECK (length(name) = 5)
    );
    CREATE INDEX idx_words_list ON words (word_list_id);
    CREATE INDEX idx_words_list_updated_at ON words (word_list_id, archived);
    CREATE INDEX idx_words_updated_at ON words (updated_at);
  `)
}

// Reverse the migrations
func (m *CreateWords_20230403_120311) Down() {
  m.SQL(`
    DROP INDEX idx_words_updated_at;
    DROP INDEX idx_words_list;
    DROP INDEX idx_words_list_updated_at;
    DROP TABLE words;
    DELETE from migrations WHERE name = 'CreateWords_20230403_120311';
  `)
}
