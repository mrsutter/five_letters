package serializers

import (
  "encoding/json"
  "five_letters/models"
)

type Attempt struct {
  Number int      `json:"number"`
  Word   string   `json:"word"`
  Result []string `json:"result"`
}

func (attemptSerializer *Attempt) Serialize(attempt models.Attempt) *Attempt {
  attemptSerializer.Number = attempt.Number
  attemptSerializer.Word = attempt.Word
  json.Unmarshal([]byte(attempt.Result), &attemptSerializer.Result)

  return attemptSerializer
}
