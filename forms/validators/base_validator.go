package validators

import (
  "github.com/beego/i18n"
)

type BaseValidator struct {
}

func (v *BaseValidator) translateErrorCode(e string) string {
  return i18n.Tr("en", "errors."+e)
}
