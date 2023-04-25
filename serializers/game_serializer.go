package serializers

import (
  "five_letters/models"
  "five_letters/utils"
  "sort"
)

type Game struct {
  Id            int        `json:"id"`
  State         string     `json:"state"`
  AttemptsCount int        `json:"attempts_count"`
  CreatedAt     string     `json:"created_at"`
  Attempts      []*Attempt `json:"attempts"`
}

func (gameSerializer *Game) Serialize(game models.Game) *Game {
  gameSerializer.Id = game.Id
  gameSerializer.State = game.State
  gameSerializer.AttemptsCount = game.AttemptsCount
  gameSerializer.CreatedAt = utils.FormatTime(game.CreatedAt)

  if len(game.Attempts) == 0 {
    gameSerializer.Attempts = []*Attempt{}
  } else {
    sortAttempts(game.Attempts)

    for _, attempt := range game.Attempts {
      attemptSerializer := Attempt{}
      gameSerializer.Attempts = append(
        gameSerializer.Attempts,
        attemptSerializer.Serialize(*attempt),
      )
    }
  }

  return gameSerializer
}

func sortAttempts(attempts []*models.Attempt) {
  sort.Slice(attempts, func(i, j int) bool {
    return attempts[i].Number < attempts[j].Number
  })
}

type GameShort struct {
  Id            int    `json:"id"`
  State         string `json:"state"`
  AttemptsCount int    `json:"attempts_count"`
  CreatedAt     string `json:"created_at"`
}

func (gameShortSerializer *GameShort) Serialize(game models.Game) *GameShort {
  gameShortSerializer.Id = game.Id
  gameShortSerializer.State = game.State
  gameShortSerializer.AttemptsCount = game.AttemptsCount
  gameShortSerializer.CreatedAt = utils.FormatTime(game.CreatedAt)

  return gameShortSerializer
}

type Games []*GameShort

func (gamesSerializer Games) Serialize(games []*models.Game) Games {
  for _, game := range games {
    gameShortSerializer := GameShort{}
    gamesSerializer = append(gamesSerializer, gameShortSerializer.Serialize(*game))
  }
  return gamesSerializer
}
