package models

import (
  "github.com/beego/beego/v2/client/orm"
  "time"
)

type User struct {
  Id                  int
  Nickname            string
  Password            string
  Language            *Language       `orm:"rel(fk)"`
  CreatedAt           time.Time       `orm:"auto_now_add;type(datetime)"`
  UpdatedAt           time.Time       `orm:"auto_now;type(datetime)"`
  NextGameAvailableAt time.Time       `orm:"type(datetime)"`
  Games               []*Game         `orm:"reverse(many);on_delete(cascade)"`
  AccessTokens        []*AccessToken  `orm:"reverse(many);on_delete(cascade)"`
  RefreshTokens       []*RefreshToken `orm:"reverse(many);on_delete(cascade)"`
}

func init() {
  orm.RegisterModel(new(User))
}

func (u *User) TableName() string {
  return "users"
}
