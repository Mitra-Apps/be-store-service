package postgres_test

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Mitra-Apps/be-store-service/domain/store/entity"
	repositoryPostgres "github.com/Mitra-Apps/be-store-service/domain/store/repository/postgres"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	storeID = "7d56be32-70a2-4f49-b66b-63e6f8e719d5"
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

func Test_postgres_OpenCloseStore(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	assert.NoError(t, err)
	repo := repositoryPostgres.NewPostgres(gormDB)
	ctx := context.Background()
	storeIDUuid, _ := uuid.Parse(storeID)

	type args struct {
		ctx      context.Context
		storeId  uuid.UUID
		isActive bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "OpenCloseStores_NoGormError_ReturnNil",
			args: args{
				ctx:      ctx,
				storeId:  storeIDUuid,
				isActive: false,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectBegin()
			mock.ExpectExec(regexp.QuoteMeta(`UPDATE "stores"`)).
				WithArgs(tt.args.isActive, tt.args.storeId).
				WillReturnResult(sqlmock.NewErrorResult(nil))
			mock.ExpectCommit()
			if err := repo.OpenCloseStore(tt.args.ctx, tt.args.storeId, tt.args.isActive); (err != nil) != tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
