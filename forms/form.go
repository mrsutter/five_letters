package forms

import (
  "five_letters/models"
  "github.com/beego/beego/v2/server/web/context"
)

type Form interface {
  Process(*context.Context) models.Error
}
