package product

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
	GetByIDWithTransaction(ctx context.Context, tx pgx.Tx, id ulid.ULID) (*domain.Product, error)
	GetByID(ctx context.Context, id ulid.ULID) (*domain.Product, error)
	SaveWithTransaction(ctx context.Context, tx pgx.Tx, prd *domain.Product) error
	EditWithTransaction(ctx context.Context, tx pgx.Tx, prd *domain.Product) error
	DeleteWithTransaction(ctx context.Context, tx pgx.Tx, prd *domain.Product) error
	GetByCursor(ctx context.Context, limit int, cursor string) ([]domain.Product, string, error)
}

type repo struct {
	db *pgxpool.Pool
}

// GetByID implements Repo.
func (*repo) GetByIDWithTransaction(ctx context.Context, tx pgx.Tx, id ulid.ULID) (*domain.Product, error) {
	query := `
		SELECT
			product_id,
			name,
			description,
			price,
			image_preview
		FROM
			Product
		WHERE
			product_id = $1 AND deleted_at IS NULL
	`
	row := tx.QueryRow(
		ctx,
		query,
		id,
	)
	var prd domain.Product
	if err := row.Scan(
		&prd.ProductID,
		&prd.Name,
		&prd.Description,
		&prd.Price,
		&prd.ImagePreview,
	); err != nil {
		return nil, err
	}
	return &prd, nil
}

// GetByID implements Repo.
func (r *repo) GetByID(ctx context.Context, id ulid.ULID) (*domain.Product, error) {
	query := `
		SELECT
			product_id,
			name,
			description,
			price,
			image_preview
		FROM
			Product
		WHERE
			product_id = $1 AND deleted_at IS NULL
	`
	row := r.db.QueryRow(
		ctx,
		query,
		id,
	)
	var prd domain.Product
	if err := row.Scan(
		&prd.ProductID,
		&prd.Name,
		&prd.Description,
		&prd.Price,
		&prd.ImagePreview,
	); err != nil {
		return nil, err
	}
	return &prd, nil
}

func (*repo) SaveWithTransaction(ctx context.Context, tx pgx.Tx, prd *domain.Product) error {
	query := `
		INSERT INTO Product
			(product_id, name, description, price)
		VALUES
			($1, $2, $3, $4)
	`
	if _, err := tx.Exec(
		ctx,
		query,
		&prd.ProductID,
		&prd.Name,
		&prd.Description,
		&prd.Price,
	); err != nil {
		return err
	}
	return nil
}

func (*repo) EditWithTransaction(ctx context.Context, tx pgx.Tx, prd *domain.Product) error {
	query := `
		UPDATE Product SET
			name = $1,
			description = $2,
			price = $3,
			image_preview = $4,
			updated_at = $5
		WHERE
			product_id = $6 AND deleted_at IS NULL
	`
	currentTime := time.Now()
	if _, err := tx.Exec(
		ctx,
		query,
		&prd.Name,
		&prd.Description,
		&prd.Price,
		&prd.ImagePreview,
		currentTime,
		&prd.ProductID,
	); err != nil {
		return err
	}
	return nil
}

func (*repo) DeleteWithTransaction(ctx context.Context, tx pgx.Tx, prd *domain.Product) error {
	query := `
		UPDATE Product SET
			deleted_at = $1
		WHERE
			product_id = $2 AND deleted_at IS NULL
	`
	currentTime := time.Now()
	if _, err := tx.Exec(
		ctx,
		query,
		currentTime,
		&prd.ProductID,
	); err != nil {
		return err
	}
	return nil
}

func (r *repo) GetByCursor(ctx context.Context, limit int, cursor string) ([]domain.Product, string, error) {
	query := `
		SELECT
			product_id,
			name,
			description,
			price,
			image_preview,
			created_at
		FROM
			Product
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

	var products []domain.Product
	for rows.Next() {
		var product domain.Product
		if err := rows.Scan(
			&product.ProductID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.ImagePreview,
			&product.CreatedAt,
		); err != nil {
			return nil, "", err
		}
		products = append(products, product)
	}

	nextCursor := ""
	if len(products) == limit {
		nextCursor = helper.EncodeCursor(products[len(products)-1].CreatedAt)
	}

	return products, nextCursor, nil
}

func NewRepo(db *pgxpool.Pool) Repo {
	return &repo{
		db: db,
	}
}
