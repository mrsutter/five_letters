package models

import (
  "github.com/beego/beego/v2/client/orm"
  "time"
)

type AccessToken struct {
  Id           int
  Jti          string
  ExpiredAt    time.Time     `orm:"type(datetime)"`
  CreatedAt    time.Time     `orm:"auto_now_add;type(datetime)"`
  UpdatedAt    time.Time     `orm:"auto_now;type(datetime)"`
  User         *User         `orm:"rel(fk)"`
  RefreshToken *RefreshToken `orm:"rel(one);on_delete(cascade)"`
}

func init() {
  orm.RegisterModel(new(AccessToken))
}

func (a *AccessToken) TableName() string {
  return "access_tokens"
}
