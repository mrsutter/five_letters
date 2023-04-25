package testschemas

type LanguageSchema struct {
  Id   int    `json:"id"`
  Slug string `json:"slug"`
  Name string `json:"name"`
}

type LanguagesSchema []LanguageSchema
