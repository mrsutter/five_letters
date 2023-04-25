package services

import (
  "five_letters/models"
  "github.com/beego/beego/v2/client/orm"
)

func DeleteTokensByAccessToken(accessToken *models.AccessToken) error {
  err := DeleteTokensByRefreshToken(accessToken.RefreshToken)
  return err
}

func DeleteTokensByRefreshToken(refreshToken *models.RefreshToken) error {
  o := orm.NewOrm()
  _, err := o.Delete(refreshToken)
  return err
}
