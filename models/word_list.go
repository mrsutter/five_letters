package models

import (
  "github.com/beego/beego/v2/client/orm"
  "time"
)

type WordList struct {
  Id        int
  Version   int
  Language  *Language `orm:"rel(one);on_delete(do_nothing)"`
  CreatedAt time.Time `orm:"auto_now_add;type(datetime)"`
  UpdatedAt time.Time `orm:"auto_now;type(datetime)"`
  Words     []*Word   `orm:"reverse(many);on_delete(do_nothing)"`
}

func init() {
  orm.RegisterModel(new(WordList))
}

func (w *WordList) TableName() string {
  return "word_lists"
}
