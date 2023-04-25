package models

import (
  "github.com/beego/beego/v2/client/orm"
  "time"
)

type Language struct {
  Id          int
  Slug        string
  Name        string
  LettersList string
  Available   bool
  CreatedAt   time.Time `orm:"auto_now_add;type(datetime)"`
  UpdatedAt   time.Time `orm:"auto_now;type(datetime)"`
  WordList    *WordList `orm:"reverse(one)"`
  Users       []*User   `orm:"reverse(many);on_delete(do_nothing)"`
}

func init() {
  orm.RegisterModel(new(Language))
}

func (l *Language) TableName() string {
  return "languages"
}
