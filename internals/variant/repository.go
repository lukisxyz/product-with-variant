package variant

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
	GetByIDWithTransaction(ctx context.Context, tx pgx.Tx, id ulid.ULID) (*domain.Variant, error)
	GetByID(ctx context.Context, id ulid.ULID) (*domain.Variant, error)
	SaveWithTransaction(ctx context.Context, tx pgx.Tx, vrn *domain.Variant) error
	EditWithTransaction(ctx context.Context, tx pgx.Tx, vrn *domain.Variant) error
	DeleteWithTransaction(ctx context.Context, tx pgx.Tx, vrn *domain.Variant) error
	GetByCursor(ctx context.Context, limit int, cursor string) ([]domain.Variant, string, error)
}

type repo struct {
	db *pgxpool.Pool
}

// GetByID implements Repo.
func (*repo) GetByIDWithTransaction(ctx context.Context, tx pgx.Tx, id ulid.ULID) (*domain.Variant, error) {
	query := `
		SELECT
			v.variant_id,
			v.name AS variant_name,
			v.description AS variant_description,
			v.price AS variant_price,
			p.product_id,
			p.name AS product_name,
			p.description AS product_description,
			p.price AS product_price,
			p.image_preview
		FROM
			Variant AS v
		LEFT JOIN
			Product AS p ON v.main_product_id = p.product_id
		WHERE
			v.variant_id = $1 AND v.deleted_at IS NULL AND p.deleted_at is NULL
		LIMIT 1;
	`
	row := tx.QueryRow(
		ctx,
		query,
		id,
	)
	var variant domain.Variant
	var mainProduct domain.Product
	if err := row.Scan(
		&variant.VariantID,
		&variant.Name,
		&variant.Description,
		&variant.Price,
		&mainProduct.ProductID,
		&mainProduct.Name,
		&mainProduct.Description,
		&mainProduct.Price,
		&mainProduct.ImagePreview,
	); err != nil {
		return nil, err
	}

	variant.MainProduct = mainProduct
	return &variant, nil
}

// GetByID implements Repo.
func (r *repo) GetByID(ctx context.Context, id ulid.ULID) (*domain.Variant, error) {
	query := `
		SELECT
			v.variant_id,
			v.name AS variant_name,
			v.description AS variant_description,
			v.price AS variant_price,
			p.product_id,
			p.name AS product_name,
			p.description AS product_description,
			p.price AS product_price,
			p.image_preview
		FROM
			Variant AS v
		LEFT JOIN
			Product AS p ON v.main_product_id = p.product_id
		WHERE
			v.variant_id = $1 AND v.deleted_at IS NULL AND p.deleted_at is NULL
		LIMIT 1;
	`
	row := r.db.QueryRow(
		ctx,
		query,
		id,
	)
	var variant domain.Variant
	var mainProduct domain.Product
	if err := row.Scan(
		&variant.VariantID,
		&variant.Name,
		&variant.Description,
		&variant.Price,
		&mainProduct.ProductID,
		&mainProduct.Name,
		&mainProduct.Description,
		&mainProduct.Price,
		&mainProduct.ImagePreview,
	); err != nil {
		return nil, err
	}
	variant.MainProduct = mainProduct
	return &variant, nil
}

func (*repo) SaveWithTransaction(ctx context.Context, tx pgx.Tx, vrn *domain.Variant) error {
	query := `
		INSERT INTO Variant
			(variant_id, main_product_id, name, description, price)
		VALUES
			($1, $2, $3, $4, $5)
	`
	if _, err := tx.Exec(
		ctx,
		query,
		&vrn.VariantID,
		&vrn.MainProduct.ProductID,
		&vrn.Name,
		&vrn.Description,
		&vrn.Price,
	); err != nil {
		return err
	}
	return nil
}

func (*repo) EditWithTransaction(ctx context.Context, tx pgx.Tx, vrn *domain.Variant) error {
	query := `
		UPDATE Variant SET
			name = $1,
			description = $2,
			price = $3 = $4,
			updated_at = $5,
			main_product_id = $7
		WHERE
			variant_id = $6 AND deleted_at IS NULL
	`
	currentTime := time.Now()
	if _, err := tx.Exec(
		ctx,
		query,
		&vrn.Name,
		&vrn.Description,
		&vrn.Price,
		currentTime,
		&vrn.VariantID,
		&vrn.MainProduct.ProductID,
	); err != nil {
		return err
	}
	return nil
}

func (*repo) DeleteWithTransaction(ctx context.Context, tx pgx.Tx, vrn *domain.Variant) error {
	query := `
		UPDATE Variant SET
			deleted_at = $1
		WHERE
			variant_id = $2 AND deleted_at IS NULL
	`
	currentTime := time.Now()
	if _, err := tx.Exec(
		ctx,
		query,
		currentTime,
		&vrn.VariantID,
	); err != nil {
		return err
	}
	return nil
}

func (r *repo) GetByCursor(ctx context.Context, limit int, cursor string) ([]domain.Variant, string, error) {
	query := `
		SELECT
			v.variant_id,
			v.name AS variant_name,
			v.description AS variant_description,
			v.price AS variant_price,
			p.product_id,
			p.name AS product_name,
			p.image_preview
		FROM
			Variant AS v
		LEFT JOIN
			Product AS p ON v.main_product_id = p.product_id
		WHERE
			v.created_at > $1 AND v.deleted_at IS NULL AND p.deleted_at is NULL
		ORDER BY
			v.created_at
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

	var variants []domain.Variant
	for rows.Next() {
		var variant domain.Variant
		var product domain.Product
		if err := rows.Scan(
			&variant.VariantID,
			&variant.Name,
			&variant.Description,
			&variant.Price,
			&product.ProductID,
			&product.Name,
			&product.ImagePreview,
		); err != nil {
			return nil, "", err
		}
		variant.MainProduct = product
		variants = append(variants, variant)
	}

	nextCursor := ""
	if len(variants) == limit {
		nextCursor = helper.EncodeCursor(variants[len(variants)-1].CreatedAt)
	}

	return variants, nextCursor, nil
}

func NewRepo(db *pgxpool.Pool) Repo {
	return &repo{
		db: db,
	}
}
