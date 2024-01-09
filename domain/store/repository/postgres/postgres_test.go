package postgres_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Mitra-Apps/be-store-service/domain/store/entity"
	repositoryPostgres "github.com/Mitra-Apps/be-store-service/domain/store/repository/postgres"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestCreateStore(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	assert.NoError(t, err)

	repo := repositoryPostgres.NewPostgres(gormDB)

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO stores").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	store := &entity.Store{
		StoreName:   "Test Store",
		Address:     "123 Main St",
		City:        "Cityville",
		State:       "CA",
		ZipCode:     "12345",
		Phone:       "555-1234",
		Email:       "test@store.com",
		Website:     "http://www.teststore.com",
		LocationLat: 0.0,
		LocationLng: 0.0,
	}

	result, err := repo.CreateStore(context.Background(), store)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}
