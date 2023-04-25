package controllers

import (
  "encoding/json"
  "five_letters/forms"
  "five_letters/models"
  "five_letters/utils"
  "github.com/beego/beego/v2/server/web"
  "github.com/beego/i18n"
)

type BaseController struct {
  web.Controller
}

func (c *BaseController) Abort422(err models.Error) {
  c.Ctx.Input.SetData("error", err)
  c.Abort("422")
}

func (c *BaseController) parseBody(body []byte, form interface{}) {
  err := json.Unmarshal(body, form)
  if err != nil {
    c.Abort422(models.Error{
      Code:    c.translateErrorCode("badJSON"),
      Message: err.Error(),
    })
  }
}

func (c *BaseController) translateErrorCode(e string) string {
  return i18n.Tr("en", "errors."+e)
}

func (c *BaseController) processForm(form forms.Form) {
  err := form.Process(c.Ctx)

  if err.Code != "" {
    c.Abort422(err)
  }
}

func (c *BaseController) currentUser() models.User {
  return c.Ctx.Input.GetData("currentUser").(models.User)
}

func (c *BaseController) setHeaders(user models.User) {
  nextGameAvailableAt := utils.FormatTime(user.NextGameAvailableAt)
  c.Ctx.Output.Header("Next-Game-Available-At", nextGameAvailableAt)
}
