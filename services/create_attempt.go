package services

import (
  "encoding/json"
  "five_letters/models"
  "github.com/beego/beego/v2/client/orm"
  "strings"
)

func CreateAttempt(game *models.Game, word string) (models.Attempt, error) {
  o := orm.NewOrm()
  tr, err := o.Begin()
  if err != nil {
    return models.Attempt{}, err
  }
  defer func() {
    if err != nil {
      tr.Rollback()
    } else {
      tr.Commit()
    }
  }()

  attempt := models.Attempt{
    Word:   strings.ToLower(word),
    Number: game.AttemptsCount + 1,
    Game:   game,
  }
  attemptResult, _ := json.Marshal(attempt.CalcResult())
  attempt.Result = string(attemptResult)
  _, err = tr.Insert(&attempt)
  if err != nil {
    return models.Attempt{}, err
  }

  game.SetValuesAfterAttempt(attempt)
  _, err = tr.Update(game, "State", "AttemptsCount")
  if err != nil {
    return models.Attempt{}, err
  }

  return attempt, nil
}
