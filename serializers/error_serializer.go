package serializers

import (
  "five_letters/models"
)

type ErrorItem struct {
  Code  string `json:"code"`
  Field string `json:"field"`
}

func (errorItemSerializer *ErrorItem) Serialize(errorItem models.ErrorItem) *ErrorItem {
  errorItemSerializer.Code = errorItem.Code
  errorItemSerializer.Field = errorItem.Field

  return errorItemSerializer
}

type Error struct {
  Status  int          `json:"status"`
  Code    string       `json:"code"`
  Message string       `json:"message"`
  Details []*ErrorItem `json:"details"`
}

func (errorSerializer *Error) Serialize(err models.Error) *Error {
  errorSerializer.Status = err.Status
  errorSerializer.Code = err.Code
  errorSerializer.Message = err.Message

  if len(err.Details) == 0 {
    errorSerializer.Details = []*ErrorItem{}
  } else {
    for _, errItem := range err.Details {
      errItemSerializer := ErrorItem{}
      errorSerializer.Details = append(
        errorSerializer.Details,
        errItemSerializer.Serialize(errItem),
      )
    }
  }
  return errorSerializer
}
