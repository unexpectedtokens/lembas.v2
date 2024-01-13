package view_util

import (
	"fmt"
)

type comparablefl interface {
	float32 | float64
}

func FloatToString[T comparablefl](fl float32) string {
	return fmt.Sprintf("%.2f", fl)
}
