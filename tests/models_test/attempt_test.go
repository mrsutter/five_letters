package models_test

import (
  "five_letters/models"
  . "github.com/onsi/ginkgo/v2"
  . "github.com/onsi/gomega"
)

var _ = Describe("Attempt", func() {
  Describe("CalcResult", func() {
    DescribeTable("CalcResult Table",
      func(word string, puzzledWord string, result []string) {
        gameWord := models.Word{Name: puzzledWord}
        game := models.Game{Word: &gameWord}

        attempt := models.Attempt{Word: word, Game: &game}

        Expect(attempt.CalcResult()).To(Equal(result))
      },
      Entry("index-pizza", "index", "pizza",
        []string{"wrong_place", "absence", "absence", "absence", "absence"},
      ),
      Entry("index-upset", "index", "upset",
        []string{"absence", "absence", "absence", "match", "absence"},
      ),
      Entry("agent-brand", "agent", "brand",
        []string{"wrong_place", "absence", "absence", "match", "absence"},
      ),
      Entry("lease-money", "lease", "money",
        []string{"absence", "wrong_place", "absence", "absence", "absence"},
      ),
      Entry("offer-order", "offer", "order",
        []string{"match", "absence", "absence", "match", "match"},
      ),
      Entry("proxy-quota", "proxy", "quota",
        []string{"absence", "absence", "match", "absence", "absence"},
      ),
      Entry("stock-value", "stock", "value",
        []string{"absence", "absence", "absence", "absence", "absence"},
      ),
      Entry("fraud-hedge", "fraud", "hedge",
        []string{"absence", "absence", "absence", "absence", "wrong_place"},
      ),
      Entry("truck-abort", "truck", "abort",
        []string{"wrong_place", "wrong_place", "absence", "absence", "absence"},
      ),
      Entry("pizza-pizza", "pizza", "pizza",
        []string{"match", "match", "match", "match", "match"},
      ),
      Entry("кошка-толпа", "кошка", "толпа",
        []string{"absence", "match", "absence", "absence", "match"},
      ),
      Entry("пицца-зелье", "пицца", "зелье",
        []string{"absence", "absence", "absence", "absence", "absence"},
      ),
      Entry("кулич-кашпо", "кулич", "кашпо",
        []string{"match", "absence", "absence", "absence", "absence"},
      ),
      Entry("штора-арбуз", "штора", "арбуз",
        []string{"absence", "absence", "absence", "wrong_place", "wrong_place"},
      ),
      Entry("плита-пушка", "плита", "пушка",
        []string{"match", "absence", "absence", "absence", "match"},
      ),
      Entry("мороз-волна", "мороз", "волна",
        []string{"absence", "match", "absence", "absence", "absence"},
      ),
      Entry("линза-левша", "линза", "левша",
        []string{"match", "absence", "absence", "absence", "match"},
      ),
      Entry("пилка-пилот", "пилка", "пилот",
        []string{"match", "match", "match", "absence", "absence"},
      ),
      Entry("товар-актор", "товар", "актор",
        []string{"wrong_place", "wrong_place", "absence", "wrong_place", "match"},
      ),
      Entry("пицца-пицца", "пицца", "пицца",
        []string{"match", "match", "match", "match", "match"},
      ),
    )
  })

  Describe("Successful", func() {
    It("returns correct bool value", func() {
      gameWord := models.Word{Name: "pizza"}
      game := models.Game{Word: &gameWord}

      firstAttempt := models.Attempt{Word: "index", Game: &game}
      secondAttempt := models.Attempt{Word: "pizza", Game: &game}

      Expect(firstAttempt.Successful()).To(Equal(false))
      Expect(secondAttempt.Successful()).To(Equal(true))
    })
  })
})
