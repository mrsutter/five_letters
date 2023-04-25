package models_test

import (
  "five_letters/models"
  "fmt"
  . "github.com/onsi/ginkgo/v2"
  . "github.com/onsi/gomega"
)

var _ = Describe("Game", func() {
  var unknownState = "unknown"

  Describe("CanTransitionTo", func() {
    DescribeTable("CanTransitionTo Table",
      func(state string, newState string, result bool) {
        game := models.Game{State: state}
        Expect(game.CanTransitionTo(newState)).To(Equal(result))
      },

      Entry("active->active", models.StateActive, models.StateActive, false),
      Entry("active->wasted", models.StateActive, models.StateWasted, true),
      Entry("active->won", models.StateActive, models.StateWon, true),
      Entry("active->unknown", models.StateActive, unknownState, false),

      Entry("won->active", models.StateWon, models.StateActive, false),
      Entry("won->wasted", models.StateWon, models.StateWasted, false),
      Entry("won->won", models.StateWon, models.StateWon, false),
      Entry("won->unknown", models.StateWon, unknownState, false),

      Entry("wasted->active", models.StateWasted, models.StateActive, false),
      Entry("wasted->wasted", models.StateWasted, models.StateWasted, false),
      Entry("wasted->won", models.StateWasted, models.StateWon, false),
      Entry("wasted->unknown", models.StateWasted, unknownState, false),
    )
  })

  Describe("TransitionTo", func() {
    DescribeTable("TransitionTo Table",
      func(state string, newState string, shouldReturnError bool) {
        game := models.Game{State: state}

        if shouldReturnError {
          result := fmt.Errorf("invalid transition from %s to %s", state, newState)
          Expect(game.TransitionTo(newState)).To(Equal(result))
          Expect(game.State).To(Equal(state))
        } else {
          Expect(game.TransitionTo(newState)).To(BeNil())
          Expect(game.State).To(Equal(newState))
        }
      },

      Entry("active->active", models.StateActive, models.StateActive, true),
      Entry("active->wasted", models.StateActive, models.StateWasted, false),
      Entry("active->wasted", models.StateActive, models.StateWon, false),
      Entry("active->wasted", models.StateActive, unknownState, true),

      Entry("won->active", models.StateWon, models.StateActive, true),
      Entry("won->wasted", models.StateWon, models.StateWasted, true),
      Entry("won->wasted", models.StateWon, models.StateWon, true),
      Entry("won->wasted", models.StateWon, unknownState, true),

      Entry("wasted->active", models.StateWasted, models.StateActive, true),
      Entry("wasted->wasted", models.StateWasted, models.StateWasted, true),
      Entry("wasted->wasted", models.StateWasted, models.StateWon, true),
      Entry("wasted->wasted", models.StateWasted, unknownState, true),
    )
  })

  Describe("SetValuesAfterAttempt", func() {
    var active = models.StateActive
    var won = models.StateWon
    var wasted = models.StateWasted

    It("does nothing if attempt is from another game", func() {
      gameWord := models.Word{Name: "pizza"}
      attemptsCount := 1

      game := models.Game{
        State:         active,
        Word:          &gameWord,
        AttemptsCount: attemptsCount,
      }
      anotherGame := models.Game{}

      attempt := models.Attempt{Game: &anotherGame}

      game.SetValuesAfterAttempt(attempt)
      Expect(game.AttemptsCount).To(Equal(attemptsCount))
      Expect(game.State).To(Equal(active))
    })

    DescribeTable("SetValuesAfterAttempt",
      func(
        initState string,
        finalState string,
        puzzledWord string,
        attemptWord string,
        initAttemptsCount int,
        finalAttemptsCount int,
      ) {

        gameWord := models.Word{Name: puzzledWord}
        game := models.Game{
          State:         initState,
          Word:          &gameWord,
          AttemptsCount: initAttemptsCount,
        }

        attempt := models.Attempt{
          Game: &game,
          Word: attemptWord,
        }

        game.SetValuesAfterAttempt(attempt)
        Expect(game.AttemptsCount).To(Equal(finalAttemptsCount))
        Expect(game.State).To(Equal(finalState))
      },

      Entry("1st wrong attempt on active", active, active, "pizza", "index", 0, 1),
      Entry("1st successful attempt on active", active, won, "pizza", "pizza", 0, 1),
      Entry("1st wrong attempt on wasted", wasted, wasted, "pizza", "index", 0, 0),
      Entry("1st successful attempt on wasted", wasted, wasted, "pizza", "pizza", 0, 0),
      Entry("1st wrong attempt on won", won, won, "pizza", "index", 0, 0),
      Entry("1st successful attempt on won", won, won, "pizza", "pizza", 0, 0),

      Entry("2nd wrong attempt on active", active, active, "pizza", "index", 1, 2),
      Entry("2nd successful attempt on active", active, won, "pizza", "pizza", 1, 2),
      Entry("2nd wrong attempt on wasted", wasted, wasted, "pizza", "index", 1, 1),
      Entry("2nd successful attempt on wasted", wasted, wasted, "pizza", "pizza", 1, 1),
      Entry("2nd wrong attempt on won", won, won, "pizza", "index", 1, 1),
      Entry("2nd successful attempt on won", won, won, "pizza", "pizza", 1, 1),

      Entry("3rd wrong attempt on active", active, active, "pizza", "index", 2, 3),
      Entry("3rd successful attempt on active", active, won, "pizza", "pizza", 2, 3),
      Entry("3rd wrong attempt on wasted", wasted, wasted, "pizza", "index", 2, 2),
      Entry("3rd successful attempt on wasted", wasted, wasted, "pizza", "pizza", 2, 2),
      Entry("3rd wrong attempt on won", won, won, "pizza", "index", 2, 2),
      Entry("3rd successful attempt on won", won, won, "pizza", "pizza", 2, 2),

      Entry("4th wrong attempt on active", active, active, "pizza", "index", 3, 4),
      Entry("4th successful attempt on active", active, won, "pizza", "pizza", 3, 4),
      Entry("4th wrong attempt on wasted", wasted, wasted, "pizza", "index", 3, 3),
      Entry("4th successful attempt on wasted", wasted, wasted, "pizza", "pizza", 3, 3),
      Entry("4th wrong attempt on won", won, won, "pizza", "index", 3, 3),
      Entry("4th successful attempt on won", won, won, "pizza", "pizza", 3, 3),

      Entry("5th wrong attempt on active", active, active, "pizza", "index", 4, 5),
      Entry("5th successful attempt on active", active, won, "pizza", "pizza", 4, 5),
      Entry("5th wrong attempt on wasted", wasted, wasted, "pizza", "index", 4, 4),
      Entry("5th successful attempt on wasted", wasted, wasted, "pizza", "pizza", 4, 4),
      Entry("5th wrong attempt on won", won, won, "pizza", "index", 4, 4),
      Entry("5th successful attempt on won", won, won, "pizza", "pizza", 4, 4),

      Entry("6th wrong attempt on active", active, wasted, "pizza", "index", 5, 6),
      Entry("6th successful attempt on active", active, won, "pizza", "pizza", 5, 6),
      Entry("6th wrong attempt on wasted", wasted, wasted, "pizza", "index", 5, 5),
      Entry("6th successful attempt on wasted", wasted, wasted, "pizza", "pizza", 5, 5),
      Entry("6th wrong attempt on won", won, won, "pizza", "index", 5, 5),
      Entry("6th successful attempt on won", won, won, "pizza", "pizza", 5, 5),

      Entry("7th wrong attempt on active", active, active, "pizza", "index", 6, 6),
      Entry("7th successful attempt on active", active, active, "pizza", "pizza", 6, 6),
      Entry("7th wrong attempt on wasted", wasted, wasted, "pizza", "index", 6, 6),
      Entry("7th successful attempt on wasted", wasted, wasted, "pizza", "pizza", 6, 6),
      Entry("7th wrong attempt on won", won, won, "pizza", "index", 6, 6),
      Entry("7th successful attempt on won", won, won, "pizza", "pizza", 6, 6),
    )
  })
})
