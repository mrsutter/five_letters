package testfactories

import (
  "five_letters/models"
  "github.com/beego/beego/v2/client/orm"
)

func CreateWord(name string, wordList models.WordList) models.Word {
  word := models.Word{WordList: &wordList, Name: name}

  o := orm.NewOrm()
  o.Insert(&word)
  return word
}
