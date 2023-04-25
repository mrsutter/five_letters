package testschemas

type AttemptSchema struct {
  Number int      `json:"number"`
  Word   string   `json:"word"`
  Result []string `json:"result"`
}
