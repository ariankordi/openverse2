package models

import (
  "github.com/go-gorp/gorp"

  "time"
)

func ProperSince(nullTime *gorp.NullTime) int64 {
  postTime := nullTime.Time
  timeProper, _ := time.Parse("Mon Jan 2 15:04:05 MST 2006", postTime.Format("Mon Jan 2 15:04:05 EDT 2006"))
  return int64(time.Now().Sub(timeProper))
}

func ProperNow() time.Time {
  properTime, _ := time.Parse("Mon Jan 2 15:04:05 EDT 2006", time.Now().Format("Mon Jan 2 15:04:05 MST 2006"))
  return properTime
}
