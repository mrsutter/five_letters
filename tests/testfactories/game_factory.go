package testfactories

import (
  "five_letters/models"
  "github.com/beego/beego/v2/client/orm"
)

func CreateGame(
  user models.User,
  word models.Word,
  state string,
  attemptsCount int) models.Game {

  game := models.Game{
    User:          &user,
    Word:          &word,
    State:         state,
    AttemptsCount: attemptsCount,
  }

  o := orm.NewOrm()
  o.Insert(&game)
  return game
}

func CreateGameWithAttempts(
  user models.User,
  word models.Word,
  state string,
  attemptWords ...string) models.Game {

  game := CreateGame(user, word, state, 0)

  for i, attemptWord := range attemptWords {
    CreateAttempt(i+1, attemptWord, &game)
  }

  return game
}
