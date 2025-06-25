package migrations

import (
	"thiennguyen.dev/welab-healthcare-app/infra/database"
	"thiennguyen.dev/welab-healthcare-app/models"
)

func Migrate() {
	var migrationModels = []interface{}{&models.Example{}}
	err := database.GetDB().AutoMigrate(migrationModels...)
	if err != nil {
		return
	}
}
