package attribute

import (
	"context"
	"flukis/product/domain"
	"flukis/product/utils/helper"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

type Repo interface {
	GetByIDWithTransaction(ctx context.Context, tx pgx.Tx, id ulid.ULID) (*domain.Attribute, error)
	GetByID(ctx context.Context, id ulid.ULID) (*domain.Attribute, error)
	SaveWithTransaction(ctx context.Context, tx pgx.Tx, attr *domain.Attribute) error
	EditWithTransaction(ctx context.Context, tx pgx.Tx, attr *domain.Attribute) error
	DeleteWithTransaction(ctx context.Context, tx pgx.Tx, attr *domain.Attribute) error
	GetByCursor(ctx context.Context, limit int, cursor string) ([]domain.Attribute, string, error)
}

type repo struct {
	db *pgxpool.Pool
}

// GetByID implements Repo.
func (*repo) GetByIDWithTransaction(ctx context.Context, tx pgx.Tx, id ulid.ULID) (*domain.Attribute, error) {
	query := `
		SELECT
			attribute_id,
			name
		FROM
			Attribute
		WHERE
			attribute_id = $1 AND deleted_at IS NULL
	`
	row := tx.QueryRow(
		ctx,
		query,
		id,
	)
	var attr domain.Attribute
	if err := row.Scan(
		&attr.AttributeID,
		&attr.Name,
	); err != nil {
		return nil, err
	}
	return &attr, nil
}

// GetByID implements Repo.
func (r *repo) GetByID(ctx context.Context, id ulid.ULID) (*domain.Attribute, error) {
	query := `
		SELECT
			attribute_id,
			name
		FROM
			Attribute
		WHERE
			attribute_id = $1 AND deleted_at IS NULL
	`
	row := r.db.QueryRow(
		ctx,
		query,
		id,
	)
	var attr domain.Attribute
	if err := row.Scan(
		&attr.AttributeID,
		&attr.Name,
	); err != nil {
		return nil, err
	}
	return &attr, nil
}

func (*repo) SaveWithTransaction(ctx context.Context, tx pgx.Tx, attr *domain.Attribute) error {
	query := `
		INSERT INTO Attribute
			(attribute_id, name, created_at)
		VALUES
			($1, $2, $3)
	`
	if _, err := tx.Exec(
		ctx,
		query,
		&attr.AttributeID,
		&attr.Name,
		&attr.CreatedAt,
	); err != nil {
		return err
	}
	return nil
}

func (*repo) EditWithTransaction(ctx context.Context, tx pgx.Tx, attr *domain.Attribute) error {
	query := `
		UPDATE Attribute SET
			name = $1,
			updated_at = $2
		WHERE
			attribute_id = $3 AND deleted_at IS NULL
	`
	currentTime := time.Now()
	if _, err := tx.Exec(
		ctx,
		query,
		&attr.Name,
		currentTime,
		&attr.AttributeID,
	); err != nil {
		return err
	}
	return nil
}

func (*repo) DeleteWithTransaction(ctx context.Context, tx pgx.Tx, attr *domain.Attribute) error {
	query := `
		UPDATE Attribute SET
			deleted_at = $1
		WHERE
			attribute_id = $2 AND deleted_at IS NULL
	`
	currentTime := time.Now()
	if _, err := tx.Exec(
		ctx,
		query,
		currentTime,
		&attr.AttributeID,
	); err != nil {
		return err
	}
	return nil
}

func (r *repo) GetByCursor(ctx context.Context, limit int, cursor string) ([]domain.Attribute, string, error) {
	query := `
		SELECT
			attribute_id, name, created_at FROM Attribute
		WHERE
			created_at > $1 AND deleted_at IS NULL
		ORDER BY
			created_at
		LIMIT $2
	`
	decodedCursor, err := helper.DecodeCursor(cursor)
	if err != nil && cursor != "" {
		log.Warn().Err(err).Msg("failed to decode cursor")
		return nil, "", err
	}

	rows, err := r.db.Query(ctx, query, decodedCursor, limit)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	var attributes []domain.Attribute
	for rows.Next() {
		var attribute domain.Attribute
		if err := rows.Scan(&attribute.AttributeID, &attribute.Name, &attribute.CreatedAt); err != nil {
			return nil, "", err
		}
		attributes = append(attributes, attribute)
	}

	nextCursor := ""
	if len(attributes) == limit {
		nextCursor = helper.EncodeCursor(attributes[len(attributes)-1].CreatedAt)
	}

	return attributes, nextCursor, nil
}

func NewRepo(db *pgxpool.Pool) Repo {
	return &repo{
		db: db,
	}
}
