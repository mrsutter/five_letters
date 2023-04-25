package utils

import (
  "reflect"
  "strings"
)

func GetJsonFieldName(struc interface{}, fieldName string) string {
  field, _ := reflect.TypeOf(struc).FieldByName(fieldName)

  tag := field.Tag.Get("json")

  if tag == "" {
    return field.Name
  }

  if tag == "-" {
    return ""
  }

  if i := strings.Index(tag, ","); i != -1 {
    if i == 0 {
      return field.Name
    } else {
      return tag[:i]
    }
  }

  return tag
}
