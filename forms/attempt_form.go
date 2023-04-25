package forms

import (
  "five_letters/forms/validators"
  "five_letters/models"
  "github.com/beego/beego/v2/server/web/context"
)

type AttemptForm struct {
  BaseForm
  Word string `json:"word" valid:"Required"`
}

func (form AttemptForm) Process(ctx *context.Context) models.Error {
  game := ctx.Input.GetData("game").(models.Game)

  validators := []validators.Validator{
    validators.GeneralValidator{Form: form},
    validators.WordExistanceValidator{Game: game, Word: form.Word},
  }
  err := form.Validate(validators)
  if err.Code != "" {
    return err
  }

  return models.Error{}
}
