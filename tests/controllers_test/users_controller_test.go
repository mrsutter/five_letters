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
  "time"
)

var _ = Describe("UsersController", func() {
  var user models.User
  var enLang models.Language

  var token string
  var expToken string

  var currTime = time.Now()

  BeforeEach(func() {
    enLang = tf.CreateEnLanguage(true)
    user, token, _ = tf.CreateUserWithTokens("nick", "password", currTime, enLang)
    expToken, _ = tf.CreateTokens(user, true)
  })

  var getTokenByType = func(tokenType string) string {
    var t string
    switch tokenType {
    case "empty":
      t = ""
    case "wrong":
      t = "wrong"
    case "expired":
      t = expToken
    default:
      t = token
    }
    return t
  }

  Describe("GET profile", func() {
    DescribeTable("When incorrect token was sent",
      func(tokenType string) {
        var body ts.ErrorSchema

        t := getTokenByType(tokenType)
        status, bodyRaw, _ := tu.HTTPGetWithToken(profilePath, t)
        json.Unmarshal(bodyRaw, &body)

        Expect(status).To(Equal(401))

        Expect(body.Status).To(Equal(401))
        Expect(body.Code).To(Equal("unathorized"))
        Expect(body.Message).To(Equal(""))
        Expect(body.Details).To(Equal([]ts.ErrorSchemaItem{}))
      },

      Entry("no access token was sent", "empty"),
      Entry("wrong token was sent", "wrong"),
      Entry("expired token was sent", "expired"),
    )

    Context("When correct token was sent", func() {
      It("returns status 200, correct data and headers", func() {
        var body ts.UserProfileSchema

        status, bodyRaw, headers := tu.HTTPGetWithToken(profilePath, token)
        json.Unmarshal(bodyRaw, &body)

        Expect(status).To(Equal(200))

        Expect(body.Nickname).To(Equal(user.Nickname))

        lang := body.Language
        Expect(lang.Id).To(Equal(enLang.Id))
        Expect(lang.Slug).To(Equal(enLang.Slug))
        Expect(lang.Name).To(Equal(enLang.Name))

        value := headers.Get("Next-Game-Available-At")
        Expect(value).To(Equal(utils.FormatTime(currTime)))
      })
    })
  })

  Describe("PUT/PATCH profile", func() {
    DescribeTable("When incorrect token was sent",
      func(tokenType string) {
        var body ts.ErrorSchema

        t := getTokenByType(tokenType)
        status, bodyRaw, _ := tu.HTTPPutWithToken(profilePath, []byte{}, t)
        json.Unmarshal(bodyRaw, &body)

        Expect(status).To(Equal(401))

        Expect(body.Status).To(Equal(401))
        Expect(body.Code).To(Equal("unathorized"))
        Expect(body.Message).To(Equal(""))
        Expect(body.Details).To(Equal([]ts.ErrorSchemaItem{}))
      },

      Entry("no access token was sent", "empty"),
      Entry("wrong token was sent", "wrong"),
      Entry("expired token was sent", "expired"),
    )

    Context("When incorrect data was sent", func() {
      Context("When bad json was sent", func() {
        It("returns status 422 and correct error", func() {
          var body ts.ErrorSchema

          reqBody := `{"wrong json"}`
          status, bodyRaw, _ := tu.HTTPPutWithToken(profilePath, []byte(reqBody), token)
          json.Unmarshal(bodyRaw, &body)

          Expect(status).To(Equal(422))

          Expect(body.Status).To(Equal(422))
          Expect(body.Code).To(Equal("bad_json"))
          Expect(body.Message).NotTo(Equal(""))
          Expect(body.Details).To(Equal([]ts.ErrorSchemaItem{}))
        })
      })

      DescribeTable("When wrong language_id was sent",
        func(reqBodyType string, detailsCode string) {

          var reqBody string
          switch reqBodyType {
          case "withNoLanguageId":
            reqBody = `{"nickname": "newNickmame"}`
          case "withWrongLanguageId":
            reqBody = `{"language_id": -2}`
          case "withNotAvailableLanguageId":
            ruLang := tf.CreateRuLanguage(false)
            reqBody = fmt.Sprintf(`{"language_id": %d}`, ruLang.Id)
          default:
            reqBody = ""
          }

          var body ts.ErrorSchema

          t := getTokenByType("default")
          status, bodyRaw, _ := tu.HTTPPutWithToken(profilePath, []byte(reqBody), t)
          json.Unmarshal(bodyRaw, &body)

          Expect(status).To(Equal(422))

          Expect(body.Status).To(Equal(422))
          Expect(body.Code).To(Equal("input_errors"))
          Expect(body.Message).To(Equal(""))

          details := body.Details
          Expect(details).To(HaveLen(1))
          Expect(details[0].Field).To(Equal("language_id"))
          Expect(details[0].Code).To(Equal(detailsCode))

          o := orm.NewOrm()
          o.Read(&user)
          Expect(user.Language.Id).To(Equal(enLang.Id))
        },

        Entry("no language_id was sent",
          "withNoLanguageId", "required"),
        Entry("non existent language_id was sent",
          "withWrongLanguageId", "not_found",
        ),
        Entry("not available language_id was sent",
          "withNotAvailableLanguageId", "not_found",
        ),
      )

      DescribeTable("When correct data was sent",
        func(oldLanguage bool) {
          var sentLanguage models.Language

          if oldLanguage {
            sentLanguage = enLang
          } else {
            sentLanguage = tf.CreateRuLanguage(true)
          }

          var body ts.UserProfileSchema

          reqBody := fmt.Sprintf(`{"language_id": %d}`, sentLanguage.Id)
          status, bodyRaw, headers := tu.HTTPPutWithToken(profilePath, []byte(reqBody), token)
          json.Unmarshal(bodyRaw, &body)

          Expect(status).To(Equal(200))

          Expect(body.Nickname).To(Equal(user.Nickname))

          lang := body.Language
          Expect(lang.Id).To(Equal(sentLanguage.Id))
          Expect(lang.Slug).To(Equal(sentLanguage.Slug))
          Expect(lang.Name).To(Equal(sentLanguage.Name))

          value := headers.Get("Next-Game-Available-At")
          Expect(value).To(Equal(utils.FormatTime(currTime)))

          o := orm.NewOrm()
          o.Read(&user)
          Expect(user.Language.Id).To(Equal(sentLanguage.Id))
        },

        Entry("when old language sent", true),
        Entry("when new language sent", false),
      )
    })
  })
})
