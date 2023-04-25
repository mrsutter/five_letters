package tasks_test

import (
  "context"
  "five_letters/models"
  "five_letters/tasks"
  tf "five_letters/tests/testfactories"
  "github.com/beego/beego/v2/client/orm"
  . "github.com/onsi/ginkgo/v2"
  . "github.com/onsi/gomega"
  "time"
)

var _ = Describe("WasteOutdatedGames", func() {
  var userWithOutdatedGame models.User
  var outdatedGame models.Game

  var userWithFreshGame models.User
  var freshGame models.Game

  BeforeEach(func() {
    o := orm.NewOrm()

    enLang := tf.CreateEnLanguage(true)
    enList := tf.CreateWordList(enLang, 1)
    word := tf.CreateWord("pizza", enList)

    userWithFreshGame = tf.CreateUser("nick", "password", time.Now(), enLang)
    freshGame = tf.CreateGame(userWithFreshGame, word, models.StateActive, 0)

    createdAt := freshGame.CreatedAt.Add((-models.MaxHoursForGame + 1) * time.Hour)
    o.Raw("UPDATE games SET created_at = ? WHERE id = ?",
      createdAt,
      freshGame.Id).
      Exec()

    userWithOutdatedGame = tf.CreateUser("nick2", "password", time.Now(), enLang)
    outdatedGame = tf.CreateGame(userWithOutdatedGame, word, models.StateActive, 0)

    createdAt = outdatedGame.CreatedAt.Add(-models.MaxHoursForGame * time.Hour)
    o.Raw("UPDATE games SET created_at = ? WHERE id = ?",
      createdAt,
      outdatedGame.Id).
      Exec()
  })

  It("wastes only outdated games", func() {
    ctx := context.Background()
    tasks.WasteOutdatedGames(ctx)

    o := orm.NewOrm()

    o.Read(&outdatedGame)
    Expect(outdatedGame.State).To(Equal(models.StateWasted))

    o.Read(&freshGame)
    Expect(freshGame.State).To(Equal(models.StateActive))
  })
})
