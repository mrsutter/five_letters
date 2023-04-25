package tasks

import (
  "context"
  "github.com/beego/beego/v2/client/orm"
  "github.com/beego/beego/v2/core/logs"
  "time"
)

func DeleteOutdatedTokens(ctx context.Context) error {
  o := orm.NewOrm()

  _, err := o.QueryTable("refresh_tokens").
    Filter("expired_at__lt", time.Now()).
    Delete()
  if err != nil {
    logErrorOnDeleteOutdatedTokens(err)
    return err
  }

  _, err = o.QueryTable("access_tokens").
    Filter("expired_at__lt", time.Now()).
    Delete()
  if err != nil {
    logErrorOnDeleteOutdatedTokens(err)
    return err
  }

  logs.Info("Outdated refresh and access tokens were deleted")
  return nil
}

func logErrorOnDeleteOutdatedTokens(err error) {
  logs.Critical("Error while deleting outdated refresh and access tokens: %s", err)
}
