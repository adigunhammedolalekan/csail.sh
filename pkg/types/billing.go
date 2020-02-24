package types

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
)

var (
	Plans = map[string]Plan {
	"SR": {Name: "STARTER", Alias: "SR", Price: 12, Info: PlanInfo{Cpu: 1, Memory: 1}},
	"B1": {Name: "BEGINNER", Alias: "B1", Price: 22, Info: PlanInfo{Cpu: 1, Memory: 3}},
	"B2": {Name: "BEGINNER 2", Alias: "B2", Price: 22, Info: PlanInfo{Cpu: 2, Memory: 2}},
	"S1": {Name: "STANDARD 1", Alias: "S1", Price: 27, Info: PlanInfo{Cpu: 2, Memory: 4}},
	"S2": {Name: "STANDARD 2", Alias: "S2", Price: 50, Info: PlanInfo{Cpu: 4, Memory: 8}},
	"S3": {Name: "STANDARD 3", Alias: "S3", Price: 90, Info: PlanInfo{Cpu: 6, Memory: 16}},
	"M1": {Name: "MASTER", Alias: "M1", Price: 170, Info: PlanInfo{Cpu: 8, Memory: 32}},
	"M2": {Name: "GRAND MASTER", Alias: "M2", Price: 340, Info: PlanInfo{Cpu: 16, Memory: 64}},
	}
	DefaultPlan = Plan{
		Name:  "TEST",
		Alias: "TST",
		Price: 0,
		Info:  PlanInfo{
			Memory: 0.1, Cpu: 0.1,
		},
	}
)

type Plan struct {
	gorm.Model
	AppId uint `json:"app_id"`
	Name string `json:"name"`
	Alias string `json:"alias"`
	Price float64 `json:"price"`
	Info PlanInfo `json:"info"`
}

func NewPlan(appId uint, planAlias string) *Plan {
	var plan = DefaultPlan
	if planAlias != "" {
		if p, ok := Plans[planAlias]; ok {
			plan = p
		}
	}
	plan.AppId = appId
	return &plan
}

func (p Plan) PriceString() string {
	return fmt.Sprintf("$%f", p.Price)
}

type PlanInfo struct {
	Cpu float64
	Memory float64
}

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
