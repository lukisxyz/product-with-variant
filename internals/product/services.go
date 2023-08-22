package product

import (
	"context"
	"flukis/product/domain"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oklog/ulid/v2"
)

type Service interface {
	GetProductByID(ctx context.Context, id ulid.ULID) (domain.ProductDTO, error)
	CreateProduct(ctx context.Context, name, desc string, price float64) (domain.ProductDTO, error)
	UpdateImageProduct(ctx context.Context, id ulid.ULID, image []byte) (domain.ProductDTO, error)
}

type service struct {
	repo Repo
	db   *pgxpool.Pool
}

// CreateProduct implements Service.
func (s *service) CreateProduct(ctx context.Context, name, desc string, price float64) (domain.ProductDTO, error) {
	newPrd, err := domain.NewProduct(name, desc, price)
	if err != nil {
		return domain.ProductDTO{}, err
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return domain.ProductDTO{}, err
	}

	err = s.repo.SaveWithTransaction(ctx, tx, &newPrd)
	if err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return domain.ProductDTO{}, err
		}
		return domain.ProductDTO{}, err
	}

	res := domain.ProductDTO{
		ID:          newPrd.ProductID,
		Name:        newPrd.Name,
		Description: newPrd.Description,
		Price:       newPrd.Price,
		Image:       newPrd.ImagePreview,
	}

	err = tx.Commit(ctx)
	if err != nil {
		return domain.ProductDTO{}, err
	}
	return res, nil
}

// GetProductByID implements Service.
func (s *service) GetProductByID(ctx context.Context, id ulid.ULID) (domain.ProductDTO, error) {
	attr, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return domain.ProductDTO{}, err
	}
	res := domain.ProductDTO{
		ID:          attr.ProductID,
		Name:        attr.Name,
		Description: attr.Description,
		Price:       attr.Price,
		Image:       attr.ImagePreview,
	}
	return res, nil
}

// CreateProduct implements Service.
func (s *service) UpdateImageProduct(ctx context.Context, id ulid.ULID, image []byte) (domain.ProductDTO, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return domain.ProductDTO{}, err
	}
	currentPrd, err := s.repo.GetByIDWithTransaction(ctx, tx, id)
	if err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return domain.ProductDTO{}, err
		}
		return domain.ProductDTO{}, err
	}
	currentPrd.ImagePreview = image

	err = s.repo.EditWithTransaction(ctx, tx, currentPrd)
	if err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return domain.ProductDTO{}, err
		}
		return domain.ProductDTO{}, err
	}

	res := domain.ProductDTO{
		ID:          currentPrd.ProductID,
		Name:        currentPrd.Name,
		Description: currentPrd.Description,
		Price:       currentPrd.Price,
		Image:       currentPrd.ImagePreview,
	}

	err = tx.Commit(ctx)
	if err != nil {
		return domain.ProductDTO{}, err
	}
	return res, nil
}

func NewService(
	repo Repo,
	db *pgxpool.Pool,
) Service {
	return &service{
		repo: repo,
		db:   db,
	}
}
