package controllers_test

import (
  "encoding/json"
  ts "five_letters/tests/testschemas"
  tu "five_letters/tests/testutils"
  . "github.com/onsi/ginkgo/v2"
  . "github.com/onsi/gomega"
)

var _ = Describe("UnexistingPath", func() {
  DescribeTable("UnexistingPathTable",
    func(method string, path string, reqBody []byte) {
      var body ts.ErrorSchema

      status, bodyRaw, _ := tu.HttpRequest(method, path, reqBody, "")
      json.Unmarshal(bodyRaw, &body)

      Expect(status).To(Equal(404))

      Expect(body.Status).To(Equal(404))
      Expect(body.Code).To(Equal("not_found"))
      Expect(body.Message).To(Equal(""))
      Expect(body.Details).To(Equal([]ts.ErrorSchemaItem{}))
    },

    Entry("Get", "GET", unExistingPath, nil),
    Entry("Post", "POST", unExistingPath, []byte{}),
    Entry("Put", "PUT", unExistingPath, []byte{}),
    Entry("Patch", "PATCH", unExistingPath, []byte{}),
    Entry("Delete", "DELETE", unExistingPath, nil),
  )
})
