package serializers

import (
  "five_letters/models"
)

type User struct {
  Nickname string    `json:"nickname"`
  Language *Language `json:"language"`
}

func (userSerializer *User) Serialize(user models.User) *User {
  userSerializer.Nickname = user.Nickname

  languageSerializer := Language{}
  userSerializer.Language = languageSerializer.Serialize(*user.Language)

  return userSerializer
}
