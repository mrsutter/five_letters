package forms

import (
  "five_letters/forms/validators"
  "five_letters/models"
  "github.com/beego/beego/v2/server/web/context"
)

type UserSignInForm struct {
  BaseForm
  Nickname string `json:"nickname" valid:"Required;AlphaDash"`
  Password string `json:"password" valid:"Required"`
}

func (form UserSignInForm) Process(ctx *context.Context) models.Error {
  validators := []validators.Validator{
    validators.GeneralValidator{Form: form},
    validators.CredentialsValidator{
      Nickname: form.Nickname,
      Password: form.Password,
    },
  }
  err := form.Validate(validators)
  if err.Code != "" {
    return err
  }

  return models.Error{}
}
