package postgres

import (
	"fmt"
	"os"
	"time"

	"github.com/Mitra-Apps/be-store-service/domain/store/entity"
	"github.com/sirupsen/logrus"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connection() *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", os.Getenv("DB_HOST"), os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		TranslateError:         true,
		SkipDefaultTransaction: true,
	})
	if err != nil {
		logrus.Fatalf("Failed to connect to database: %v", err)
	}

	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")

	err = db.AutoMigrate(
		&entity.Store{},
		&entity.StoreImage{},
		&entity.StoreTag{},
		&entity.StoreHour{},
	)
	if err != nil {
		logrus.Fatalf("Failed to migrate table: %v", err)
	}

	sqlDb, err := db.DB()
	if err != nil {
		logrus.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDb.SetMaxIdleConns(10)
	sqlDb.SetMaxOpenConns(100)
	sqlDb.SetConnMaxLifetime(time.Hour * 6)

	return db
}
