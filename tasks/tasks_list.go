package tasks

import (
  "context"
)

type Task struct {
  Id    string
  RunAt string
  Func  func(ctx context.Context) error
}

func InitTasks() (tasks []*Task) {
  return []*Task{
    &Task{
      Id:    "SeedLanguagesAndListsAndWords",
      RunAt: "0 0 * * * *",
      Func:  SeedLanguagesAndListsAndWords,
    },
    &Task{
      Id:    "wasteOutdatedGames",
      RunAt: "*/30 * * * * *",
      Func:  WasteOutdatedGames,
    },
    &Task{
      Id:    "deleteOutdatedTokens",
      RunAt: "0 30 * * * *",
      Func:  DeleteOutdatedTokens,
    },
  }
}
