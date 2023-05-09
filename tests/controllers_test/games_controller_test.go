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
  "path"
  "strconv"
  "strings"
  "time"
)

var _ = Describe("GamesController", func() {
  var user models.User
  var token string
  var expToken string
  var currTime = time.Now()
  var nextGameAvailableAt = time.Now().Add(models.MaxHoursForGame * time.Hour)

  var game models.Game
  var gameWord models.Word
  var state = models.StateActive

  var attemptWords = []string{"index"}
  var attemptsCount = len(attemptWords)

  var enList models.WordList
  var listWords = []string{"pizza", "brain", "cross", "lover", "beach", "bread"}

  var gameOfAnotherUser models.Game

  BeforeEach(func() {
    enLang := tf.CreateEnLanguage(true)
    enList = tf.CreateWordList(enLang, 1)

    gameWord = tf.CreateWord(listWords[0], enList)

    user, token, _ = tf.CreateUserWithTokens("nick", "password", nextGameAvailableAt, enLang)
    expToken, _ = tf.CreateTokens(user, true)

    game = tf.CreateGameWithAttempts(user, gameWord, state, attemptWords...)

    anotherUser := tf.CreateUser("nick2", "password", currTime, enLang)
    gameOfAnotherUser = tf.CreateGameWithAttempts(anotherUser, gameWord, state, attemptWords...)
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

  Describe("GET index", func() {
    DescribeTable("When incorrect token was sent",
      func(tokenType string) {
        var body ts.ErrorSchema

        t := getTokenByType(tokenType)
        status, bodyRaw, _ := tu.HTTPGetWithToken(gamesPath, t)
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
      It("returns status 200, correct data in correct order(created_at) and headers", func() {
        wastedGame := tf.CreateGameWithAttempts(user, gameWord, models.StateWasted, attemptWords...)

        var body ts.GamesSchema

        status, bodyRaw, headers := tu.HTTPGetWithToken(gamesPath, token)
        json.Unmarshal(bodyRaw, &body)

        Expect(status).To(Equal(200))

        Expect(body).To(HaveLen(2))

        firstGame := body[0]
        Expect(firstGame.Id).To(Equal(game.Id))
        Expect(firstGame.State).To(Equal(state))
        Expect(firstGame.AttemptsCount).To(Equal(attemptsCount))
        Expect(firstGame.CreatedAt).To(Equal(utils.FormatTime(game.CreatedAt)))

        secondGame := body[1]
        Expect(secondGame.Id).To(Equal(wastedGame.Id))
        Expect(secondGame.State).To(Equal(wastedGame.State))
        Expect(secondGame.AttemptsCount).To(Equal(wastedGame.AttemptsCount))
        Expect(secondGame.CreatedAt).To(Equal(utils.FormatTime(wastedGame.CreatedAt)))
        nextGameHeaderValue := headers.Get("Next-Game-Available-At")
        Expect(nextGameHeaderValue).To(Equal(utils.FormatTime(nextGameAvailableAt)))
      })
    })
  })

  Describe("GET read", func() {
    var url string

    var _ = BeforeEach(func() {
      url = path.Join(gamesPath, strconv.Itoa(game.Id))
    })

    DescribeTable("When incorrect token was sent",
      func(tokenType string) {
        var body ts.ErrorSchema

        t := getTokenByType(tokenType)
        status, bodyRaw, _ := tu.HTTPGetWithToken(url, t)
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

    DescribeTable("incorrect id was sent",
      func(idType string) {
        var body ts.ErrorSchema
        var id string

        if idType == "unexistingId" {
          id = "unexistingId"
        } else {
          id = strconv.Itoa(gameOfAnotherUser.Id)
        }

        url = path.Join(gamesPath, id)
        status, bodyRaw, _ := tu.HTTPGetWithToken(url, token)
        json.Unmarshal(bodyRaw, &body)

        Expect(status).To(Equal(404))

        Expect(body.Status).To(Equal(404))
        Expect(body.Code).To(Equal("not_found"))
        Expect(body.Message).To(Equal(""))
        Expect(body.Details).To(Equal([]ts.ErrorSchemaItem{}))
      },

      Entry("unexisting id was sent", "unexistingId"),
      Entry("id of game of another user was sent", "gameOfAnotherUserId"),
    )

    Context("When correct id was sent", func() {
      It("returns status 200, correct data and headers", func() {
        var body ts.GameSchema

        status, bodyRaw, headers := tu.HTTPGetWithToken(url, token)
        json.Unmarshal(bodyRaw, &body)

        Expect(status).To(Equal(200))

        Expect(body.Id).To(Equal(game.Id))
        Expect(body.State).To(Equal(game.State))
        Expect(body.AttemptsCount).To(Equal(game.AttemptsCount))
        Expect(body.CreatedAt).To(Equal(utils.FormatTime(game.CreatedAt)))

        Expect(body.Attempts).To(HaveLen(1))

        attempt := body.Attempts[0]
        Expect(attempt.Number).To(Equal(1))
        Expect(attempt.Word).To(Equal(attemptWords[0]))
        Expect(attempt.Result).To(Equal(
          []string{"wrong_place", "absence", "absence", "absence", "absence"},
        ))

        nextGameHeaderValue := headers.Get("Next-Game-Available-At")
        Expect(nextGameHeaderValue).To(Equal(utils.FormatTime(nextGameAvailableAt)))
      })
    })
  })

  Describe("GET active", func() {
    DescribeTable("When incorrect token was sent",
      func(tokenType string) {
        var body ts.ErrorSchema

        t := getTokenByType(tokenType)
        status, bodyRaw, _ := tu.HTTPGetWithToken(gamesActivePath, t)
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

    Context("When there is no active game", func() {
      var _ = BeforeEach(func() {
        game.TransitionTo(models.StateWasted)
        o := orm.NewOrm()
        o.Update(&game, "State")
      })

      It("returns status 404 and correct error", func() {
        var body ts.ErrorSchema

        status, bodyRaw, _ := tu.HTTPGetWithToken(gamesActivePath, token)
        json.Unmarshal(bodyRaw, &body)

        Expect(status).To(Equal(404))

        Expect(body.Status).To(Equal(404))
        Expect(body.Code).To(Equal("not_found"))
        Expect(body.Message).To(Equal(""))
        Expect(body.Details).To(Equal([]ts.ErrorSchemaItem{}))
      })
    })

    Context("When there is active game", func() {
      It("returns status 200, correct data and headers", func() {
        var body ts.GameSchema

        status, bodyRaw, headers := tu.HTTPGetWithToken(gamesActivePath, token)
        json.Unmarshal(bodyRaw, &body)

        Expect(status).To(Equal(200))

        Expect(body.Id).To(Equal(game.Id))
        Expect(body.State).To(Equal(game.State))
        Expect(body.AttemptsCount).To(Equal(game.AttemptsCount))
        Expect(body.CreatedAt).To(Equal(utils.FormatTime(game.CreatedAt)))
        Expect(body.Attempts).To(HaveLen(1))

        attempt := body.Attempts[0]
        Expect(attempt.Number).To(Equal(1))
        Expect(attempt.Word).To(Equal(attemptWords[0]))
        Expect(attempt.Result).To(Equal(
          []string{"wrong_place", "absence", "absence", "absence", "absence"},
        ))

        nextGameHeaderValue := headers.Get("Next-Game-Available-At")
        Expect(nextGameHeaderValue).To(Equal(utils.FormatTime(nextGameAvailableAt)))
      })
    })
  })

  Describe("POST create", func() {
    DescribeTable("When incorrect token was sent",
      func(tokenType string) {
        var body ts.ErrorSchema

        t := getTokenByType(tokenType)
        status, bodyRaw, _ := tu.HTTPPostWithToken(gamesPath, []byte{}, t)
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

    Context("When it's too early to start a new game", func() {
      Context("When NextGameAvailableAt is in the future", func() {
        It("returns status 422 and correct error", func() {
          var body ts.ErrorSchema

          status, bodyRaw, _ := tu.HTTPPostWithToken(gamesPath, []byte{}, token)
          json.Unmarshal(bodyRaw, &body)

          Expect(status).To(Equal(422))

          Expect(body.Status).To(Equal(422))
          Expect(body.Code).To(Equal("too_early"))
          Expect(body.Message).To(Equal(""))
          Expect(body.Details).To(Equal([]ts.ErrorSchemaItem{}))
        })
      })

      Context("When NextGameAvailableAt is not in the future, but there is an active game", func() {
        var _ = BeforeEach(func() {
          user.NextGameAvailableAt = time.Now()

          o := orm.NewOrm()
          o.Update(&user, "NextGameAvailableAt")
        })

        It("returns status 500 and correct error", func() {
          var body ts.ErrorSchema

          status, bodyRaw, _ := tu.HTTPPostWithToken(gamesPath, []byte{}, token)
          json.Unmarshal(bodyRaw, &body)

          Expect(status).To(Equal(500))

          Expect(body.Status).To(Equal(500))
          Expect(body.Code).To(Equal("internal_server_error"))
          Expect(body.Message).To(Equal(""))
          Expect(body.Details).To(Equal([]ts.ErrorSchemaItem{}))
        })
      })
    })

    Context("When user can start a new game", func() {
      var _ = BeforeEach(func() {
        o := orm.NewOrm()

        user.NextGameAvailableAt = time.Now()
        o.Update(&user, "NextGameAvailableAt")

        game.TransitionTo(models.StateWasted)
        o.Update(&game, "State")
      })

      DescribeTable("saves game in db and returns correct data, status and headers",
        func(noNewWords bool) {
          var w models.Word

          if !noNewWords {
            w = tf.CreateWord("lover", enList)
          } else {
            w = gameWord
          }

          var body ts.GameSchema

          status, bodyRaw, headers := tu.HTTPPostWithToken(gamesPath, []byte{}, token)
          json.Unmarshal(bodyRaw, &body)

          o := orm.NewOrm()
          o.LoadRelated(&user, "Games")

          Expect(user.Games).To(HaveLen(2))

          game := user.Games[1]
          Expect(game.State).To(Equal(models.StateActive))
          Expect(game.AttemptsCount).To(Equal(0))
          Expect(game.Word.Id).To(Equal(w.Id))

          Expect(status).To(Equal(201))

          Expect(body.Id).To(Equal(game.Id))
          Expect(body.State).To(Equal(game.State))
          Expect(body.AttemptsCount).To(Equal(game.AttemptsCount))
          Expect(body.CreatedAt).To(Equal(utils.FormatTime(game.CreatedAt)))
          Expect(body.Attempts).To(HaveLen(0))

          nextGameAvailableAt := game.CreatedAt.Add(models.MaxHoursForGame * time.Hour)
          o.Read(&user)
          Expect(user.NextGameAvailableAt.Unix()).To(Equal(nextGameAvailableAt.Unix()))

          nextGameHeaderValue := headers.Get("Next-Game-Available-At")
          Expect(nextGameHeaderValue).To(Equal(utils.FormatTime(nextGameAvailableAt)))
        },

        Entry("there are no new words in word list", true),
        Entry("there are new words in word list", false),
      )
    })
  })

  Describe("Post createAttempt", func() {
    var _ = BeforeEach(func() {
      for _, value := range listWords[1:] {
        tf.CreateWord(value, enList)
      }
    })

    DescribeTable("When incorrect token was sent",
      func(tokenType string) {
        var body ts.ErrorSchema

        t := getTokenByType(tokenType)
        status, bodyRaw, _ := tu.HTTPPostWithToken(attemptsPath, []byte{}, t)

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

    Context("When there is no active game", func() {
      var _ = BeforeEach(func() {
        game.TransitionTo(models.StateWasted)

        o := orm.NewOrm()
        o.Update(&game, "State")
      })

      It("returns status 404 and correct error", func() {
        var body ts.ErrorSchema

        status, bodyRaw, _ := tu.HTTPPostWithToken(attemptsPath, []byte{}, token)
        json.Unmarshal(bodyRaw, &body)

        Expect(status).To(Equal(404))

        Expect(body.Status).To(Equal(404))
        Expect(body.Code).To(Equal("not_found"))
        Expect(body.Message).To(Equal(""))
        Expect(body.Details).To(Equal([]ts.ErrorSchemaItem{}))
      })
    })

    Context("When incorrect data was sent", func() {
      Context("When bad json was sent", func() {
        It("returns status 422 and correct error", func() {
          var body ts.ErrorSchema

          reqBody := `{"wrong json"}`
          status, bodyRaw, _ := tu.HTTPPostWithToken(attemptsPath, []byte(reqBody), token)
          json.Unmarshal(bodyRaw, &body)

          Expect(status).To(Equal(422))

          Expect(body.Status).To(Equal(422))
          Expect(body.Code).To(Equal("bad_json"))
          Expect(body.Message).NotTo(Equal(""))
          Expect(body.Details).To(Equal([]ts.ErrorSchemaItem{}))
        })
      })

      DescribeTable("When incorrect word was sent",
        func(reqBodyType string, detailsCode string) {
          var reqBody string
          switch reqBodyType {
          case "withNoWord":
            reqBody = `{"nickname": "newNickmame"}`
          case "withWordNotFromList":
            reqBody = fmt.Sprintf(`{"word": "%s"}`, "crowd")
          case "withArchivedWord":
            word := models.Word{Name: listWords[1]}
            o := orm.NewOrm()
            o.Read(&word, "Name")
            word.Archived = true
            o.Update(&word, "Archived")
            reqBody = fmt.Sprintf(`{"word": "%s"}`, listWords[1])
          default:
            reqBody = ""
          }

          var body ts.ErrorSchema

          status, bodyRaw, _ := tu.HTTPPostWithToken(attemptsPath, []byte(reqBody), token)
          json.Unmarshal(bodyRaw, &body)

          Expect(status).To(Equal(422))

          Expect(body.Status).To(Equal(422))
          Expect(body.Code).To(Equal("input_errors"))
          Expect(body.Message).To(Equal(""))

          details := body.Details
          Expect(details).To(HaveLen(1))
          Expect(details[0].Field).To(Equal("word"))
          Expect(details[0].Code).To(Equal(detailsCode))
        },

        Entry("no word was sent",
          "withNoWord", "required"),
        Entry("non existent word was sent",
          "withWordNotFromList", "not_found",
        ),
        Entry("archived word was sent",
          "withArchivedWord", "not_found",
        ),
      )
    })
    Context("When correct data was sent", func() {
      DescribeTable("saves attempt in db, returns correct status, body and headers",
        func(
          attemptWord string,
          attemptWordArchived bool,
          lastAttempt bool,
          expectedResult []string,
          expectedState string,
          expectedAttemptNumber int,
        ) {

          o := orm.NewOrm()

          if attemptWordArchived {
            w := models.Word{Name: attemptWord}
            o.Read(&w, "Name")
            w.Archived = true
            o.Update(&w)
          }

          if lastAttempt {
            for i, value := range listWords[2:] {
              tf.CreateAttempt(i+2, value, &game)
            }
          }

          var body ts.GameSchema

          reqBody := fmt.Sprintf(`{"word": "%s"}`, attemptWord)
          status, bodyRaw, headers := tu.HTTPPostWithToken(attemptsPath, []byte(reqBody), token)
          json.Unmarshal(bodyRaw, &body)

          o.LoadRelated(&user, "Games")

          game = models.Game{Id: game.Id}
          o.Read(&game)
          o.LoadRelated(&game, "Attempts")

          Expect(game.State).To(Equal(expectedState))
          Expect(game.AttemptsCount).To(Equal(expectedAttemptNumber))

          expectedResultJson, _ := json.Marshal(expectedResult)
          expectedResultJsonString := strings.ReplaceAll(string(expectedResultJson), ",", ", ")

          newAttempt := game.Attempts[expectedAttemptNumber-1]

          Expect(newAttempt.Number).To(Equal(expectedAttemptNumber))
          Expect(newAttempt.Word).To(Equal(attemptWord))
          Expect(newAttempt.Result).To(Equal(expectedResultJsonString))

          Expect(status).To(Equal(201))

          Expect(body.Id).To(Equal(game.Id))
          Expect(body.State).To(Equal(game.State))
          Expect(body.AttemptsCount).To(Equal(game.AttemptsCount))
          Expect(body.CreatedAt).To(Equal(utils.FormatTime(game.CreatedAt)))
          Expect(body.Attempts).To(HaveLen(expectedAttemptNumber))

          Expect(body.Attempts[expectedAttemptNumber-1].Number).To(Equal(expectedAttemptNumber))
          Expect(body.Attempts[expectedAttemptNumber-1].Word).To(Equal(attemptWord))
          Expect(body.Attempts[expectedAttemptNumber-1].Result).To(Equal(expectedResult))

          nextGameHeaderValue := headers.Get("Next-Game-Available-At")
          Expect(nextGameHeaderValue).To(Equal(utils.FormatTime(nextGameAvailableAt)))
        },

        Entry("word is not correct and it is not the last attempt",
          listWords[1],
          false,
          false,
          []string{"absence", "absence", "wrong_place", "wrong_place", "absence"},
          models.StateActive,
          2,
        ),
        Entry("word is not correct and it is the last attempt",
          listWords[1],
          false,
          true,
          []string{"absence", "absence", "wrong_place", "wrong_place", "absence"},
          models.StateWasted,
          6,
        ),
        Entry("word is correct and it is not the last attempt",
          listWords[0],
          false,
          false,
          []string{"match", "match", "match", "match", "match"},
          models.StateWon,
          2,
        ),
        Entry("word is correct and it is last attempt",
          listWords[0],
          false,
          true,
          []string{"match", "match", "match", "match", "match"},
          models.StateWon,
          6,
        ),
        Entry("word is correct but archived",
          listWords[0],
          true,
          false,
          []string{"match", "match", "match", "match", "match"},
          models.StateWon,
          2,
        ),
      )
    })
  })
})
