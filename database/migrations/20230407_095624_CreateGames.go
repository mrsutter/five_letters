package main

import (
  "github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type CreateGames_20230407_095624 struct {
  migration.Migration
}

// DO NOT MODIFY
func init() {
  m := &CreateGames_20230407_095624{}
  m.Created = "20230407_095624"

  migration.Register("CreateGames_20230407_095624", m)
}

// Run the migrations
func (m *CreateGames_20230407_095624) Up() {
  m.SQL(`
    CREATE TABLE "games" (
      "id" serial NOT NULL PRIMARY KEY,
      "state" varchar(255) NOT NULL DEFAULT 'active',
      "attempts_count" integer NOT NULL DEFAULT 0 ,
      "created_at" timestamp with time zone NOT NULL,
      "updated_at" timestamp with time zone NOT NULL,
      "word_id" integer NOT NULL REFERENCES words (id),
      "user_id" integer NOT NULL REFERENCES users (id) ON DELETE CASCADE,
      CHECK (state IN ('active', 'wasted', 'won'))
    );
    CREATE INDEX idx_games_user ON games (user_id);
    CREATE UNIQUE INDEX idx_games_one_active_for_user on games (user_id) WHERE state = 'active';
  `)
}

// Reverse the migrations
func (m *CreateGames_20230407_095624) Down() {
  m.SQL(`
    DROP INDEX idx_games_one_active_for_user;
    DROP INDEX idx_games_user;
    DROP TABLE games;
    DELETE from migrations WHERE name = 'CreateGames_20230407_095624';
  `)
}
