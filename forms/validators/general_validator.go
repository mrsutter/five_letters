package validators

import (
  "five_letters/models"
  "five_letters/utils"
  "github.com/beego/beego/v2/core/validation"
  "golang.org/x/exp/slices"
  "strings"
)

type GeneralValidator struct {
  BaseValidator
  Form interface{}
}

func (v GeneralValidator) Validate() models.Error {
  validator := validation.Validation{}
  result, _ := validator.Valid(v.Form)

  if result {
    return models.Error{}
  }

  errs := make([]models.ErrorItem, 0, len(validator.Errors))
  for _, err := range validator.Errors {
    msg := strings.ReplaceAll(err.Message, err.Field+" ", "")
    code := v.translateErrorCode(msg)

    jsonFieldName := utils.GetJsonFieldName(v.Form, err.Field)
    alreadyPresented := v.isFieldAlreadyPresented(errs, jsonFieldName)

    if !alreadyPresented {
      errs = append(errs, models.ErrorItem{Field: jsonFieldName, Code: code})
    }
  }

  return models.Error{
    Code:    v.translateErrorCode("inputErrors"),
    Details: errs,
  }
}

func (v *GeneralValidator) isFieldAlreadyPresented(
  errs []models.ErrorItem,
  field string) bool {

  idx := slices.IndexFunc(errs, func(i models.ErrorItem) bool {
    return i.Field == field
  })
  return idx != -1
}
