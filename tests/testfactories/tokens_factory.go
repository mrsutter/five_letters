package testfactories

import (
  "five_letters/models"
  tu "five_letters/tests/testutils"
  "five_letters/utils"
  "github.com/beego/beego/v2/client/orm"
)

func CreateTokens(user models.User, expired bool) (string, string) {
  o := orm.NewOrm()

  aTokenTTL := tu.GetConfigIntValue("accessTokenTTLSec")
  aTokenPrKey := tu.GetConfigStringValue("accessTokenPrivateKey")

  rTokenTTL := tu.GetConfigIntValue("refreshTokenTTLSec")
  rTokenPrKey := tu.GetConfigStringValue("refreshTokenPrivateKey")

  if expired {
    aTokenTTL = -aTokenTTL
    rTokenTTL = -rTokenTTL
  }

  aTokenValue, aTokenJti, aTokenExp, _ := utils.GenerateToken(aTokenTTL, user.Id, aTokenPrKey)
  rTokenValue, rTokenJti, rTokenExp, _ := utils.GenerateToken(rTokenTTL, user.Id, rTokenPrKey)

  rToken := models.RefreshToken{
    User:      &user,
    Jti:       rTokenJti,
    ExpiredAt: rTokenExp,
  }

  o.Insert(&rToken)

  aToken := models.AccessToken{
    User:         &user,
    Jti:          aTokenJti,
    ExpiredAt:    aTokenExp,
    RefreshToken: &rToken,
  }

  o.Insert(&aToken)
  return aTokenValue, rTokenValue
}
