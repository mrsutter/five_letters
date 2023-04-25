package testfactories

import (
  "encoding/json"
  "five_letters/models"
  "github.com/beego/beego/v2/client/orm"
)

func CreateAttempt(
  number int,
  word string,
  game *models.Game) models.Attempt {

  attempt := models.Attempt{Number: number, Word: word, Game: game}
  result := attempt.CalcResult()

  jsonResult, _ := json.Marshal(result)
  attempt.Result = string(jsonResult)

  o := orm.NewOrm()
  o.Insert(&attempt)

  game.SetValuesAfterAttempt(attempt)
  o.Update(game)

  return attempt
}
