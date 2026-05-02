package migration

import (
	"log"

	"github.com/deltrexgg/ai-code-editor-server/internals/infra"
	"github.com/deltrexgg/ai-code-editor-server/internals/models"
)


func AutoMigrate() {
	if err := infra.DataBaseClient.AutoMigrate(
		//users
		&models.Users{},
		//projects
		&models.Projects{},
	); err != nil {
		log.Fatalf("Migration Failed : %v ", err)
	}

	log.Println("Migration Successful")
}
