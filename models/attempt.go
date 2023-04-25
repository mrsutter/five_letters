package models

import (
  "github.com/beego/beego/v2/client/orm"
  "github.com/beego/i18n"
  "time"
)

type Attempt struct {
  Id        int
  Number    int
  Word      string
  Result    string    `orm:"type(jsonb)"`
  CreatedAt time.Time `orm:"auto_now_add;type(datetime)"`
  UpdatedAt time.Time `orm:"auto_now;type(datetime)"`
  Game      *Game     `orm:"rel(fk)"`
}

func (a *Attempt) CalcResult() []string {
  puzzledWord := a.Game.Word.Name

  result := make([]string, WordMaxLength)
  chars := []rune(a.Word)
  puzzledChars := []rune(puzzledWord)

  unMatchedPositions := make(map[rune][]int)

  for index, char := range puzzledChars {
    if char == chars[index] {
      result[index] = i18n.Tr("en", "attemptLetterStatuses.match")
      continue
    }
    unMatchedPositions[char] = append(unMatchedPositions[char], index)
  }

  for index, char := range chars {
    if result[index] != "" {
      continue
    }

    if len(unMatchedPositions[char]) != 0 {
      unMatchedPositions[char] = unMatchedPositions[char][1:]
      result[index] = i18n.Tr("en", "attemptLetterStatuses.wrongPlace")
    } else {
      result[index] = i18n.Tr("en", "attemptLetterStatuses.absence")
    }
  }
  return result
}

func (a *Attempt) Successful() bool {
  return a.Word == a.Game.Word.Name
}

func init() {
  orm.RegisterModel(new(Attempt))
}

func (a *Attempt) TableName() string {
  return "attempts"
}
