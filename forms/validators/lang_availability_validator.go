package validators

import (
  "five_letters/models"
  "github.com/beego/beego/v2/client/orm"
)

type LangAvailabilityValidator struct {
  BaseValidator
  Language *models.Language
}

func (v LangAvailabilityValidator) Validate() models.Error {
  o := orm.NewOrm()

  err := o.Read(v.Language)
  if err != nil || !v.Language.Available {
    return v.customError()
  }

  return models.Error{}
}

func (v *LangAvailabilityValidator) customError() models.Error {
  eDetails := []models.ErrorItem{
    {
      Field: "language_id",
      Code:  v.translateErrorCode("noLanguageFound"),
    },
  }

  return models.Error{
    Code:    v.translateErrorCode("inputErrors"),
    Details: eDetails,
  }
}
