package forms

import (
  "five_letters/forms/validators"
  "five_letters/models"
)

type BaseForm struct {
}

func (form BaseForm) Validate(validators []validators.Validator) models.Error {
  for _, v := range validators {
    err := v.Validate()

    if err.Code != "" {
      return err
    }
  }

  return models.Error{}
}
