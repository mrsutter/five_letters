package controllers

import (
  "five_letters/forms"
  "five_letters/models"
  "five_letters/serializers"
  "github.com/beego/beego/v2/client/orm"
)

type UsersController struct {
  BaseController
}

func (c *UsersController) Get() {
  o := orm.NewOrm()

  currentUser := c.currentUser()
  o.LoadRelated(&currentUser, "Language")

  c.setHeaders(currentUser)
  serializer := serializers.User{}
  c.Data["json"] = serializer.Serialize(currentUser)
  c.ServeJSON()
}

func (c *UsersController) Update() {
  o := orm.NewOrm()

  currentUser := c.currentUser()

  var userUpdateForm forms.UserUpdateForm
  c.parseBody(c.Ctx.Input.RequestBody, &userUpdateForm)
  c.processForm(userUpdateForm)
  language := c.Ctx.Input.GetData("language").(models.Language)

  currentUser.Language = &language
  o.Update(&currentUser)

  c.setHeaders(currentUser)
  serializer := serializers.User{}
  c.Data["json"] = serializer.Serialize(currentUser)
  c.ServeJSON()
}
