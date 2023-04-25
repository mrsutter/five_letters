package models

type ErrorItem struct {
  Code  string
  Field string
}

type Error struct {
  Status  int
  Code    string
  Message string
  Details []ErrorItem
}
