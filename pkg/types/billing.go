package types

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Billing struct {
	gorm.Model
	AppId       uint
	Cpu         float64
	Memory      float64
	TimeStarted time.Time
}

type BillingAccount struct {
	gorm.Model
	AccountId uint    `json:"account_id"`
	Total     float64 `json:"total"`
}
