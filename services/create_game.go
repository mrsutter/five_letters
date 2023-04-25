package services

import (
  "five_letters/models"
  "github.com/beego/beego/v2/client/orm"
  "time"
)

func CreateGame(user models.User) (models.Game, error) {
  o := orm.NewOrm()
  tr, err := o.Begin()
  if err != nil {
    return models.Game{}, err
  }
  defer func() {
    if err != nil {
      tr.Rollback()
    } else {
      tr.Commit()
    }
  }()

  word := getWord(user)

  game := models.Game{
    User:  &user,
    Word:  &word,
    State: models.StateActive,
  }

  _, err = tr.Insert(&game)
  if err != nil {
    return models.Game{}, err
  }

  user.NextGameAvailableAt = game.CreatedAt.Add(models.MaxHoursForGame * time.Hour)

  _, err = tr.Update(&user, "NextGameAvailableAt")
  if err != nil {
    return models.Game{}, err
  }

  return game, nil
}

func getWord(user models.User) models.Word {
  o := orm.NewOrm()

  var word models.Word

  err := o.Raw(
    sqlUnusedRandomWord(),
    user.Id,
    user.Language.Id).
    QueryRow(&word)

  if err != nil {
    o.Raw(
      sqlRandomWord(),
      user.Language.Id).
      QueryRow(&word)
  }

  return word
}

func sqlUnusedRandomWord() string {
  sql := `
    SELECT words.*
      FROM words
      JOIN word_lists ON words.word_list_id = word_lists.id
      LEFT JOIN games ON words.id = games.word_id AND games.user_id = ?
      WHERE word_lists.language_id = ?
        AND words.archived = false
        AND games.id IS NULL
      ORDER BY RANDOM()
      LIMIT 1;
  `
  return sql
}

func sqlRandomWord() string {
  sql := `
    SELECT words.*
      FROM words
      JOIN word_lists ON words.word_list_id = word_lists.id
      WHERE word_lists.language_id = ?
        AND words.archived = false
      ORDER BY RANDOM()
      LIMIT 1;
  `
  return sql
}
