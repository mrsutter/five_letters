package models

import (
  "fmt"
  "github.com/beego/beego/v2/client/orm"
  "time"
)

const (
  StateActive = "active"
  StateWasted = "wasted"
  StateWon    = "won"
)

const MaxAttemptsCount = 6
const MaxHoursForGame = 24

type Game struct {
  Id            int
  State         string
  AttemptsCount int
  CreatedAt     time.Time  `orm:"auto_now_add;type(datetime)"`
  UpdatedAt     time.Time  `orm:"auto_now;type(datetime)"`
  Word          *Word      `orm:"rel(fk)"`
  User          *User      `orm:"rel(fk)"`
  Attempts      []*Attempt `orm:"reverse(many);on_delete(cascade)"`
}

func (g *Game) CanTransitionTo(newState string) bool {
  switch g.State {
  case StateActive:
    return newState == StateWasted || newState == StateWon
  case StateWasted:
    return false
  case StateWon:
    return false
  }
  return false
}

func (g *Game) TransitionTo(newState string) error {
  if !g.CanTransitionTo(newState) {
    return fmt.Errorf("invalid transition from %s to %s", g.State, newState)
  }
  g.State = newState
  return nil
}

func (g *Game) SetValuesAfterAttempt(attempt Attempt) {
  if g.State != StateActive ||
    g != attempt.Game ||
    g.AttemptsCount == MaxAttemptsCount {
    return
  }

  g.AttemptsCount = g.AttemptsCount + 1

  if attempt.Successful() {
    g.TransitionTo(StateWon)
  } else {
    if g.AttemptsCount == MaxAttemptsCount {
      g.TransitionTo(StateWasted)
    }
  }
}

func init() {
  orm.RegisterModel(new(Game))
}

func (g *Game) TableName() string {
  return "games"
}
