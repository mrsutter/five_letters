package controllers

import (
  "five_letters/forms"
  "five_letters/models"
  "five_letters/serializers"
  "five_letters/services"
  "github.com/beego/beego/v2/client/orm"
  "time"
)

type GamesController struct {
  BaseController
}

func (c *GamesController) Index() {
  o := orm.NewOrm()

  currentUser := c.currentUser()

  var games []*models.Game

  o.QueryTable("games").
    Filter("user", currentUser).
    OrderBy("created_at").
    All(&games)

  c.setHeaders(currentUser)
  serializer := serializers.Games{}
  c.Data["json"] = serializer.Serialize(games)
  c.ServeJSON()
}

func (c *GamesController) Read() {
  o := orm.NewOrm()

  currentUser := c.currentUser()

  gameId, err := c.GetInt(":id")
  if err != nil {
    c.Abort("404")
  }

  game := models.Game{Id: gameId, User: &currentUser}
  err = o.Read(&game, "Id", "User")
  if err != nil {
    c.Abort("404")
  }
  o.LoadRelated(&game, "Attempts")

  c.setHeaders(currentUser)
  serializer := serializers.Game{}
  c.Data["json"] = serializer.Serialize(game)
  c.ServeJSON()
}

func (c *GamesController) Active() {
  o := orm.NewOrm()

  currentUser := c.currentUser()

  game := models.Game{State: models.StateActive, User: &currentUser}
  err := o.Read(&game, "State", "User")
  if err != nil {
    c.Abort("404")
  }

  o.LoadRelated(&game, "Attempts")

  c.setHeaders(currentUser)
  serializer := serializers.Game{}
  c.Data["json"] = serializer.Serialize(game)
  c.ServeJSON()
}

func (c *GamesController) Create() {
  currentUser := c.currentUser()

  if currentUser.NextGameAvailableAt.Unix() > time.Now().Unix() {
    c.Abort422(models.Error{
      Code: c.translateErrorCode("tooEarly"),
    })
  }

  game, err := services.CreateGame(currentUser)
  if err != nil {
    c.Abort("500")
  }

  c.Ctx.Output.SetStatus(201)
  c.setHeaders(*game.User)

  serializer := serializers.Game{}
  c.Data["json"] = serializer.Serialize(game)
  c.ServeJSON()
}

func (c *GamesController) CreateAttempt() {
  o := orm.NewOrm()

  currentUser := c.currentUser()

  game := models.Game{State: models.StateActive, User: &currentUser}
  err := o.Read(&game, "State", "User")
  if err != nil {
    c.Abort("404")
  }
  o.LoadRelated(&game, "Word")
  c.Ctx.Input.SetData("game", game)

  var attemptForm forms.AttemptForm
  c.parseBody(c.Ctx.Input.RequestBody, &attemptForm)
  c.processForm(attemptForm)

  _, err = services.CreateAttempt(&game, attemptForm.Word)
  if err != nil {
    c.Abort("500")
  }

  o.LoadRelated(&game, "Attempts")

  c.Ctx.Output.SetStatus(201)
  c.setHeaders(currentUser)

  serializer := serializers.Game{}
  c.Data["json"] = serializer.Serialize(game)
  c.ServeJSON()
}
