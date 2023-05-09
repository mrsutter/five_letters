package controllers_test

import (
  "encoding/json"
  "five_letters/models"
  tf "five_letters/tests/testfactories"
  ts "five_letters/tests/testschemas"
  tu "five_letters/tests/testutils"
  "five_letters/utils"
  "fmt"
  "github.com/beego/beego/v2/client/orm"
  . "github.com/onsi/ginkgo/v2"
  . "github.com/onsi/gomega"
  "golang.org/x/exp/slices"
  "time"
)

var _ = Describe("AuthController", func() {
  Describe("POST logout", func() {
    var user models.User
    var token string

    BeforeEach(func() {
      enLang := tf.CreateEnLanguage(true)
      user, token, _ = tf.CreateUserWithTokens("nick", "password", time.Now(), enLang)
    })

    DescribeTable("When incorrect token was sent",
      func(tokenType string) {
        var t string
        switch tokenType {
        case "empty":
          t = ""
        case "wrong":
          t = "wrong"
        case "expired":
          t, _ = tf.CreateTokens(user, true)
        }

        var body ts.ErrorSchema

        status, bodyRaw, _ := tu.HTTPPostWithToken(logoutPath, []byte{}, t)
        json.Unmarshal(bodyRaw, &body)

        Expect(status).To(Equal(401))

        Expect(body.Status).To(Equal(401))
        Expect(body.Code).To(Equal("unauthorized"))
        Expect(body.Message).To(Equal(""))
        Expect(body.Details).To(Equal([]ts.ErrorSchemaItem{}))
      },

      Entry("no access token was sent", "empty"),
      Entry("wrong token was sent", "wrong"),
      Entry("expired token was sent", "expired"),
    )

    Context("When correct token was sent", func() {
      It("returns status 204, empty data and deletes tokens from db", func() {
        var accessTokens []models.AccessToken
        var refreshTokens []models.RefreshToken

        status, bodyRaw, headers := tu.HTTPPostWithToken(logoutPath, []byte{}, token)

        Expect(status).To(Equal(204))
        Expect(len(bodyRaw)).To(Equal(0))

        value := headers.Get("Next-Game-Available-At")
        Expect(value).To(Equal(""))

        o := orm.NewOrm()
        o.QueryTable("access_tokens").All(&accessTokens)
        Expect(len(accessTokens)).To(Equal(0))
        o.QueryTable("refresh_tokens").All(&refreshTokens)
        Expect(len(refreshTokens)).To(Equal(0))
      })
    })
  })

  Describe("POST refresh", func() {
    var user models.User
    var refrToken string

    BeforeEach(func() {
      enLang := tf.CreateEnLanguage(true)
      user, _, refrToken = tf.CreateUserWithTokens("nick", "password", time.Now(), enLang)
    })

    Context("When incorrect data was sent", func() {
      Context("When bad json was sent", func() {
        It("returns status 422 and correct error", func() {
          var body ts.ErrorSchema

          reqBody := `{"wrong json"}`
          status, bodyRaw, _ := tu.HTTPPost(refreshPath, []byte(reqBody))
          json.Unmarshal(bodyRaw, &body)

          Expect(status).To(Equal(422))

          Expect(body.Status).To(Equal(422))
          Expect(body.Code).To(Equal("bad_json"))
          Expect(body.Message).NotTo(Equal(""))
          Expect(body.Details).To(Equal([]ts.ErrorSchemaItem{}))
        })
      })

      DescribeTable("When incorrect token was sent",
        func(bodyType string, errCode string) {
          var body ts.ErrorSchema
          var reqBody string

          o := orm.NewOrm()

          switch bodyType {
          case "tokenMissed":
            reqBody = `{"nickname": "newNickmame"}`
          case "tokenUnexisting":
            reqBody = `{"refresh_token": "refresh_token"}`
          case "tokenExpired":
            _, expiredToken := tf.CreateTokens(user, true)
            reqBody = fmt.Sprintf(`{"refresh_token": "%s"}`, expiredToken)
          case "tokenDeleted":
            o.QueryTable("refresh_tokens").
              Filter("user_id", user.Id).
              Delete()
            reqBody = fmt.Sprintf(`{"refresh_token": "%s"}`, refrToken)
          }

          status, bodyRaw, _ := tu.HTTPPost(refreshPath, []byte(reqBody))
          json.Unmarshal(bodyRaw, &body)

          Expect(status).To(Equal(422))

          Expect(body.Status).To(Equal(422))
          Expect(body.Code).To(Equal("input_errors"))
          Expect(body.Message).To(Equal(""))

          details := body.Details
          Expect(details).To(HaveLen(1))
          Expect(details[0].Field).To(Equal("refresh_token"))
          Expect(details[0].Code).To(Equal(errCode))
        },

        Entry("refresh token was missed", "tokenMissed", "required"),
        Entry("refresh token does not exist", "tokenUnexisting", "wrong"),
        Entry("refresh token was expired", "tokenExpired", "wrong"),
        Entry("refresh token was deleted from db", "tokenDeleted", "wrong"),
      )
    })

    Context("When correct data was sent", func() {
      It("returns status 200, correct data and headers and saves values in db", func() {
        var body ts.UserTokensSchema

        reqBody := fmt.Sprintf(`{"refresh_token": "%s"}`, refrToken)
        status, bodyRaw, headers := tu.HTTPPost(refreshPath, []byte(reqBody))
        json.Unmarshal(bodyRaw, &body)

        Expect(status).To(Equal(200))

        Expect(body.AccessToken).NotTo(Equal(""))
        Expect(body.RefreshToken).NotTo(Equal(""))

        value := headers.Get("Next-Game-Available-At")
        Expect(value).To(Equal(""))

        o := orm.NewOrm()

        var accessTokens []models.AccessToken
        o.QueryTable("access_tokens").All(&accessTokens)
        Expect(len(accessTokens)).To(Equal(1))

        var refreshTokens []models.RefreshToken
        o.QueryTable("refresh_tokens").All(&refreshTokens)
        Expect(len(refreshTokens)).To(Equal(1))

        refreshTokenKey := tu.GetConfigStringValue("refreshTokenPublicKey")
        _, refreshTokenJti, _ := utils.ValidateToken(body.RefreshToken, refreshTokenKey)
        refreshTokenIdx := slices.IndexFunc(refreshTokens, func(t models.RefreshToken) bool {
          return t.Jti == refreshTokenJti
        })
        refreshToken := refreshTokens[refreshTokenIdx]
        Expect(refreshToken.User.Id).To(Equal(user.Id))

        accessTokenKey := tu.GetConfigStringValue("accessTokenPublicKey")
        _, accessTokenJti, _ := utils.ValidateToken(body.AccessToken, accessTokenKey)
        accessTokenIdx := slices.IndexFunc(accessTokens, func(t models.AccessToken) bool {
          return t.Jti == accessTokenJti
        })
        accessToken := accessTokens[accessTokenIdx]
        Expect(accessToken.User.Id).To(Equal(user.Id))
        Expect(accessToken.RefreshToken.Id).To(Equal(refreshToken.Id))
      })
    })
  })
  Describe("POST login", func() {
    var user models.User
    var nick = "nick"
    var pass = "password"

    BeforeEach(func() {
      enLang := tf.CreateEnLanguage(true)
      user = tf.CreateUser(nick, pass, time.Now(), enLang)
    })

    Context("When incorrect data was sent", func() {
      Context("When bad json was sent", func() {
        It("returns status 422 and correct error", func() {
          var body ts.ErrorSchema

          reqBody := `{"wrong json"}`
          status, bodyRaw, _ := tu.HTTPPost(loginPath, []byte(reqBody))
          json.Unmarshal(bodyRaw, &body)

          Expect(status).To(Equal(422))

          Expect(body.Status).To(Equal(422))
          Expect(body.Code).To(Equal("bad_json"))
          Expect(body.Message).NotTo(Equal(""))
          Expect(body.Details).To(Equal([]ts.ErrorSchemaItem{}))
        })
      })

      DescribeTable("When nickname and password are incorrect together",
        func(reqBody string, errNicknameCode string, errPasswordCode string) {
          var body ts.ErrorSchema

          status, bodyRaw, _ := tu.HTTPPost(loginPath, []byte(reqBody))
          json.Unmarshal(bodyRaw, &body)

          Expect(status).To(Equal(422))

          Expect(body.Status).To(Equal(422))
          Expect(body.Code).To(Equal("input_errors"))
          Expect(body.Message).To(Equal(""))

          Expect(len(body.Details)).To(Equal(2))

          nickNameError := body.FindItem("nickname")
          Expect(nickNameError.Code).To(Equal(errNicknameCode))

          passwordError := body.FindItem("password")
          Expect(passwordError.Code).To(Equal(errPasswordCode))
        },

        Entry("when nickname and pass were missed",
          `{"someKey": "value"}`,
          "required", "required"),
        Entry("when wrong password was sent",
          fmt.Sprintf(`{"nickname": "%s", "password": "wrong_pass"}`, nick),
          "no_user_with_such_credentials",
          "no_user_with_such_credentials"),
        Entry("when unexisting nickname was sent",
          fmt.Sprintf(`{"nickname": "unexistent", "password": "%s"}`, pass),
          "no_user_with_such_credentials",
          "no_user_with_such_credentials"),
      )

      DescribeTable("When one field is incorrect",
        func(reqBody string, errField string, errCode string) {
          var body ts.ErrorSchema

          status, bodyRaw, _ := tu.HTTPPost(loginPath, []byte(reqBody))
          json.Unmarshal(bodyRaw, &body)

          Expect(status).To(Equal(422))

          Expect(body.Status).To(Equal(422))
          Expect(body.Code).To(Equal("input_errors"))
          Expect(body.Message).To(Equal(""))

          Expect(len(body.Details)).To(Equal(1))

          errItem := body.FindItem(errField)
          Expect(errItem.Code).To(Equal(errCode))
        },

        Entry("when nickname was missed",
          fmt.Sprintf(`{"password": "%s"}`, pass),
          "nickname", "required"),
        Entry("when nickname was incorrect",
          fmt.Sprintf(`{"nickname": "$авар", "password": "%s"}`, pass),
          "nickname", "wrong"),
        Entry("when password was missed",
          fmt.Sprintf(`{"nickname": "%s"}`, nick),
          "password", "required"),
      )
    })

    Context("When correct data was sent", func() {
      It("returns status 200, correct data and headers and saves values in db", func() {
        var body ts.UserTokensSchema

        reqBody := fmt.Sprintf(`{"nickname": "%s", "password": "%s"}`, nick, pass)
        status, bodyRaw, headers := tu.HTTPPost(loginPath, []byte(reqBody))
        json.Unmarshal(bodyRaw, &body)

        Expect(status).To(Equal(200))

        Expect(body.AccessToken).NotTo(Equal(""))
        Expect(body.RefreshToken).NotTo(Equal(""))

        value := headers.Get("Next-Game-Available-At")
        Expect(value).To(Equal(""))

        o := orm.NewOrm()

        var accessTokens []models.AccessToken
        o.QueryTable("access_tokens").All(&accessTokens)
        Expect(len(accessTokens)).To(Equal(1))

        var refreshTokens []models.RefreshToken
        o.QueryTable("refresh_tokens").All(&refreshTokens)
        Expect(len(refreshTokens)).To(Equal(1))

        refreshTokenKey := tu.GetConfigStringValue("refreshTokenPublicKey")
        _, refreshTokenJti, _ := utils.ValidateToken(body.RefreshToken, refreshTokenKey)
        refreshTokenIdx := slices.IndexFunc(refreshTokens, func(t models.RefreshToken) bool {
          return t.Jti == refreshTokenJti
        })
        refreshToken := refreshTokens[refreshTokenIdx]
        Expect(refreshToken.User.Id).To(Equal(user.Id))

        accessTokenKey := tu.GetConfigStringValue("accessTokenPublicKey")
        _, accessTokenJti, _ := utils.ValidateToken(body.AccessToken, accessTokenKey)
        accessTokenIdx := slices.IndexFunc(accessTokens, func(t models.AccessToken) bool {
          return t.Jti == accessTokenJti
        })
        accessToken := accessTokens[accessTokenIdx]
        Expect(accessToken.User.Id).To(Equal(user.Id))
        Expect(accessToken.RefreshToken.Id).To(Equal(refreshToken.Id))
      })
    })
  })

  Describe("POST register", func() {
    var enLang models.Language
    var ruLang models.Language

    var nick = "nick"
    var pass = "password1"

    BeforeEach(func() {
      enLang = tf.CreateEnLanguage(true)
      ruLang = tf.CreateRuLanguage(false)
    })

    Context("When incorrect data was sent", func() {
      Context("When bad json was sent", func() {
        It("returns status 422 and correct error", func() {
          var body ts.ErrorSchema

          reqBody := `{"wrong json"}`
          status, bodyRaw, _ := tu.HTTPPost(registerPath, []byte(reqBody))
          json.Unmarshal(bodyRaw, &body)

          Expect(status).To(Equal(422))

          Expect(body.Status).To(Equal(422))
          Expect(body.Code).To(Equal("bad_json"))
          Expect(body.Message).NotTo(Equal(""))
          Expect(body.Details).To(Equal([]ts.ErrorSchemaItem{}))
        })
      })

      Context("When all params were missed", func() {
        It("returns status 422 and correct error", func() {
          var body ts.ErrorSchema

          reqBody := fmt.Sprintf(`{"someKey": "value"}`)
          status, bodyRaw, _ := tu.HTTPPost(registerPath, []byte(reqBody))
          json.Unmarshal(bodyRaw, &body)

          Expect(status).To(Equal(422))

          Expect(body.Status).To(Equal(422))
          Expect(body.Code).To(Equal("input_errors"))
          Expect(body.Message).To(Equal(""))

          Expect(len(body.Details)).To(Equal(4))

          nickNameError := body.FindItem("nickname")
          Expect(nickNameError.Code).To(Equal("required"))
          passwordError := body.FindItem("password")
          Expect(passwordError.Code).To(Equal("required"))
          passwordConfirmationError := body.FindItem("password_confirmation")
          Expect(passwordConfirmationError.Code).To(Equal("required"))
          languageError := body.FindItem("language_id")
          Expect(languageError.Code).To(Equal("required"))
        })
      })

      DescribeTable("When incorrect param was sent",
        func(bodyType string, errField string, errDetailsCode string) {
          var body ts.ErrorSchema
          var reqBody string

          switch bodyType {
          case "nickNameMissed":
            reqBody = fmt.Sprintf(
              `{"password": "%s", "password_confirmation": "%s",
              "language_id": %d}`,
              pass, pass, enLang.Id)
          case "passwordMissed":
            reqBody = fmt.Sprintf(
              `{"nickname": "%s", "password_confirmation": "%s",
              "language_id": %d}`,
              nick, pass, enLang.Id)
          case "passwordConfirmationMissed":
            reqBody = fmt.Sprintf(
              `{"nickname": "%s", "password": "%s", "language_id": %d}`,
              nick, pass, enLang.Id)
          case "languageIdMissed":
            reqBody = fmt.Sprintf(
              `{"nickname": "%s", "password": "%s",
              "password_confirmation": "%s"}`,
              nick, pass, pass)
          case "passwordWasTooShort":
            reqBody = fmt.Sprintf(
              `{"nickname": "%s", "password": "123",
              "password_confirmation": "%s", "language_id": %d}`,
              nick, pass, enLang.Id)
          case "passwordConfirmationTooShort":
            reqBody = fmt.Sprintf(
              `{"nickname": "%s", "password": "%s",
              "password_confirmation": "123", "language_id": %d}`,
              nick, pass, enLang.Id)
          case "nicknameIncorrect":
            reqBody = fmt.Sprintf(
              `{"nickname": "$атов", "password": "%s",
              "password_confirmation": "%s", "language_id": %d}`,
              pass, pass, enLang.Id)
          case "nicknameAlreadyTaken":
            tf.CreateUser(nick, pass, time.Now(), enLang)
            reqBody = fmt.Sprintf(
              `{"nickname": "%s", "password": "%s",
              "password_confirmation": "%s", "language_id": %d}`,
              nick, pass, pass, enLang.Id)
          case "languageIdUnexisted":
            reqBody = fmt.Sprintf(
              `{"nickname": "%s", "password": "%s",
              "password_confirmation": "%s", "language_id": %d}`,
              nick, pass, pass, -2)
          case "languageIdUnavailable":
            reqBody = fmt.Sprintf(
              `{"nickname": "%s", "password": "%s",
              "password_confirmation": "%s", "language_id": %d}`,
              nick, pass, pass, ruLang.Id)
          }

          status, bodyRaw, _ := tu.HTTPPost(registerPath, []byte(reqBody))
          json.Unmarshal(bodyRaw, &body)

          Expect(status).To(Equal(422))

          Expect(body.Status).To(Equal(422))
          Expect(body.Code).To(Equal("input_errors"))
          Expect(body.Message).To(Equal(""))

          Expect(len(body.Details)).To(Equal(1))

          errItem := body.FindItem(errField)
          Expect(errItem.Code).To(Equal(errDetailsCode))
        },

        Entry("when nickname was missed", "nickNameMissed",
          "nickname", "required"),
        Entry("when password was missed", "passwordMissed",
          "password", "required"),
        Entry("when password was missed", "passwordConfirmationMissed",
          "password_confirmation", "required"),
        Entry("when language_id was missed", "languageIdMissed",
          "language_id", "required"),
        Entry("when password was too short", "passwordWasTooShort",
          "password", "too_short"),
        Entry("when password_confirmation was too short", "passwordConfirmationTooShort",
          "password_confirmation", "too_short"),
        Entry("when nickname was incorrect", "nicknameIncorrect",
          "nickname", "wrong"),
        Entry("when nickname was already taken", "nicknameAlreadyTaken",
          "nickname", "already_taken"),
        Entry("when languageId is unexisted", "languageIdUnexisted",
          "language_id", "not_found"),
        Entry("when languageId is unavailable", "languageIdUnavailable",
          "language_id", "not_found"),
      )

      Context("When password and password confirmation do not match", func() {
        It("returns status 422 and correct error", func() {
          var body ts.ErrorSchema

          reqBody := fmt.Sprintf(
            `{"nickname": "%s", "password": "%s", "password_confirmation": "confirmation", "language_id": %d}`,
            nick, pass, enLang.Id)
          status, bodyRaw, _ := tu.HTTPPost(registerPath, []byte(reqBody))
          json.Unmarshal(bodyRaw, &body)

          Expect(status).To(Equal(422))

          Expect(body.Status).To(Equal(422))
          Expect(body.Code).To(Equal("input_errors"))
          Expect(body.Message).To(Equal(""))

          Expect(len(body.Details)).To(Equal(2))

          passwordError := body.FindItem("password")
          Expect(passwordError.Code).To(Equal("passwords_are_not_equal"))
          passwordConfirmationError := body.FindItem("password_confirmation")
          Expect(passwordConfirmationError.Code).To(Equal("passwords_are_not_equal"))
        })
      })
    })

    Context("When correct data was sent", func() {
      It("returns status 201, correct data and headers and saves values in db", func() {
        var body ts.UserProfileSchema

        reqBody := fmt.Sprintf(
          `{"nickname": "%s", "password": "%s", "password_confirmation": "%s", "language_id": %d}`,
          nick, pass, pass, enLang.Id)
        status, bodyRaw, headers := tu.HTTPPost(registerPath, []byte(reqBody))
        json.Unmarshal(bodyRaw, &body)

        Expect(status).To(Equal(201))

        Expect(body.Nickname).To(Equal(nick))

        lang := body.Language
        Expect(lang.Id).To(Equal(enLang.Id))
        Expect(lang.Slug).To(Equal(enLang.Slug))
        Expect(lang.Name).To(Equal(enLang.Name))

        value := headers.Get("Next-Game-Available-At")
        Expect(value).To(Equal(""))

        o := orm.NewOrm()
        u := models.User{Nickname: nick}
        o.Read(&u, "Nickname")
        Expect(u.Language.Id).To(Equal(enLang.Id))
      })
    })
  })
})
