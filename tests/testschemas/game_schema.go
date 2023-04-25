package testschemas

type GameSchema struct {
  Id            int              `json:"id"`
  State         string           `json:"state"`
  AttemptsCount int              `json:"attempts_count"`
  CreatedAt     string           `json:"created_at"`
  Attempts      []*AttemptSchema `json:"attempts"`
}

type GameShortSchema struct {
  Id            int    `json:"id"`
  State         string `json:"state"`
  AttemptsCount int    `json:"attempts_count"`
  CreatedAt     string `json:"created_at"`
}

type GamesSchema []GameShortSchema
