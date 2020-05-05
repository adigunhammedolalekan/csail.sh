package services

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/saas/hostgolang/pkg/types"
)

func CreateDatabaseConnection(connectUri string) (*gorm.DB, error) {
	db, err := gorm.Open("postgres", connectUri)
	if err != nil {
		return nil, err
	}
	if err := db.DB().Ping(); err != nil {
		return nil, err
	}
	migrate(db)
	return db, nil
}

func migrate(db *gorm.DB) {
	db.Debug().AutoMigrate(&types.Account{},
		&types.App{}, &types.Environment{},
		&types.Release{}, &types.DeploymentSettings{},
		&types.Resource{}, &types.ResourceEnv{},
		&types.ResourceConfig{}, &types.Quota{}, &types.Plan{}, &types.Domain{})
}

func destroy(db *gorm.DB) {
	db.Debug().DropTable(&types.Account{},
		&types.App{}, &types.Environment{},
		&types.Release{}, &types.DeploymentSettings{},
		&types.Resource{}, &types.ResourceEnv{},
		&types.ResourceConfig{}, &types.Quota{}, &types.Plan{}, &types.Domain{})
}