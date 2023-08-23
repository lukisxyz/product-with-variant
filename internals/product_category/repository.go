package product_category

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
	GetByIDWithTransaction(ctx context.Context, tx pgx.Tx, id ulid.ULID) (*domain.ProductCategory, error)
	GetByID(ctx context.Context, id ulid.ULID) (*domain.ProductCategory, error)
	GetByProductIDCategoryIDWithTransaction(ctx context.Context, tx pgx.Tx, productId, categoryId ulid.ULID) (*domain.ProductCategory, error)
	GetByProductIDCategoryID(ctx context.Context, productId, categoryId ulid.ULID) (*domain.ProductCategory, error)
	GetByProductIDWithTransaction(ctx context.Context, tx pgx.Tx, id ulid.ULID) ([]domain.ProductCategory, error)
	GetByProductID(ctx context.Context, id ulid.ULID) ([]domain.ProductCategory, error)
	SaveWithTransaction(ctx context.Context, tx pgx.Tx, prd *domain.ProductCategory) error
	EditWithTransaction(ctx context.Context, tx pgx.Tx, prd *domain.ProductCategory) error
	DeleteWithTransaction(ctx context.Context, tx pgx.Tx, prd *domain.ProductCategory) error
	GetByCursor(ctx context.Context, limit int, cursor string) ([]domain.ProductCategory, string, error)
	DeleteCategoryWithTransaction(ctx context.Context, tx pgx.Tx, prd *domain.ProductCategory) error
	DeleteProductWithTransaction(ctx context.Context, tx pgx.Tx, prd *domain.ProductCategory) error
}

type repo struct {
	db *pgxpool.Pool
}

// GetByID implements Repo.
func (*repo) GetByProductIDCategoryIDWithTransaction(ctx context.Context, tx pgx.Tx, productId, categoryId ulid.ULID) (*domain.ProductCategory, error) {
	query := `
		SELECT
			pc.product_category_id,
			p.product_id,
			p.name AS product_name,
			p.description AS product_description,
			p.price,
			p.image_preview,
			c.category_id,
			c.name AS category_name,
			c.description AS category_description
		FROM Product_Category pc
		JOIN Product p ON pc.product_id = p.product_id
		JOIN Category c ON pc.category_id = c.category_id
		WHERE pc.product_id = $1
			AND pc.category_id = $2
			AND pc.deleted_at IS NULL
	`
	row := tx.QueryRow(
		ctx,
		query,
		productId,
		categoryId,
	)
	var pc domain.ProductCategory
	if err := row.Scan(
		&pc.ProductCategoryID,
		&pc.Product.ProductID,
		&pc.Product.Name,
		&pc.Product.Description,
		&pc.Product.Price,
		&pc.Product.ImagePreview,
		&pc.Category.CategoryID,
		&pc.Category.Name,
		&pc.Category.Description,
	); err != nil {
		return nil, err
	}
	return &pc, nil
}

// GetByID implements Repo.
func (r *repo) GetByProductIDCategoryID(ctx context.Context, productId, categoryId ulid.ULID) (*domain.ProductCategory, error) {
	query := `
		SELECT
			pc.product_category_id,
			p.product_id,
			p.name AS product_name,
			p.description AS product_description,
			p.price,
			p.image_preview,
			c.category_id,
			c.name AS category_name,
			c.description AS category_description
		FROM Product_Category pc
		JOIN Product p ON pc.product_id = p.product_id
		JOIN Category c ON pc.category_id = c.category_id
		WHERE pc.product_id = $1
			AND pc.category_id = $2
			AND pc.deleted_at IS NULL
	`
	row := r.db.QueryRow(
		ctx,
		query,
		productId,
		categoryId,
	)
	var pc domain.ProductCategory
	if err := row.Scan(
		&pc.ProductCategoryID,
		&pc.Product.ProductID,
		&pc.Product.Name,
		&pc.Product.Description,
		&pc.Product.Price,
		&pc.Product.ImagePreview,
		&pc.Category.CategoryID,
		&pc.Category.Name,
		&pc.Category.Description,
	); err != nil {
		return nil, err
	}
	return &pc, nil
}

