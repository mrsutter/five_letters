package services

import (
  "five_letters/models"
  "five_letters/utils"
  "github.com/beego/beego/v2/client/orm"
  "github.com/beego/beego/v2/server/web"
)

func CreateTokens(user models.User) (string, string, error) {
  o := orm.NewOrm()

  aTokenTTL, _ := web.AppConfig.Int("accessTokenTTLSec")
  aTokenPrKey, _ := web.AppConfig.String("accessTokenPrivateKey")

  rTokenTTL, _ := web.AppConfig.Int("refreshTokenTTLSec")
  rTokenPrKey, _ := web.AppConfig.String("refreshTokenPrivateKey")

  aTokenValue, aTokenJti, aTokenExp, _ := utils.GenerateToken(aTokenTTL, user.Id, aTokenPrKey)
  rTokenValue, rTokenJti, rTokenExp, _ := utils.GenerateToken(rTokenTTL, user.Id, rTokenPrKey)

  rToken := models.RefreshToken{
    User:      &user,
    Jti:       rTokenJti,
    ExpiredAt: rTokenExp,
  }

  _, err := o.Insert(&rToken)
  if err != nil {
    return "", "", err
  }

  aToken := models.AccessToken{
    User:         &user,
    Jti:          aTokenJti,
    ExpiredAt:    aTokenExp,
    RefreshToken: &rToken,
  }

  _, err = o.Insert(&aToken)
  if err != nil {
    return "", "", err
  }

  return aTokenValue, rTokenValue, nil
}
