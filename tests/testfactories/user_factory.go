package testfactories

import (
  "five_letters/models"
  "five_letters/utils"
  "github.com/beego/beego/v2/client/orm"
  "time"
)

func CreateUser(
  nickname string,
  password string,
  nextGameAvailableAt time.Time,
  language models.Language) models.User {

  user := models.User{
    Nickname:            nickname,
    Password:            utils.HashPassword(password),
    NextGameAvailableAt: nextGameAvailableAt,
    Language:            &language,
  }

  o := orm.NewOrm()
  o.Insert(&user)

  return user
}

func CreateUserWithTokens(
  nickname string,
  password string,
  nextGameAvailableAt time.Time,
  language models.Language) (models.User, string, string) {

  user := CreateUser(nickname, password, nextGameAvailableAt, language)

  accessToken, refreshToken := CreateTokens(user, false)

  return user, accessToken, refreshToken
}
