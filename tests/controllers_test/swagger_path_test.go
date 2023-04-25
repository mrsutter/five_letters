package controllers_test

import (
  tu "five_letters/tests/testutils"
  . "github.com/onsi/ginkgo/v2"
  . "github.com/onsi/gomega"
)

var _ = Describe("SwaggerPath", func() {
  It("returns status 200", func() {
    status, _, _ := tu.HTTPGet(swaggerPath)
    Expect(status).To(Equal(200))
  })
})
