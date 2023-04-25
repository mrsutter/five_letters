package middlewares

import (
  "five_letters/models"
  "five_letters/utils"
  "github.com/beego/beego/v2/client/orm"
  "github.com/beego/beego/v2/server/web"
  "github.com/beego/beego/v2/server/web/context"
  "strings"
)

func Authenticate(ctx *context.Context) {
  var token string

  header := ctx.Request.Header.Get("Authorization")
  token = strings.ReplaceAll(header, "Bearer ", "")

  if token == "" {
    ctx.Abort(401, "401")
  }

  aTokenPubKey, _ := web.AppConfig.String("accessTokenPublicKey")

  userId, jti, err := utils.ValidateToken(token, aTokenPubKey)
  if err != nil {
    ctx.Abort(401, "401")
  }

  o := orm.NewOrm()

  user := models.User{Id: int(userId.(float64))}
  err = o.Read(&user)
  if err != nil {
    ctx.Abort(401, "401")
  }

  accessToken := models.AccessToken{Jti: jti.(string), User: &user}
  err = o.Read(&accessToken, "Jti", "User")

  if err != nil {
    ctx.Abort(401, "401")
  }

  ctx.Input.SetData("currentUser", user)
  ctx.Input.SetData("accessToken", accessToken)
}
