package forms

import (
  "five_letters/forms/validators"
  "five_letters/models"
  "github.com/beego/beego/v2/server/web/context"
)

type UserSignUpForm struct {
  BaseForm
  Nickname             string `json:"nickname" valid:"Required;AlphaDash"`
  Password             string `json:"password" valid:"Required;MinSize(6)"`
  PasswordConfirmation string `json:"password_confirmation" valid:"Required;MinSize(6)"`
  LanguageId           int    `json:"language_id" valid:"Required"`
}

func (form UserSignUpForm) Process(ctx *context.Context) models.Error {
  language := models.Language{Id: form.LanguageId}

  validators := []validators.Validator{
    validators.GeneralValidator{Form: form},
    validators.PasswordEqualityValidator{
      Password:             form.Password,
      PasswordConfirmation: form.PasswordConfirmation,
    },
    validators.NicknameExistanceValidator{Nickname: form.Nickname},
    validators.LangAvailabilityValidator{Language: &language},
  }
  err := form.Validate(validators)
  if err.Code != "" {
    return err
  }

  form.setValues(ctx, language)
  return models.Error{}
}

func (form UserSignUpForm) setValues(ctx *context.Context, language models.Language) {
  ctx.Input.SetData("language", language)
}
