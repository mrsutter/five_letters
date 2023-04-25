package validators

import (
  "five_letters/models"
  "github.com/beego/beego/v2/client/orm"
)

type NicknameExistanceValidator struct {
  BaseValidator
  Nickname string
}

func (v NicknameExistanceValidator) Validate() models.Error {
  o := orm.NewOrm()

  user := models.User{Nickname: v.Nickname}
  err := o.Read(&user, "Nickname")
  if err != nil {
    return models.Error{}
  }

  return v.customError()
}

func (v *NicknameExistanceValidator) customError() models.Error {
  eDetails := []models.ErrorItem{
    {
      Field: "nickname",
      Code:  v.translateErrorCode("alreadyTaken"),
    },
  }

  return models.Error{
    Code:    v.translateErrorCode("inputErrors"),
    Details: eDetails,
  }
}
