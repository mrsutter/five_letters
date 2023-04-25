package controllers_test

import (
  "five_letters/tests/testutils"
  . "github.com/onsi/ginkgo/v2"
  . "github.com/onsi/gomega"
  "gopkg.in/khaiql/dbcleaner.v2"
  "testing"
)

const (
  languagesPath = "/api/v1/languages"

  profilePath = "/api/v1/profile"

  gamesPath       = "/api/v1/games"
  gamesActivePath = "/api/v1/games/active"
  attemptsPath    = "/api/v1/games/active/attempts"

  refreshPath  = "/api/v1/auth/refresh"
  logoutPath   = "/api/v1/auth/logout"
  loginPath    = "/api/v1/auth/login"
  registerPath = "/api/v1/auth/register"

  unExistingPath = "/api/v1/somePath"

  swaggerPath = "/swagger/"
)

var Cleaner = dbcleaner.New()

func TestControllers(t *testing.T) {
  RegisterFailHandler(Fail)
  RunSpecs(t, "Controllers Suite")
}

var _ = BeforeSuite(func() {
  testutils.InitBeego()
  testutils.InitDbCleaner(Cleaner)
})

var _ = BeforeEach(func() {
  Cleaner.Acquire(
    "languages",
    "users",
    "access_tokens",
    "refresh_tokens",
    "games",
    "attempts",
  )
})

var _ = AfterEach(func() {
  Cleaner.Clean(
    "languages",
    "users",
    "access_tokens",
    "refresh_tokens",
    "games",
    "attempts",
  )
})
