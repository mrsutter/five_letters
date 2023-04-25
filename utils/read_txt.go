package utils

import (
  "bufio"
  "os"
)

func ReadTXT(fileName string) ([]string, error) {
  var lines []string

  file, err := os.Open(fileName)
  if err != nil {
    return lines, err
  }
  defer file.Close()

  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    lines = append(lines, scanner.Text())
  }

  if err := scanner.Err(); err != nil {
    return lines, err
  }

  return lines, err
}
