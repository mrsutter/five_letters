package models

import (
  "github.com/beego/beego/v2/client/orm"
  "time"
)

type RefreshToken struct {
  Id          int
  Jti         string
  ExpiredAt   time.Time    `orm:"type(datetime)"`
  CreatedAt   time.Time    `orm:"auto_now_add;type(datetime)"`
  UpdatedAt   time.Time    `orm:"auto_now;type(datetime)"`
  User        *User        `orm:"rel(fk)"`
  AccessToken *AccessToken `orm:"reverse(one)"`
}

func init() {
  orm.RegisterModel(new(RefreshToken))
}

func (r *RefreshToken) TableName() string {
  return "refresh_tokens"
}
