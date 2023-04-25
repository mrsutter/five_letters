package utils

import (
  "encoding/csv"
  "os"
)

func ReadCSV(fileName string, separator rune, withHeaders bool) ([][]string, error) {
  file, err := os.Open(fileName)
  if err != nil {
    return [][]string{}, err
  }
  defer file.Close()

  reader := csv.NewReader(file)
  reader.Comma = separator
  rows, err := reader.ReadAll()
  if err != nil {
    return [][]string{}, err
  }

  if withHeaders {
    return rows, nil
  }

  return rows[1:], nil
}