// GetByID implements Repo.
func (*repo) GetByProductIDWithTransaction(ctx context.Context, tx pgx.Tx, id ulid.ULID) ([]domain.ProductCategory, error) {
	query := `
		SELECT
			pc.product_category_id,
			p.product_id,
			p.name AS product_name,
			p.description AS product_description,
			p.price,
			p.image_preview,
			c.category_id,
			c.name AS category_name,
			c.description AS category_description
		FROM Product_Category pc
		JOIN Product p ON pc.product_id = p.product_id
		JOIN Category c ON pc.category_id = c.category_id
		WHERE pc.product_id = $1
			AND pc.deleted_at IS NULL
	`
	rows, err := tx.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var relations []domain.ProductCategory
	for rows.Next() {
		var product domain.ProductCategory
		if err := rows.Scan(
			&product.ProductCategoryID,
			&product.Product.ProductID,
			&product.Product.Name,
			&product.Product.Description,
			&product.Product.Price,
			&product.Product.ImagePreview,
			&product.Category.CategoryID,
			&product.Category.Name,
			&product.Category.Description,
		); err != nil {
			return nil, err
		}
		relations = append(relations, product)
	}
	return relations, nil
}

// GetByID implements Repo.
func (r *repo) GetByProductID(ctx context.Context, id ulid.ULID) ([]domain.ProductCategory, error) {
	query := `
		SELECT
			pc.product_category_id,
			p.product_id,
			p.name AS product_name,
			p.description AS product_description,
			p.price,
			p.image_preview,
			c.category_id,
			c.name AS category_name,
			c.description AS category_description
		FROM Product_Category pc
		JOIN Product p ON pc.product_id = p.product_id
		JOIN Category c ON pc.category_id = c.category_id
		WHERE pc.product_id = $1
			AND pc.deleted_at IS NULL
	`
	rows, err := r.db.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var relations []domain.ProductCategory
	for rows.Next() {
		var product domain.ProductCategory
		if err := rows.Scan(
			&product.ProductCategoryID,
			&product.Product.ProductID,
			&product.Product.Name,
			&product.Product.Description,
			&product.Product.Price,
			&product.Product.ImagePreview,
			&product.Category.CategoryID,
			&product.Category.Name,
			&product.Category.Description,
		); err != nil {
			return nil, err
		}
		relations = append(relations, product)
	}
	return relations, nil
}

// GetByID implements Repo.
func (*repo) GetByIDWithTransaction(ctx context.Context, tx pgx.Tx, id ulid.ULID) (*domain.ProductCategory, error) {
	query := `
		SELECT
			pc.product_category_id,
			p.product_id,
			p.name AS product_name,
			p.description AS product_description,
			p.price,
			p.image_preview,
			c.category_id,
			c.name AS category_name,
			c.description AS category_description
		FROM Product_Category pc
		JOIN Product p ON pc.product_id = p.product_id
		JOIN Category c ON pc.category_id = c.category_id
		WHERE pc.product_category_id = $1
			AND pc.deleted_at IS NULL
	`
	row := tx.QueryRow(
		ctx,
		query,
		id,
	)
	var pc domain.ProductCategory
	if err := row.Scan(
		&pc.ProductCategoryID,
		&pc.Product.ProductID,
		&pc.Product.Name,
		&pc.Product.Description,
		&pc.Product.Price,
		&pc.Product.ImagePreview,
		&pc.Category.CategoryID,
		&pc.Category.Name,
		&pc.Category.Description,
	); err != nil {
		return nil, err
	}
	return &pc, nil
}

// GetByID implements Repo.
func (r *repo) GetByID(ctx context.Context, id ulid.ULID) (*domain.ProductCategory, error) {
	query := `
		SELECT
			pc.product_category_id,
			p.product_id,
			p.name AS product_name,
			p.description AS product_description,
			p.price,
			p.image_preview,
			c.category_id,
			c.name AS category_name,
			c.description AS category_description
		FROM Product_Category pc
		JOIN Product p ON pc.product_id = p.product_id
		JOIN Category c ON pc.category_id = c.category_id
		WHERE pc.product_category_id = $1
			AND pc.deleted_at IS NULL
	`
	row := r.db.QueryRow(
		ctx,
		query,
		id,
	)
	var pc domain.ProductCategory
	if err := row.Scan(
		&pc.ProductCategoryID,
		&pc.Product.ProductID,
		&pc.Product.Name,
		&pc.Product.Description,
		&pc.Product.Price,
		&pc.Product.ImagePreview,
		&pc.Category.CategoryID,
		&pc.Category.Name,
		&pc.Category.Description,
	); err != nil {
		return nil, err
	}
	return &pc, nil
}

