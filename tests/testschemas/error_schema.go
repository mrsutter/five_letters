package testschemas

import (
  "golang.org/x/exp/slices"
)

type ErrorSchema struct {
  Status  int               `json:"status"`
  Code    string            `json:"code"`
  Message string            `json:"message"`
  Details []ErrorSchemaItem `json:"details"`
}

type ErrorSchemaItem struct {
  Code  string `json:"code"`
  Field string `json:"field"`
}

func (e *ErrorSchema) FindItem(field string) ErrorSchemaItem {
  idx := slices.IndexFunc(e.Details,
    func(i ErrorSchemaItem) bool {
      return i.Field == field
    })
  return e.Details[idx]
}
