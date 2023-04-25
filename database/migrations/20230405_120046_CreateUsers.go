package main

import (
  "github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type CreateUsers_20230405_120046 struct {
  migration.Migration
}

// DO NOT MODIFY
func init() {
  m := &CreateUsers_20230405_120046{}
  m.Created = "20230405_120046"

  migration.Register("CreateUsers_20230405_120046", m)
}

//curl -X POST http://localhost:8080/api/v1/auth/logout  -H 'Content-Type: application/json' -H 'Authorization: eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2ODEyODY4NjgsImlhdCI6MTY4MTI4MzI2OCwianRpIjoiNWVlM2YwOGEtZjJjYS00Nzg5LWEwMTAtMjUzNjczZjA1ZDcxIiwibmJmIjoxNjgxMjgzMjY4LCJzdWIiOjF9.VWqyYy3gywhFRRIEvhRZErhkImgXT9JLHWKugHCIKI_mt5nOjMvjiG1x6XMhsxPuf1oj0N0zzyYBnfCOpUcsaw'

// Run the migrations
func (m *CreateUsers_20230405_120046) Up() {
  m.SQL(`
    CREATE TABLE users (
      "id" serial NOT NULL PRIMARY KEY,
      "nickname" varchar(255) NOT NULL UNIQUE,
      "password" varchar(255) NOT NULL,
      "language_id" integer NOT NULL REFERENCES languages (id),
      "created_at" timestamp with time zone NOT NULL,
      "updated_at" timestamp with time zone NOT NULL,
      "next_game_available_at" timestamp with time zone NOT NULL,
      CHECK (length(password) >= 6)
    );
  `)
}

// Reverse the migrations
func (m *CreateUsers_20230405_120046) Down() {
  m.SQL(`
    DROP TABLE users;
    DELETE from migrations WHERE name = 'CreateUsers_20230405_120046';
  `)
}