func (*repo) SaveWithTransaction(ctx context.Context, tx pgx.Tx, pc *domain.ProductCategory) error {
	query := `
		INSERT INTO Product_Category (product_category_id, product_id, category_id)
		VALUES ($1, $2, $3)
	`
	if _, err := tx.Exec(
		ctx,
		query,
		&pc.ProductCategoryID,
		&pc.Product.ProductID,
		&pc.Category.CategoryID,
	); err != nil {
		return err
	}
	return nil
}

func (*repo) EditWithTransaction(ctx context.Context, tx pgx.Tx, prd *domain.ProductCategory) error {
	query := `
		UPDATE Product_Category SET
			product_id = $1,
			category_id = $2
			updated_at = $3
		WHERE
			product_category_id = $4 AND deleted_at IS NULL
	`
	currentTime := time.Now()
	if _, err := tx.Exec(
		ctx,
		query,
		&prd.Product.ProductID,
		&prd.Category.CategoryID,
		currentTime,
		&prd.ProductCategoryID,
	); err != nil {
		return err
	}
	return nil
}

func (*repo) DeleteWithTransaction(ctx context.Context, tx pgx.Tx, prd *domain.ProductCategory) error {
	query := `
		UPDATE Product_Category SET
			deleted_at = $1
		WHERE
			product_category_id = $2 AND deleted_at IS NULL
	`
	currentTime := time.Now()
	if _, err := tx.Exec(
		ctx,
		query,
		currentTime,
		&prd.ProductCategoryID,
	); err != nil {
		return err
	}
	return nil
}

func (*repo) DeleteCategoryWithTransaction(ctx context.Context, tx pgx.Tx, prd *domain.ProductCategory) error {
	query := `
		UPDATE Product_Category SET
			deleted_at = $1
		WHERE
			category_id = $2 AND deleted_at IS NULL
	`
	currentTime := time.Now()
	if _, err := tx.Exec(
		ctx,
		query,
		currentTime,
		&prd.Category.CategoryID,
	); err != nil {
		return err
	}
	return nil
}

func (*repo) DeleteProductWithTransaction(ctx context.Context, tx pgx.Tx, prd *domain.ProductCategory) error {
	query := `
		UPDATE Product_Category SET
			deleted_at = $1
		WHERE
			product_id = $2 AND deleted_at IS NULL
	`
	currentTime := time.Now()
	if _, err := tx.Exec(
		ctx,
		query,
		currentTime,
		&prd.Product.ProductID,
	); err != nil {
		return err
	}
	return nil
}

func (r *repo) GetByCursor(ctx context.Context, limit int, cursor string) ([]domain.ProductCategory, string, error) {
	query := `
		SELECT
			pc.product_category_id,
			p.product_id,
			p.name AS product_name,
			p.description AS product_description,
			p.price,
			p.image_preview,
			c.category_id,
			c.name AS category_name,
			c.description AS category_description
		FROM Product_Category pc
		JOIN Product p ON pc.product_id = p.product_id
		JOIN Category c ON pc.category_id = c.category_id
		WHERE pc.product_category_id = $1
		AND pc.deleted_at IS NULL
			pc.created_at > $1 AND pc.deleted_at IS NULL
		ORDER BY
			pc.created_at
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

	var products []domain.ProductCategory
	for rows.Next() {
		var product domain.ProductCategory
		if err := rows.Scan(
			&product.ProductCategoryID,
			&product.Product.ProductID,
			&product.Product.Name,
			&product.Product.Description,
			&product.Product.Price,
			&product.Product.ImagePreview,
			&product.Category.CategoryID,
			&product.Category.Name,
			&product.Category.Description,
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
