package controllers

import (
  "five_letters/forms"
  "five_letters/models"
  "five_letters/serializers"
  "five_letters/services"
  "five_letters/utils"
  "github.com/beego/beego/v2/client/orm"
  "time"
)

type AuthController struct {
  BaseController
}

func (c *AuthController) Register() {
  var userSignUp forms.UserSignUpForm

  c.parseBody(c.Ctx.Input.RequestBody, &userSignUp)
  c.processForm(userSignUp)

  o := orm.NewOrm()

  language := c.Ctx.Input.GetData("language").(models.Language)
  user := models.User{
    Nickname:            userSignUp.Nickname,
    Password:            utils.HashPassword(userSignUp.Password),
    Language:            &language,
    NextGameAvailableAt: time.Now(),
  }

  _, err := o.Insert(&user)
  if err != nil {
    c.Abort("500")
  }

  c.Ctx.Output.SetStatus(201)
  serializer := serializers.User{}
  c.Data["json"] = serializer.Serialize(user)
  c.ServeJSON()
}

func (c *AuthController) Login() {
  var userSignInForm forms.UserSignInForm

  c.parseBody(c.Ctx.Input.RequestBody, &userSignInForm)
  c.processForm(userSignInForm)

  o := orm.NewOrm()

  user := models.User{Nickname: userSignInForm.Nickname}
  o.Read(&user, "Nickname")

  accessT, refreshT, err := services.CreateTokens(user)
  if err != nil {
    c.Abort("500")
  }

  c.Data["json"] = serializers.Tokens{AccessToken: accessT, RefreshToken: refreshT}
  c.ServeJSON()
}

func (c *AuthController) Logout() {
  accessToken := c.Ctx.Input.GetData("accessToken").(models.AccessToken)

  err := services.DeleteTokensByAccessToken(&accessToken)
  if err != nil {
    c.Abort("500")
  }

  c.Ctx.Output.SetStatus(204)
}

func (c *AuthController) Refresh() {
  var refreshTokenForm forms.RefreshTokenForm

  c.parseBody(c.Ctx.Input.RequestBody, &refreshTokenForm)
  c.processForm(refreshTokenForm)

  refreshToken := c.Ctx.Input.GetData("refreshToken").(models.RefreshToken)

  accessT, refreshT, err := services.CreateTokens(*refreshToken.User)
  if err != nil {
    c.Abort("500")
  }

  err = services.DeleteTokensByRefreshToken(&refreshToken)
  if err != nil {
    c.Abort("500")
  }

  c.Data["json"] = serializers.Tokens{AccessToken: accessT, RefreshToken: refreshT}
  c.ServeJSON()
}
