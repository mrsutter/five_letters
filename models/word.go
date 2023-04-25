package models

import (
  "github.com/beego/beego/v2/client/orm"
  "time"
)

const WordMaxLength = 5

type Word struct {
  Id        int
  Name      string
  Archived  bool
  CreatedAt time.Time `orm:"auto_now_add;type(datetime)"`
  UpdatedAt time.Time `orm:"auto_now;type(datetime)"`
  WordList  *WordList `orm:"rel(fk)"`
  Games     []*Game   `orm:"reverse(many);on_delete(do_nothing)"`
}

func init() {
  orm.RegisterModel(new(Word))
}

func (w *Word) TableName() string {
  return "words"
}
