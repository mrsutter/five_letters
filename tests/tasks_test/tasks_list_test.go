package tasks_test

import (
  t "five_letters/tasks"
  . "github.com/onsi/ginkgo/v2"
  . "github.com/onsi/gomega"
  "reflect"
)

var _ = Describe("InitTasks", func() {
  It("returns correct task list", func() {
    tasks := t.InitTasks()

    Expect(len(tasks)).To(Equal(3))

    Expect(tasks[0].Id).To(Equal("SeedLanguagesAndListsAndWords"))
    Expect(tasks[0].RunAt).To(Equal("0 0 * * * *"))
    Expect(reflect.TypeOf(tasks[0].Func)).To(Equal(
      reflect.TypeOf(t.SeedLanguagesAndListsAndWords),
    ))

    Expect(tasks[1].Id).To(Equal("wasteOutdatedGames"))
    Expect(tasks[1].RunAt).To(Equal("*/30 * * * * *"))
    Expect(reflect.TypeOf(tasks[1].Func)).To(Equal(
      reflect.TypeOf(t.WasteOutdatedGames),
    ))

    Expect(tasks[2].Id).To(Equal("deleteOutdatedTokens"))
    Expect(tasks[2].RunAt).To(Equal("0 30 * * * *"))
    Expect(reflect.TypeOf(tasks[2].Func)).To(Equal(
      reflect.TypeOf(t.DeleteOutdatedTokens),
    ))
  })
})
