package validators

import (
  "five_letters/models"
  "github.com/beego/beego/v2/client/orm"
  "strings"
)

type WordExistanceValidator struct {
  BaseValidator
  Game models.Game
  Word string
}

func (v WordExistanceValidator) Validate() models.Error {
  o := orm.NewOrm()

  word := models.Word{
    WordList: v.Game.Word.WordList,
    Name:     strings.ToLower(v.Word),
  }

  err := o.Read(&word, "WordList", "Name")

  if err != nil {
    return v.customError()
  }

  if word.Archived == true && v.Game.Word.Id != word.Id {
    return v.customError()
  }

  return models.Error{}
}

func (v *WordExistanceValidator) customError() models.Error {
  eDetails := []models.ErrorItem{
    {
      Field: "word",
      Code:  v.translateErrorCode("noWordFound"),
    },
  }

  return models.Error{
    Code:    v.translateErrorCode("inputErrors"),
    Details: eDetails,
  }
}
