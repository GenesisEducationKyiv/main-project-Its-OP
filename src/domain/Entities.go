package domain

import "time"

type Price struct {
	Amount    float64
	Currency  string
	Timestamp time.Time
}
