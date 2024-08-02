package delete_request_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jfelipearaujo-org/ms-customer-management/internal/adapter/database"
	"github.com/jfelipearaujo-org/ms-customer-management/internal/entity"
	"github.com/jfelipearaujo-org/ms-customer-management/internal/environment"
	"github.com/jfelipearaujo-org/ms-customer-management/internal/repository/delete_request"
	"github.com/stretchr/testify/assert"
)

func TestGetByCustomerId(t *testing.T) {
	t.Run("Should return a deletion request", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		mock.ExpectQuery("SELECT (.+) FROM (.+)?customer_deletion_requests(.+)?").
			WillReturnRows(sqlmock.NewRows([]string{"id", "customer_id", "name", "address", "phone", "executed", "created_at", "updated_at"}).
				AddRow("id", "customer_id", "name", "address", "phone", false, time.Now(), time.Now()))

		config := &environment.Config{
			DbConfig: &environment.DatabaseConfig{
				Url: "postgres://host:1234",
			},
		}

		service := database.NewDatabase(config)
		service.(*database.Service).Client = db

		repo := delete_request.NewRepository(db)

		// Act
		res, err := repo.GetByCustomerId(ctx, "customer_id")

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "id", res.Id)
	})

	t.Run("Should return an error when try to get a deletion request", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		mock.ExpectQuery("SELECT (.+) FROM (.+)?customer_deletion_requests(.+)?").
			WillReturnError(errors.New("error"))

		config := &environment.Config{
			DbConfig: &environment.DatabaseConfig{
				Url: "postgres://host:1234",
			},
		}

		service := database.NewDatabase(config)
		service.(*database.Service).Client = db

		repo := delete_request.NewRepository(db)

		// Act
		res, err := repo.GetByCustomerId(ctx, "customer_id")

		// Assert
		assert.Error(t, err)
		assert.Empty(t, res)
	})

	t.Run("Should return an error when deletion request is not found", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		mock.ExpectQuery("SELECT (.+) FROM (.+)?customer_deletion_requests(.+)?").
			WillReturnRows(sqlmock.NewRows([]string{"id", "customer_id", "name", "address", "phone", "executed", "created_at", "updated_at"}))

		config := &environment.Config{
			DbConfig: &environment.DatabaseConfig{
				Url: "postgres://host:1234",
			},
		}

		service := database.NewDatabase(config)
		service.(*database.Service).Client = db

		repo := delete_request.NewRepository(db)

		// Act
		res, err := repo.GetByCustomerId(ctx, "customer_id")

		// Assert
		assert.Error(t, err)
		assert.Empty(t, res)
	})

	t.Run("Should return an error when try to scan the results", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		mock.ExpectQuery("SELECT (.+) FROM (.+)?customer_deletion_requests(.+)?").
			WillReturnRows(sqlmock.NewRows([]string{"id", "customer_id", "name", "address", "phone", "executed", "created_at", "updated_at"}).
				AddRow("id", "customer_id", "name", "address", "phone", false, 123, time.Now()))

		config := &environment.Config{
			DbConfig: &environment.DatabaseConfig{
				Url: "postgres://host:1234",
			},
		}

		service := database.NewDatabase(config)
		service.(*database.Service).Client = db

		repo := delete_request.NewRepository(db)

		// Act
		res, err := repo.GetByCustomerId(ctx, "customer_id")

		// Assert
		assert.Error(t, err)
		assert.Empty(t, res)
	})
}

func TestCreate(t *testing.T) {
	t.Run("Should create a deletion request", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO (.+)?customer_deletion_requests(.+)?").
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		repo := delete_request.NewRepository(db)

		// Act
		err = repo.Create(ctx, entity.DeletionRequest{
			Id:         "id",
			CustomerId: "customer_id",
			Name:       "name",
			Address:    "address",
			Phone:      "phone",
			Executed:   false,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		})

		// Assert
		assert.NoError(t, err)
	})

	t.Run("Should return an error when try to create a deletion request", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO (.+)?customer_deletion_requests(.+)?").
			WillReturnError(errors.New("error"))
		mock.ExpectRollback()

		repo := delete_request.NewRepository(db)

		// Act
		err = repo.Create(ctx, entity.DeletionRequest{
			Id:         "id",
			CustomerId: "customer_id",
			Name:       "name",
			Address:    "address",
			Phone:      "phone",
			Executed:   false,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		})

		// Assert
		assert.Error(t, err)
	})

	t.Run("Should return an error when try to begin a transaction", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		mock.ExpectBegin().WillReturnError(errors.New("error"))

		repo := delete_request.NewRepository(db)

		// Act
		err = repo.Create(ctx, entity.DeletionRequest{
			Id:         "id",
			CustomerId: "customer_id",
			Name:       "name",
			Address:    "address",
			Phone:      "phone",
			Executed:   false,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		})

		// Assert
		assert.Error(t, err)
	})

	t.Run("Should return an error when try to commit a deletion request", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO (.+)?customer_deletion_requests(.+)?").
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit().WillReturnError(errors.New("error"))

		repo := delete_request.NewRepository(db)

		// Act
		err = repo.Create(ctx, entity.DeletionRequest{
			Id:         "id",
			CustomerId: "customer_id",
			Name:       "name",
			Address:    "address",
			Phone:      "phone",
			Executed:   false,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		})

		// Assert
		assert.Error(t, err)
	})
}
