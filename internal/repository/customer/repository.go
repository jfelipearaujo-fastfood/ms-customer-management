package customer

import (
	"context"
	"database/sql"

	"github.com/doug-martin/goqu/v9"
	"github.com/jfelipearaujo-org/ms-customer-management/internal/entity"
	"github.com/jfelipearaujo-org/ms-customer-management/internal/shared/custom_error"
)

const (
	tableName = "customers"
)

type repository struct {
	conn *sql.DB
}

func NewRepository(conn *sql.DB) Repository {
	return &repository{
		conn: conn,
	}
}

func (r *repository) Get(ctx context.Context, id string) (entity.Customer, error) {
	customer := entity.Customer{}

	sql, params, err := goqu.
		From(tableName).
		Select("id", "document_id", "password", "is_anonymous", "created_at", "updated_at").
		Where(goqu.Ex{
			"id": id,
		}).
		ToSQL()

	if err != nil {
		return entity.Customer{}, err
	}

	statement, err := r.conn.QueryContext(ctx, sql, params...)

	if err != nil {
		return entity.Customer{}, err
	}

	defer statement.Close()

	for statement.Next() {
		err = statement.Scan(
			&customer.Id,
			&customer.DocumentId,
			&customer.Password,
			&customer.IsAnonymous,
			&customer.CreatedAt,
			&customer.UpdatedAt)

		if err != nil {
			return entity.Customer{}, err
		}
	}

	if customer.Id == "" {
		return entity.Customer{}, custom_error.ErrCustomerNotFound
	}

	return customer, nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	sql, params, err := goqu.
		Delete(tableName).
		Where(goqu.Ex{
			"id": id,
		}).
		ToSQL()

	if err != nil {
		return err
	}

	result, err := r.conn.ExecContext(ctx, sql, params...)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return custom_error.ErrCustomerNotFound
	}

	return nil
}
