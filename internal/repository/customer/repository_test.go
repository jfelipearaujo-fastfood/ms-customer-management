package customer_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jfelipearaujo-org/ms-customer-management/internal/adapter/database"
	"github.com/jfelipearaujo-org/ms-customer-management/internal/environment"
	"github.com/jfelipearaujo-org/ms-customer-management/internal/repository/customer"
	"github.com/stretchr/testify/assert"
)

func TestRepository_Get(t *testing.T) {
	t.Run("Should return a customer", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery("SELECT (.+) FROM (.+)?customers(.+)?").
			WillReturnRows(sqlmock.NewRows([]string{"id", "document_id", "password", "is_anonymous", "created_at", "updated_at"}).
				AddRow("id", "document_id", "password", true, time.Now(), time.Now()))

		config := &environment.Config{
			DbConfig: &environment.DatabaseConfig{
				Url: "postgres://host:1234",
			},
		}

		service := database.NewDatabase(config)
		service.(*database.Service).Client = db

		repo := customer.NewRepository(db)

		// Act
		res, err := repo.Get(context.Background(), "id")

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "id", res.Id)
	})

	t.Run("Should return an error", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery("SELECT (.+) FROM (.+)?customers(.+)?").
			WillReturnError(errors.New("error"))

		config := &environment.Config{
			DbConfig: &environment.DatabaseConfig{
				Url: "postgres://host:1234",
			},
		}

		service := database.NewDatabase(config)
		service.(*database.Service).Client = db

		repo := customer.NewRepository(db)

		// Act
		res, err := repo.Get(context.Background(), "id")

		// Assert
		assert.Error(t, err)
		assert.Empty(t, res)
	})

	t.Run("Should return an error when the customer is not found", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery("SELECT (.+) FROM (.+)?customers(.+)?").
			WillReturnRows(sqlmock.NewRows([]string{"id", "document_id", "password", "is_anonymous", "created_at", "updated_at"}))

		config := &environment.Config{
			DbConfig: &environment.DatabaseConfig{
				Url: "postgres://host:1234",
			},
		}

		service := database.NewDatabase(config)
		service.(*database.Service).Client = db

		repo := customer.NewRepository(db)

		// Act
		res, err := repo.Get(context.Background(), "id")

		// Assert
		assert.Error(t, err)
		assert.Empty(t, res)
	})
}

func TestRepository_Delete(t *testing.T) {
	t.Run("Should delete a customer", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectExec("DELETE FROM (.+)?customers(.+)?").
			WillReturnResult(sqlmock.NewResult(1, 1))

		config := &environment.Config{
			DbConfig: &environment.DatabaseConfig{
				Url: "postgres://host:1234",
			},
		}

		service := database.NewDatabase(config)
		service.(*database.Service).Client = db

		repo := customer.NewRepository(db)

		// Act
		err = repo.Delete(context.Background(), "id")

		// Assert
		assert.NoError(t, err)
	})

	t.Run("Should return an error", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectExec("DELETE FROM (.+)?customers(.+)?").
			WillReturnError(errors.New("error"))

		config := &environment.Config{
			DbConfig: &environment.DatabaseConfig{
				Url: "postgres://host:1234",
			},
		}

		service := database.NewDatabase(config)
		service.(*database.Service).Client = db

		repo := customer.NewRepository(db)

		// Act
		err = repo.Delete(context.Background(), "id")

		// Assert
		assert.Error(t, err)
	})

	t.Run("Should return an error when the customer is not found", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectExec("DELETE FROM (.+)?customers(.+)?").
			WillReturnResult(sqlmock.NewResult(0, 0))

		config := &environment.Config{
			DbConfig: &environment.DatabaseConfig{
				Url: "postgres://host:1234",
			},
		}

		service := database.NewDatabase(config)
		service.(*database.Service).Client = db

		repo := customer.NewRepository(db)

		// Act
		err = repo.Delete(context.Background(), "id")

		// Assert
		assert.Error(t, err)
	})
}
