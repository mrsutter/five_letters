package main

import (
  "github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type CreateAttempts_20230407_111258 struct {
  migration.Migration
}

// DO NOT MODIFY
func init() {
  m := &CreateAttempts_20230407_111258{}
  m.Created = "20230407_111258"

  migration.Register("CreateAttempts_20230407_111258", m)
}

// Run the migrations
func (m *CreateAttempts_20230407_111258) Up() {
  m.SQL(`
    CREATE TABLE "attempts" (
      "id" serial NOT NULL PRIMARY KEY,
      "number" integer NOT NULL,
      "word" varchar(255) NOT NULL,
      "result" jsonb NOT NULL,
      "created_at" timestamp with time zone NOT NULL,
      "updated_at" timestamp with time zone NOT NULL,
      "game_id" integer NOT NULL REFERENCES games (id) ON DELETE CASCADE,
      CHECK (length(word) = 5),
      CHECK (number >= 1 AND number <= 6),
      CONSTRAINT unique_attempt_number_in_game UNIQUE (number, game_id)
    );
    CREATE INDEX idx_attempts_game ON attempts (game_id);
  `)
}

// Reverse the migrations
func (m *CreateAttempts_20230407_111258) Down() {
  m.SQL(`
    DROP INDEX idx_attempts_game;
    DROP TABLE attempts;
    DELETE from migrations WHERE name = 'CreateAttempts_20230407_111258';
  `)
}
