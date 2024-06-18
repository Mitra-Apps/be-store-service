package postgres

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	prodEntity "github.com/Mitra-Apps/be-store-service/domain/product/entity"
	"github.com/Mitra-Apps/be-store-service/domain/store/entity"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	integrationTest = flag.Bool("integration", false, "run integration test")
	db              *gorm.DB
	err             error
)

func TestMain(m *testing.M) {
	if *integrationTest {
		if _, err := os.Stat("./../../../../.env"); !os.IsNotExist(err) {
			err := godotenv.Load(os.ExpandEnv("./../../../../.env"))
			if err != nil {
				log.Fatalf("Error getting env %v\n", err)
			}
		}

		username := os.Getenv("DB_USERNAME")
		password := os.Getenv("DB_PASSWORD")
		host := os.Getenv("DB_HOST")
		dbName := os.Getenv("DB_NAME_TEST")
		port := os.Getenv("DB_PORT")
		portStr := ""
		if strings.Trim(port, " ") != "" {
			portStr = fmt.Sprintf("port=%s ", port)
		}
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s %ssslmode=disable TimeZone=Asia/Jakarta", host, username, password, dbName, portStr)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			TranslateError:         false,
			SkipDefaultTransaction: true,
		})
		if err != nil {
			logrus.Fatalf("Failed to connect to database: %v", err)
		}

		db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")

		db.Migrator().DropTable(
			&entity.Store{},
			&entity.StoreImage{},
			&entity.StoreHour{},
			&prodEntity.ProductCategory{},
			&prodEntity.ProductType{},
			&prodEntity.UnitOfMeasure{},
			&prodEntity.Product{},
			&prodEntity.ProductImage{},
		)
		db.AutoMigrate(
			&entity.Store{},
			&entity.StoreImage{},
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
	}

	// Run tests
	code := m.Run()
	os.Exit(code)
}
