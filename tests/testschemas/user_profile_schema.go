package testschemas

type UserProfileSchema struct {
  Nickname string         `json:"nickname"`
  Language LanguageSchema `json:"language"`
}
