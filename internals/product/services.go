package product

import (
	"context"
	"errors"
	"flukis/product/domain"
	"flukis/product/internals/product_category"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oklog/ulid/v2"
)

type Service interface {
	GetProductByID(ctx context.Context, id ulid.ULID) (domain.ProductDetailDTO, error)
	CreateProduct(ctx context.Context, name, desc string, price float64) (domain.ProductDTO, error)
	UpdateImageProduct(ctx context.Context, id ulid.ULID, image []byte) (domain.ProductDTO, error)
	UpdateDataProduct(ctx context.Context, id ulid.ULID, name, desc string, price float64) (domain.ProductDTO, error)
	GetProductsByCursor(ctx context.Context, limit int, cursor string) (res []domain.ProductDetailDTO, length int, nextCursor string, err error)
	DeleteProduct(ctx context.Context, id ulid.ULID) error
	UpdateCategoryProduct(ctx context.Context, id ulid.ULID, categoryIds []ulid.ULID) error
	DeleteCategoryProductBatch(ctx context.Context, id ulid.ULID, categoryIds []ulid.ULID) error
}

type service struct {
	repo                 Repo
	categoryRelationrepo product_category.Repo
	db                   *pgxpool.Pool
}

func (s *service) DeleteCategoryProductBatch(ctx context.Context, id ulid.ULID, categoryIds []ulid.ULID) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}

	for idx := range categoryIds {
		prd, err := s.categoryRelationrepo.GetByProductIDCategoryID(ctx, id, categoryIds[idx])
		if err != nil {
			if err := tx.Rollback(ctx); err != nil {
				return err
			}
			return err
		}
		err = s.categoryRelationrepo.DeleteWithTransaction(ctx, tx, prd)
		if err != nil {
			if err := tx.Rollback(ctx); err != nil {
				return err
			}
			return err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) UpdateCategoryProduct(ctx context.Context, id ulid.ULID, categoryIds []ulid.ULID) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}

	for idx := range categoryIds {
		_, err := s.categoryRelationrepo.GetByProductIDCategoryID(ctx, id, categoryIds[idx])
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				newRelation, err := domain.NewRelationProductCategory(
					id,
					categoryIds[idx],
				)
				if err != nil {
					if err := tx.Rollback(ctx); err != nil {
						return err
					}
					return err
				}
				err = s.categoryRelationrepo.SaveWithTransaction(ctx, tx, &newRelation)
				if err != nil {
					if err := tx.Rollback(ctx); err != nil {
						return err
					}
					return err
				}
			}
			continue
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

// DeleteProduct implements Service.
func (s *service) DeleteProduct(ctx context.Context, id ulid.ULID) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}

	currPrd, err := s.repo.GetByIDWithTransaction(ctx, tx, id)
	if err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return err
		}
		return err
	}

	err = s.repo.DeleteWithTransaction(ctx, tx, currPrd)
	if err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return err
		}
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetProductsByCursor(ctx context.Context, limit int, cursor string) (res []domain.ProductDetailDTO, length int, nextCursor string, err error) {
	prd, nextCursor, err := s.repo.GetByCursor(ctx, limit, cursor)
	if err != nil {
		return []domain.ProductDetailDTO{}, 0, "", err
	}
	dataLen := len(prd)
	if dataLen == 0 {
		return []domain.ProductDetailDTO{}, 0, "", nil
	}
	var data = make([]domain.ProductDetailDTO, dataLen)
	for i := range prd {
		data[i].ID = prd[i].ProductID
		data[i].Name = prd[i].Name
		data[i].Description = prd[i].Description
		data[i].Image = prd[i].ImagePreview
		data[i].Price = prd[i].Price
	}
	return data, dataLen, nextCursor, nil
}

// CreateProduct implements Service.
func (s *service) UpdateDataProduct(ctx context.Context, id ulid.ULID, name, desc string, price float64) (domain.ProductDTO, error) {
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

	currentPrd.Name = name
	currentPrd.Description = desc
	currentPrd.Price = price

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
func (s *service) GetProductByID(ctx context.Context, id ulid.ULID) (domain.ProductDetailDTO, error) {
	prd, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return domain.ProductDetailDTO{}, err
	}
	category, err := s.categoryRelationrepo.GetByProductID(ctx, id)
	if err != nil {
		return domain.ProductDetailDTO{}, err
	}
	var categories = make([]domain.CategoriesDTO, 0)
	for idx := range category {
		buf := domain.CategoriesDTO{
			ID:          category[idx].Category.CategoryID,
			Name:        category[idx].Category.Name,
			Description: category[idx].Category.Description,
		}
		categories = append(categories, buf)
	}
	res := domain.ProductDetailDTO{
		ProductDTO: domain.ProductDTO{
			ID:          id,
			Name:        prd.Name,
			Description: prd.Description,
			Price:       prd.Price,
			Image:       prd.ImagePreview,
		},
		Category:  categories,
		Attribute: []domain.AttributesDTO{},
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
	categoryRelationrepo product_category.Repo,
	db *pgxpool.Pool,
) Service {
	return &service{
		repo:                 repo,
		db:                   db,
		categoryRelationrepo: categoryRelationrepo,
	}
}
