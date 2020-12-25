// Copyright @ 2020 OPS Inc.
//
// Author: Jinlong Yang
//

package handler

import (
	"time"
)

func ConvertTime(date time.Time) string {
	return date.Format("2006-01-02T15:04:05Z")
}
