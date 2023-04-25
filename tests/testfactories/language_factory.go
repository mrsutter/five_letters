package testfactories

import (
  "five_letters/models"
  "github.com/beego/beego/v2/client/orm"
)

func CreateEnLanguage(available bool) models.Language {
  return CreateLanguage("en", "English", "^[a-z]+$", available)
}

func CreateRuLanguage(available bool) models.Language {
  return CreateLanguage("ru", "Русский", "^[а-я]+$", available)
}

func CreateLanguage(
  slug string,
  name string,
  lettersList string,
  available bool) models.Language {

  lang := models.Language{
    Slug:        slug,
    Name:        name,
    LettersList: lettersList,
    Available:   available,
  }

  o := orm.NewOrm()
  o.Insert(&lang)
  return lang
}
