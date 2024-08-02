package delete_request

import (
	"context"
	"database/sql"

	"github.com/doug-martin/goqu/v9"
	"github.com/jfelipearaujo-org/ms-customer-management/internal/entity"
	"github.com/jfelipearaujo-org/ms-customer-management/internal/shared/custom_error"
)

const (
	tableName = "customer_deletion_requests"
)

type repository struct {
	conn *sql.DB
}

func NewRepository(conn *sql.DB) Repository {
	return &repository{
		conn: conn,
	}
}

func (r *repository) GetByCustomerId(ctx context.Context, customerId string) (entity.DeletionRequest, error) {
	deletionRequest := entity.DeletionRequest{}

	sql, params, err := goqu.
		From(tableName).
		Select("id", "customer_id", "name", "address", "phone", "executed", "created_at", "updated_at").
		Where(goqu.Ex{
			"customer_id": customerId,
			"executed":    false,
		}).
		ToSQL()
	if err != nil {
		return entity.DeletionRequest{}, err
	}

	statement, err := r.conn.QueryContext(ctx, sql, params...)

	if err != nil {
		return entity.DeletionRequest{}, err
	}

	defer statement.Close()

	for statement.Next() {
		err = statement.Scan(
			&deletionRequest.Id,
			&deletionRequest.CustomerId,
			&deletionRequest.Name,
			&deletionRequest.Address,
			&deletionRequest.Phone,
			&deletionRequest.Executed,
			&deletionRequest.CreatedAt,
			&deletionRequest.UpdatedAt)

		if err != nil {
			return entity.DeletionRequest{}, err
		}
	}

	if deletionRequest.Id == "" {
		return entity.DeletionRequest{}, custom_error.ErrDeletionRequestNotFound
	}

	return deletionRequest, nil
}

func (r *repository) Create(ctx context.Context, request entity.DeletionRequest) error {
	tx, err := r.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	sql, params, err := goqu.Insert(tableName).
		Cols("id", "customer_id", "name", "address", "phone", "executed", "created_at", "updated_at").
		Vals(goqu.Vals{
			request.Id,
			request.CustomerId,
			request.Name,
			request.Address,
			request.Phone,
			request.Executed,
			request.CreatedAt,
			request.UpdatedAt,
		}).
		ToSQL()
	if err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, sql, params...); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
