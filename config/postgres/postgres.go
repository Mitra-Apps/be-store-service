package postgres

import (
	"fmt"
	"os"
	"strings"
	"time"

	prodEntity "github.com/Mitra-Apps/be-store-service/domain/product/entity"
	"github.com/Mitra-Apps/be-store-service/domain/store/entity"
	"github.com/sirupsen/logrus"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connection() *gorm.DB {
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	portStr := ""
	if strings.Trim(port, " ") != "" {
		portStr = fmt.Sprintf("port=%s ", port)
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s %ssslmode=disable TimeZone=Asia/Jakarta", host, username, password, dbName, portStr)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		TranslateError:         false,
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
		&prodEntity.ProductCategory{},
		&prodEntity.ProductType{},
		&prodEntity.UnitOfMeasure{},
		&prodEntity.Product{},
		&prodEntity.ProductImage{},
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
