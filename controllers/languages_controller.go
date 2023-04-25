package controllers

import (
  "five_letters/models"
  "five_letters/serializers"
  "github.com/beego/beego/v2/client/orm"
)

type LanguagesController struct {
  BaseController
}

func (c *LanguagesController) Index() {
  var languages []*models.Language

  o := orm.NewOrm()
  o.QueryTable("languages").
    Filter("Available", true).
    OrderBy("created_at").
    All(&languages)

  serializer := serializers.Languages{}
  c.Data["json"] = serializer.Serialize(languages)
  c.ServeJSON()
}
