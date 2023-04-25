package tasks

import (
  "context"
  "five_letters/models"
  "five_letters/utils"
  "github.com/beego/beego/v2/client/orm"
  "github.com/beego/beego/v2/core/logs"
  "github.com/beego/beego/v2/server/web"
  "regexp"
  "strconv"
  "strings"
  "sync"
  "time"
  "unicode/utf8"
)

func SeedLanguagesAndListsAndWords(ctx context.Context) error {
  languages := ReadOrCreateLanguages()
  lists := readOrCreateWordLists(languages)

  var seedWordGroup sync.WaitGroup
  for _, list := range lists {
    list := list
    seedWordGroup.Add(1)
    go func() {
      defer seedWordGroup.Done()
      seedWordList(list)
    }()
  }
  seedWordGroup.Wait()
  return nil
}

func ReadOrCreateLanguages() (languages []*models.Language) {
  languageData, err := utils.ReadCSV("tasks/seeds/languages.csv", ';', false)

  if err != nil {
    logs.Critical("unable to read languages file")
  }

  o := orm.NewOrm()

  for _, lang := range languageData {
    language := models.Language{Slug: lang[0]}

    err := o.Read(&language, "Slug")

    language.Name = lang[1]
    language.LettersList = lang[2]
    language.Available, _ = strconv.ParseBool(lang[3])

    if err == nil {
      o.Update(&language)
    } else {
      o.Insert(&language)
    }

    languages = append(languages, &language)
  }
  return languages
}

func readOrCreateWordLists(languages []*models.Language) (lists []*models.WordList) {
  o := orm.NewOrm()

  for _, lang := range languages {
    list := models.WordList{Language: lang, Version: 0}
    _, _, err := o.ReadOrCreate(&list, "Language")
    if err == nil {
      o.LoadRelated(&list, "Language")
      lists = append(lists, &list)
    }
  }
  return
}

func seedWordList(list *models.WordList) {
  o := orm.NewOrm()

  lang := list.Language.Slug

  filePattern := "tasks/seeds/" + web.BConfig.RunMode + "/words_" + lang + "_v*"
  file, version, err := utils.FindFileWithHighestVersion(filePattern, ".txt")
  if err != nil {
    logs.Critical("unable to find file for language = %s", lang)
    return
  }

  if version > list.Version {
    err = seedWords(list, file)
    if err != nil {
      logs.Critical("unable to read and save words from file = %s", file)
      return
    }

    list.Version = version
    o.Update(list, "Version")
    logs.Info("Wordlist for %s language succesfully updated", lang)
  } else {
    logs.Info("There is no need to update Wordlist for %s language", lang)
  }
}

func seedWords(list *models.WordList, file string) error {
  words, err := readWords(file, list)

  if err != nil {
    return err
  }

  o := orm.NewOrm()

  tr, err := o.Begin()
  if err != nil {
    return err
  }

  defer func() {
    if err != nil {
      tr.Rollback()
    } else {
      tr.Commit()
    }
  }()

  startUpdate := time.Now().UTC()
  for _, word := range words {
    if created, _, err := tr.ReadOrCreate(&word, "Name"); err == nil {
      if !created {
        word.Archived = false
        word.UpdatedAt = time.Now()
        tr.Update(&word, "UpdatedAt", "Archived")
      }
    }
  }
  _, err = tr.QueryTable("words").
    Filter("updated_at__lt", startUpdate).
    Filter("word_list_id", list.Id).
    Update(orm.Params{"archived": true})

  return err
}

func readWords(wordsFile string, list *models.WordList) ([]models.Word, error) {
  var words []models.Word

  wordsRaw, err := utils.ReadTXT(wordsFile)

  if err != nil {
    return words, err
  }

  for _, wordRaw := range wordsRaw {
    name := strings.ToLower(wordRaw)

    if utf8.RuneCountInString(name) != models.WordMaxLength {
      continue
    }

    letters := regexp.MustCompile(list.Language.LettersList)
    if !letters.MatchString(name) {
      continue
    }

    words = append(words, models.Word{Name: name, WordList: list})
  }

  return words, nil
}
