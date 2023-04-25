package utils

import (
  "fmt"
  "path/filepath"
  "strconv"
  "strings"
)

func FindFileWithHighestVersion(fileNamePattern string, ext string) (string, int, error) {
  var highestVersion int
  var highestVersionFile string

  files, err := filepath.Glob(fileNamePattern + ext)
  if err != nil {
    return "", 0, err
  }
  if len(files) == 0 {
    return "", 0, fmt.Errorf("no matching files found in directory")
  }

  for _, file := range files {
    parts := strings.Split(file, "_")
    versionStr := parts[len(parts)-1]

    versionStr = strings.TrimSuffix(versionStr, ext)
    versionStr = strings.TrimPrefix(versionStr, "v")

    version, err := strconv.Atoi(versionStr)

    if err != nil {
      return "", 0, err
    }

    if version > highestVersion {
      highestVersion = version
      highestVersionFile = file
    }
  }

  return highestVersionFile, highestVersion, nil
}
