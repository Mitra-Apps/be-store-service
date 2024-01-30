package postgres_test

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Mitra-Apps/be-store-service/domain/product/entity"
	repositoryPostgres "github.com/Mitra-Apps/be-store-service/domain/product/repository/postgres"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	storeID = "7d56be32-70a2-4f49-b66b-63e6f8e719d5"
)

func getRepoMock(t *testing.T) (sqlmock.Sqlmock, *repositoryPostgres.Postgres) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	assert.NoError(t, err)
	repo := repositoryPostgres.NewPostgres(gormDB)
	return mock, repo
}

func Test_postgres_GetProductsByStoreIdAndNames(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	assert.NoError(t, err)
	repo := repositoryPostgres.NewPostgres(gormDB)
	ctx := context.Background()
	type args struct {
		ctx     context.Context
		names   []string
		storeId uuid.UUID
	}
	tests := []struct {
		storeId uuid.UUID
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "GetProductsByStoreIdAndNames_NoGormError_Success",
			args: args{
				ctx: ctx,
				names: []string{
					"indomie",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectExec(`SELECT * FROM "products"`).
				WithArgs(sqlmock.AnyArg).
				WillReturnResult(sqlmock.NewResult(0, 0))
			_, err := repo.GetProductsByStoreIdAndNames(tt.args.ctx, tt.args.storeId, tt.args.names)
			if (err != nil) == tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func Test_postgres_UpsertProducts(t *testing.T) {
	mock, repo := getRepoMock(t)
	ctx := context.Background()
	type args struct {
		ctx      context.Context
		products []*entity.Product
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "UpsertProducts_NoGormError_Success",
			args: args{
				ctx:      ctx,
				products: []*entity.Product{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectExec(regexp.QuoteMeta(`insert "products"`)).
				WithArgs(tt.args.products).
				WillReturnResult(sqlmock.NewErrorResult(nil))
			err := repo.UpsertProducts(tt.args.ctx, tt.args.products)
			if (err != nil) == tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

/*
func Test_postgres_UpsertUnitOfMeasure(t *testing.T) {
	type args struct {
		ctx context.Context
		uom *entity.UnitOfMeasure
	}
	tests := []struct {
		name    string
		p       *postgres
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.p.UpsertUnitOfMeasure(tt.args.ctx, tt.args.uom); (err != nil) != tt.wantErr {
				t.Errorf("postgres.UpsertUnitOfMeasure() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_postgres_UpsertProductCategory(t *testing.T) {
	type args struct {
		ctx          context.Context
		prodCategory *entity.ProductCategory
	}
	tests := []struct {
		name    string
		p       *postgres
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.p.UpsertProductCategory(tt.args.ctx, tt.args.prodCategory); (err != nil) != tt.wantErr {
				t.Errorf("postgres.UpsertProductCategory() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_postgres_UpsertProductType(t *testing.T) {
	type args struct {
		ctx      context.Context
		prodType *entity.ProductType
	}
	tests := []struct {
		name    string
		p       *postgres
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.p.UpsertProductType(tt.args.ctx, tt.args.prodType); (err != nil) != tt.wantErr {
				t.Errorf("postgres.UpsertProductType() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
*/
