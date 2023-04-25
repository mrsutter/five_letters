package controllers_test

import (
  "encoding/json"
  "five_letters/models"
  tf "five_letters/tests/testfactories"
  ts "five_letters/tests/testschemas"
  tu "five_letters/tests/testutils"
  . "github.com/onsi/ginkgo/v2"
  . "github.com/onsi/gomega"
)

var _ = Describe("LanguagesController", func() {
  var enLang models.Language
  var ruLang models.Language

  BeforeEach(func() {
    enLang = tf.CreateEnLanguage(true)
  })

  Context("When all languages are available", func() {
    BeforeEach(func() {
      ruLang = tf.CreateRuLanguage(true)
    })

    It("returns correct status and sorted languages by created_at", func() {
      var body ts.LanguagesSchema

      status, bodyRaw, _ := tu.HTTPGet(languagesPath)
      json.Unmarshal(bodyRaw, &body)

      Expect(status).To(Equal(200))

      Expect(len(body)).To(Equal(2))

      enL := body[0]
      Expect(enL.Id).To(Equal(enLang.Id))
      Expect(enL.Slug).To(Equal(enLang.Slug))
      Expect(enL.Name).To(Equal(enLang.Name))

      ruL := body[1]
      Expect(ruL.Id).To(Equal(ruLang.Id))
      Expect(ruL.Slug).To(Equal(ruLang.Slug))
      Expect(ruL.Name).To(Equal(ruLang.Name))
    })
  })

  Context("When not all languages are available", func() {
    BeforeEach(func() {
      ruLang = tf.CreateRuLanguage(false)
    })

    It("returns correct status and sorted available languages only", func() {
      var body ts.LanguagesSchema

      status, bodyRaw, _ := tu.HTTPGet(languagesPath)
      json.Unmarshal(bodyRaw, &body)

      Expect(status).To(Equal(200))

      Expect(len(body)).To(Equal(1))

      lang := body[0]
      Expect(lang.Id).To(Equal(enLang.Id))
      Expect(lang.Slug).To(Equal(enLang.Slug))
      Expect(lang.Name).To(Equal(enLang.Name))
    })
  })
})
