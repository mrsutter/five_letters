package testfactories

import (
  "five_letters/models"
  "github.com/beego/beego/v2/client/orm"
)

func CreateWordList(language models.Language, version int) models.WordList {
  list := models.WordList{Language: &language, Version: version}

  o := orm.NewOrm()
  o.Insert(&list)
  return list
}
