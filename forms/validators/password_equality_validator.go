package validators

import (
  "five_letters/models"
)

type PasswordEqualityValidator struct {
  BaseValidator
  Password             string
  PasswordConfirmation string
}

func (v PasswordEqualityValidator) Validate() models.Error {
  if v.Password == v.PasswordConfirmation {
    return models.Error{}
  }

  return v.customError()
}

func (v *PasswordEqualityValidator) customError() models.Error {
  eDetailsCode := v.translateErrorCode("passwordsAreNotEqual")

  eDetails := []models.ErrorItem{
    {
      Field: "password",
      Code:  eDetailsCode,
    },
    {
      Field: "password_confirmation",
      Code:  eDetailsCode,
    },
  }

  return models.Error{
    Code:    v.translateErrorCode("inputErrors"),
    Details: eDetails,
  }
}
