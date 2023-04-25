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

var _ = Describe("SeedLanguagesAndListsAndWords", func() {
  Context("When it's first launch", func() {
    It("creates languages, words and word lists correctly", func() {
      ctx := context.Background()
      tasks.SeedLanguagesAndListsAndWords(ctx)

      var languages []*models.Language
      o := orm.NewOrm()
      o.QueryTable("languages").
        Filter("Available", true).
        All(&languages)

      enLang := languages[0]
      Expect(enLang.Slug).To(Equal("en"))
      Expect(enLang.Name).To(Equal("English"))
      Expect(enLang.LettersList).To(Equal("^[a-z]+$"))
      Expect(enLang.Available).To(Equal(true))

      enWordList := models.WordList{Language: enLang}
      o.Read(&enWordList, "Language")
      Expect(enWordList.Version).To(Equal(2))

      o.LoadRelated(&enWordList, "Words")
      Expect(enWordList.Words).To(HaveLen(2))
      firstEnWord := enWordList.Words[0]
      Expect(firstEnWord.Name).To(Equal("lover"))
      secondEnWord := enWordList.Words[1]
      Expect(secondEnWord.Name).To(Equal("brain"))

      ruLang := languages[1]
      Expect(ruLang.Slug).To(Equal("ru"))
      Expect(ruLang.Name).To(Equal("Русский"))
      Expect(ruLang.LettersList).To(Equal("^[а-я]+$"))
      Expect(ruLang.Available).To(Equal(true))

      ruWordList := models.WordList{Language: ruLang}
      o.Read(&ruWordList, "Language")
      Expect(ruWordList.Version).To(Equal(1))

      o.LoadRelated(&ruWordList, "Words")
      Expect(ruWordList.Words).To(HaveLen(2))
      firstRuWord := ruWordList.Words[0]
      Expect(firstRuWord.Name).To(Equal("пирог"))
      secondRuWord := ruWordList.Words[1]
      Expect(secondRuWord.Name).To(Equal("замок"))
    })
  })

  Context("When it's not first launch", func() {
    It("doesn't create words if word_list already have highest version", func() {
      enLang := tf.CreateEnLanguage(true)
      enList := tf.CreateWordList(enLang, 2)

      ruLang := tf.CreateRuLanguage(true)
      ruList := tf.CreateWordList(ruLang, 1)

      o := orm.NewOrm()

      ctx := context.Background()
      tasks.SeedLanguagesAndListsAndWords(ctx)

      o.Read(&ruList)
      o.LoadRelated(&enList, "Words")
      Expect(enList.Words).To(HaveLen(0))

      o.Read(&enList)
      o.LoadRelated(&ruList, "Words")
      Expect(enList.Words).To(HaveLen(0))
    })
  })

  Context("When language was archived", func() {
    It("updates it correctly", func() {
      enLang := tf.CreateEnLanguage(false)

      ctx := context.Background()
      tasks.SeedLanguagesAndListsAndWords(ctx)

      o := orm.NewOrm()

      o.Read(&enLang)
      Expect(enLang.Available).To(Equal(true))
    })
  })

  Context("When word was archived in previous version", func() {
    It("makes it available", func() {
      enLang := tf.CreateEnLanguage(true)
      enList := tf.CreateWordList(enLang, 1)
      word := tf.CreateWord("lover", enList)

      o := orm.NewOrm()
      word.Archived = true
      o.Update(&word)

      ctx := context.Background()
      tasks.SeedLanguagesAndListsAndWords(ctx)

      o.Read(&word)
      Expect(word.Archived).To(Equal(false))
    })
  })

  Context("When word was actual previously but now it is absent", func() {
    It("makes it archived", func() {
      enLang := tf.CreateEnLanguage(true)
      enList := tf.CreateWordList(enLang, 1)

      o := orm.NewOrm()

      word := tf.CreateWord("pizza", enList)
      updatedAt := time.Now().Add(-1 * time.Second)
      o.Raw("UPDATE words SET updated_at = ? WHERE id = ?", updatedAt, word.Id).Exec()

      ctx := context.Background()
      tasks.SeedLanguagesAndListsAndWords(ctx)

      o.Read(&word)
      Expect(word.Archived).To(Equal(true))
    })
  })
})
