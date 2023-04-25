package forms

import (
  "five_letters/forms/validators"
  "five_letters/models"
  "github.com/beego/beego/v2/server/web/context"
)

type UserUpdateForm struct {
  BaseForm
  LanguageId int `json:"language_id" valid:"Required"`
}

type NewPetArg struct {
  UserUpdateForm *UserUpdateForm
}

func (form UserUpdateForm) Process(ctx *context.Context) models.Error {
  language := models.Language{Id: form.LanguageId}

  validators := []validators.Validator{
    validators.GeneralValidator{Form: form},
    validators.LangAvailabilityValidator{Language: &language},
  }
  err := form.Validate(validators)
  if err.Code != "" {
    return err
  }

  form.setValues(ctx, language)
  return models.Error{}
}

func (form UserUpdateForm) setValues(ctx *context.Context, language models.Language) {
  ctx.Input.SetData("language", language)
}
