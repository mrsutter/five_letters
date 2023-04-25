package tasks

import (
  "context"
  "five_letters/models"
  "github.com/beego/beego/v2/client/orm"
  "github.com/beego/beego/v2/core/logs"
  "time"
)

func WasteOutdatedGames(ctx context.Context) error {
  o := orm.NewOrm()
  _, err := o.QueryTable("games").
    Filter("state", models.StateActive).
    Filter("created_at__lte", time.Now().Add(-models.MaxHoursForGame*time.Hour)).
    Update(orm.Params{"state": models.StateWasted})

  if err == nil {
    logs.Info("Outdated games were wasted")
  } else {
    logs.Critical("Error while wasting outdated games: %s", err)
  }

  return err
}
