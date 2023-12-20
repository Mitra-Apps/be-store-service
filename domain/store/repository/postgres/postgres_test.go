package postgres_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Mitra-Apps/be-store-service/domain/store/entity"
	repository "github.com/Mitra-Apps/be-store-service/domain/store/repository/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestCreateStore(t *testing.T) {
	// Create a new SQL mock
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(t, err)

	// Create a new instance of your postgres struct with the mocked DB
	p := repository.NewPostgres(gormDB)

	// Expect the Begin, Save, and Commit queries to be executed
	mock.ExpectBegin()

	// Expect the INSERT INTO "stores" query
	mock.ExpectExec(`INSERT INTO "stores" ("created_at","created_by","updated_at","updated_by","deleted_at","deleted_by","store_name","address","city","state","zip_code","phone","email","website","map_location") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15) RETURNING "id"`).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Expect the INSERT INTO "store_hours" query
	mock.ExpectExec(`INSERT INTO "store_hours" (.+) VALUES (.+) RETURNING "id"`).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Expect the INSERT INTO "store_images" query
	mock.ExpectExec(`INSERT INTO "store_images" (.+) VALUES (.+) RETURNING "id"`).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Expect the INSERT INTO "store_tags" query
	mock.ExpectExec(`INSERT INTO "store_tags" (.+) VALUES (.+) RETURNING "id"`).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	// Call the CreateStore function with a sample store
	ctx := context.Background()
	store := &entity.Store{
		StoreName:   "Test Store",
		Address:     "Test Address",
		City:        "Test City",
		State:       "Test State",
		ZipCode:     "Test ZipCode",
		Phone:       "Test Phone",
		Email:       "test@store.com",
		Website:     "http://teststore.com",
		MapLocation: "Test MapLocation",
		Hours: []entity.StoreHour{
			{
				DayOfWeek: entity.Monday,
				Open:      "09:00 AM",
				Close:     "05:00 PM",
			},
			// Add other StoreHour entries as needed
		},
		Images: []entity.StoreImage{
			{
				ImageType: "Logo",
				ImageURL:  "http://teststore.com/logo.png",
			},
			// Add other StoreImage entries as needed
		},
		Tags: []entity.StoreTag{
			{
				TagName: "Tag1",
			},
			// Add other StoreTag entries as needed
		},
	}

	createdStore, err := p.CreateStore(ctx, store)

	// Assert that there are no errors
	assert.NoError(t, err)
	assert.NotNil(t, createdStore)

	// Ensure all expectations were met
	require.NoError(t, mock.ExpectationsWereMet())
}
