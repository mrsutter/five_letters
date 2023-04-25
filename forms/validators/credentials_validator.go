package validators

import (
  "five_letters/models"
  "five_letters/utils"
  "github.com/beego/beego/v2/client/orm"
)

type CredentialsValidator struct {
  BaseValidator
  Nickname string
  Password string
}

func (v CredentialsValidator) Validate() models.Error {
  o := orm.NewOrm()

  user := models.User{Nickname: v.Nickname}

  err := o.Read(&user, "Nickname")
  if err != nil {
    return v.customError()
  }

  err = utils.ComparePassword(user.Password, v.Password)
  if err != nil {
    return v.customError()
  }

  return models.Error{}
}

func (v *CredentialsValidator) customError() models.Error {
  eDetailsCode := v.translateErrorCode("nickAndPassDoNotMatch")
  eDetails := []models.ErrorItem{
    {
      Field: "nickname",
      Code:  eDetailsCode,
    },
    {
      Field: "password",
      Code:  eDetailsCode,
    },
  }

  return models.Error{
    Code:    v.translateErrorCode("inputErrors"),
    Details: eDetails,
  }
}
