package models_test

import (
  "five_letters/tests/testutils"
  . "github.com/onsi/ginkgo/v2"
  . "github.com/onsi/gomega"
  "testing"
)

func TestModels(t *testing.T) {
  RegisterFailHandler(Fail)
  RunSpecs(t, "Models Suite")
}

var _ = BeforeSuite(func() {
  testutils.InitBeego()
})
