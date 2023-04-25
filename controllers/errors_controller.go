package controllers

import (
  "five_letters/models"
  "five_letters/serializers"
  "strconv"
)

type ErrorsController struct {
  BaseController
}

func (c *ErrorsController) Error404() {
  status := 404

  e := models.Error{
    Status: status,
    Code:   c.translateErrorCode(strconv.Itoa(status)),
  }

  c.Ctx.Output.Status = status

  serializer := serializers.Error{}
  c.Data["json"] = serializer.Serialize(e)
  c.ServeJSON()
}

func (c *ErrorsController) Error401() {
  status := 401

  e := models.Error{
    Status: status,
    Code:   c.translateErrorCode(strconv.Itoa(status)),
  }

  c.Ctx.Output.Status = status

  serializer := serializers.Error{}
  c.Data["json"] = serializer.Serialize(e)
  c.ServeJSON()
}

func (c *ErrorsController) Error500() {
  status := 500

  e := models.Error{
    Status: status,
    Code:   c.translateErrorCode(strconv.Itoa(status)),
  }

  c.Ctx.Output.Status = status

  serializer := serializers.Error{}
  c.Data["json"] = serializer.Serialize(e)
  c.ServeJSON()
}

func (c *ErrorsController) Error422() {
  status := 422

  e := c.Ctx.Input.GetData("error").(models.Error)
  e.Status = status

  c.Ctx.Output.Status = status

  serializer := serializers.Error{}
  c.Data["json"] = serializer.Serialize(e)
  c.ServeJSON()
}
