package serializers

import (
  "five_letters/models"
)

type Language struct {
  Id   int    `json:"id"`
  Slug string `json:"slug"`
  Name string `json:"name"`
}

func (languageSerializer *Language) Serialize(language models.Language) *Language {
  languageSerializer.Id = language.Id
  languageSerializer.Slug = language.Slug
  languageSerializer.Name = language.Name

  return languageSerializer
}

type Languages []*Language

func (languagesSerializer Languages) Serialize(languages []*models.Language) Languages {
  for _, lang := range languages {
    langSerializer := Language{}
    languagesSerializer = append(languagesSerializer, langSerializer.Serialize(*lang))
  }
  return languagesSerializer
}
