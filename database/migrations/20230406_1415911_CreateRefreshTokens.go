package main

import (
  "github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type CreateRefreshTokens_20230406_145911 struct {
  migration.Migration
}

// DO NOT MODIFY
func init() {
  m := &CreateRefreshTokens_20230406_145911{}
  m.Created = "20230406_145911"

  migration.Register("CreateRefreshTokens_20230406_145911", m)
}

// Run the migrations
func (m *CreateRefreshTokens_20230406_145911) Up() {
  m.SQL(`
    CREATE TABLE "refresh_tokens" (
      "id" serial NOT NULL PRIMARY KEY,
      "jti" varchar(255) NOT NULL UNIQUE,
      "expired_at" timestamp with time zone NOT NULL,
      "created_at" timestamp with time zone NOT NULL,
      "updated_at" timestamp with time zone NOT NULL,
      "user_id" integer NOT NULL REFERENCES users (id) ON DELETE CASCADE
    );
    CREATE INDEX idx_refresh_tokens_user ON refresh_tokens(user_id);
  `)
}

// Reverse the migrations
func (m *CreateRefreshTokens_20230406_145911) Down() {
  m.SQL(`
    DROP INDEX idx_refresh_tokens_user;
    DROP TABLE refresh_tokens;
    DELETE from migrations WHERE name = 'CreateRefreshTokens_20230406_145911';
  `)
}
