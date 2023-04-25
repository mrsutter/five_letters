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

var _ = Describe("DeleteOutdatedTokens", func() {
  var activeAccessToken *models.AccessToken
  var activeRefreshToken *models.RefreshToken

  BeforeEach(func() {
    o := orm.NewOrm()

    enLang := tf.CreateEnLanguage(true)

    user, _, _ := tf.CreateUserWithTokens("nick", "password", time.Now(), enLang)
    o.LoadRelated(&user, "AccessTokens")
    o.LoadRelated(&user, "RefreshTokens")

    activeAccessToken = user.AccessTokens[0]
    activeRefreshToken = user.RefreshTokens[0]

    tf.CreateTokens(user, true)
  })

  It("deletes only outdated tokens", func() {
    o := orm.NewOrm()

    ctx := context.Background()
    tasks.DeleteOutdatedTokens(ctx)

    var accessTokens []models.AccessToken
    o.QueryTable("access_tokens").All(&accessTokens)
    Expect(len(accessTokens)).To(Equal(1))

    var refreshTokens []models.RefreshToken
    o.QueryTable("refresh_tokens").All(&refreshTokens)
    Expect(len(refreshTokens)).To(Equal(1))

    Expect(accessTokens[0].Id).To(Equal(activeAccessToken.Id))
    Expect(refreshTokens[0].Id).To(Equal(activeRefreshToken.Id))
  })
})
