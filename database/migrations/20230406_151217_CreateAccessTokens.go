package main

import (
  "github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type CreateAccessTokens_20230406_151217 struct {
  migration.Migration
}

// DO NOT MODIFY
func init() {
  m := &CreateAccessTokens_20230406_151217{}
  m.Created = "20230406_151217"

  migration.Register("CreateAccessTokens_20230406_151217", m)
}

// Run the migrations
func (m *CreateAccessTokens_20230406_151217) Up() {
  m.SQL(`
    CREATE TABLE "access_tokens" (
      "id" serial NOT NULL PRIMARY KEY,
      "jti" varchar(255) NOT NULL UNIQUE,
      "expired_at" timestamp with time zone NOT NULL,
      "created_at" timestamp with time zone NOT NULL,
      "updated_at" timestamp with time zone NOT NULL,
      "refresh_token_id" integer NOT NULL REFERENCES refresh_tokens (id) ON DELETE CASCADE,
      "user_id" integer NOT NULL REFERENCES users (id) ON DELETE CASCADE
    );
    CREATE INDEX idx_access_tokens_user ON access_tokens (user_id);
    CREATE INDEX idx_access_tokens_refresh_token ON access_tokens (refresh_token_id);
  `)
}

// Reverse the migrations
func (m *CreateAccessTokens_20230406_151217) Down() {
  m.SQL(`
    DROP INDEX idx_access_tokens_user;
    DROP INDEX idx_access_tokens_refresh_token;
    DROP TABLE access_tokens;
    DELETE from migrations WHERE name = 'CreateAccessTokens_20230406_151217';
  `)
}
