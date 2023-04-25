package tasks_test

import (
  "five_letters/tests/testutils"
  . "github.com/onsi/ginkgo/v2"
  . "github.com/onsi/gomega"
  "gopkg.in/khaiql/dbcleaner.v2"

  "testing"
)

func TestTasks(t *testing.T) {
  RegisterFailHandler(Fail)
  RunSpecs(t, "Tasks Suite")
}

var Cleaner = dbcleaner.New()

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
