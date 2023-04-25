package validators

import (
  "five_letters/models"
  "five_letters/utils"
  "github.com/beego/beego/v2/client/orm"
  "github.com/beego/beego/v2/server/web"
)

type RefreshTokenValidator struct {
  BaseValidator
  Token      *models.RefreshToken
  TokenValue string
}

func (v RefreshTokenValidator) Validate() models.Error {
  o := orm.NewOrm()

  rTokenPubKey, _ := web.AppConfig.String("refreshTokenPublicKey")
  userId, jti, err := utils.ValidateToken(v.TokenValue, rTokenPubKey)
  if err != nil {
    return v.customError()
  }

  user := models.User{Id: int(userId.(float64))}
  err = o.Read(&user)
  if err != nil {
    return v.customError()
  }

  v.Token.Jti = jti.(string)
  v.Token.User = &user

  err = o.Read(v.Token, "Jti", "User")
  if err != nil {
    return v.customError()
  }

  return models.Error{}
}

func (v *RefreshTokenValidator) customError() models.Error {
  eDetails := []models.ErrorItem{
    {
      Field: "refresh_token",
      Code:  v.translateErrorCode("wrong"),
    },
  }

  return models.Error{
    Code:    v.translateErrorCode("inputErrors"),
    Details: eDetails,
  }
}
