package category

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
	GetByIDWithTransaction(ctx context.Context, tx pgx.Tx, id ulid.ULID) (*domain.Category, error)
	GetByID(ctx context.Context, id ulid.ULID) (*domain.Category, error)
	SaveWithTransaction(ctx context.Context, tx pgx.Tx, cat *domain.Category) error
	EditWithTransaction(ctx context.Context, tx pgx.Tx, cat *domain.Category) error
	DeleteWithTransaction(ctx context.Context, tx pgx.Tx, cat *domain.Category) error
	GetByCursor(ctx context.Context, limit int, cursor string) ([]domain.Category, string, error)
}

type repo struct {
	db *pgxpool.Pool
}

// GetByID implements Repo.
func (*repo) GetByIDWithTransaction(ctx context.Context, tx pgx.Tx, id ulid.ULID) (*domain.Category, error) {
	query := `
		SELECT
			category_id,
			name
		FROM
			category
		WHERE
			category_id = $1 AND deleted_at IS NULL
	`
	row := tx.QueryRow(
		ctx,
		query,
		id,
	)
	var cat domain.Category
	if err := row.Scan(
		&cat.CategoryID,
		&cat.Name,
	); err != nil {
		return nil, err
	}
	return &cat, nil
}

// GetByID implements Repo.
func (r *repo) GetByID(ctx context.Context, id ulid.ULID) (*domain.Category, error) {
	query := `
		SELECT
			category_id,
			name,
			description
		FROM
			category
		WHERE
			category_id = $1  AND deleted_at IS NULL
	`
	row := r.db.QueryRow(
		ctx,
		query,
		id,
	)
	var cat domain.Category
	if err := row.Scan(
		&cat.CategoryID,
		&cat.Name,
		&cat.Description,
	); err != nil {
		return nil, err
	}
	return &cat, nil
}

func (*repo) SaveWithTransaction(ctx context.Context, tx pgx.Tx, cat *domain.Category) error {
	query := `
		INSERT INTO category
			(category_id, name, description, created_at)
		VALUES
			($1, $2, $3, $4)
	`
	if _, err := tx.Exec(
		ctx,
		query,
		&cat.CategoryID,
		&cat.Name,
		&cat.Description,
		&cat.CreatedAt,
	); err != nil {
		return err
	}
	return nil
}

func (*repo) EditWithTransaction(ctx context.Context, tx pgx.Tx, cat *domain.Category) error {
	query := `
		UPDATE category SET
			name = $1,
			updated_at = $2,
			description = $4
		WHERE
			category_id = $3 AND deleted_at IS NULL
	`
	currentTime := time.Now()
	if _, err := tx.Exec(
		ctx,
		query,
		&cat.Name,
		currentTime,
		&cat.CategoryID,
		&cat.Description,
	); err != nil {
		return err
	}
	return nil
}

func (*repo) DeleteWithTransaction(ctx context.Context, tx pgx.Tx, cat *domain.Category) error {
	query := `
		UPDATE category SET
			deleted_at = $1
		WHERE
			category_id = $2 AND deleted_at IS NULL
	`
	currentTime := time.Now()
	if _, err := tx.Exec(
		ctx,
		query,
		currentTime,
		&cat.CategoryID,
	); err != nil {
		return err
	}
	return nil
}

func (r *repo) GetByCursor(ctx context.Context, limit int, cursor string) ([]domain.Category, string, error) {
	query := `
		SELECT
			category_id, name, description, created_at FROM category
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

	var categories []domain.Category
	for rows.Next() {
		var category domain.Category
		if err := rows.Scan(&category.CategoryID, &category.Name, &category.Description, &category.CreatedAt); err != nil {
			return nil, "", err
		}
		categories = append(categories, category)
	}

	nextCursor := ""
	if len(categories) == limit {
		nextCursor = helper.EncodeCursor(categories[len(categories)-1].CreatedAt)
	}

	return categories, nextCursor, nil
}

func NewRepo(db *pgxpool.Pool) Repo {
	return &repo{
		db: db,
	}
}
