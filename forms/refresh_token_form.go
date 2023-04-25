package forms

import (
  "five_letters/forms/validators"
  "five_letters/models"
  "github.com/beego/beego/v2/server/web/context"
)

type RefreshTokenForm struct {
  BaseForm
  Value string `json:"refresh_token" valid:"Required"`
}

func (form RefreshTokenForm) Process(ctx *context.Context) models.Error {
  token := models.RefreshToken{}

  validators := []validators.Validator{
    validators.GeneralValidator{Form: form},
    validators.RefreshTokenValidator{TokenValue: form.Value, Token: &token},
  }
  err := form.Validate(validators)
  if err.Code != "" {
    return err
  }

  form.setValues(ctx, token)
  return models.Error{}
}

func (form RefreshTokenForm) setValues(ctx *context.Context, token models.RefreshToken) {
  ctx.Input.SetData("refreshToken", token)
}
