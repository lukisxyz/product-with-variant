package product_variant

import (
	"context"
	"flukis/product/domain"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oklog/ulid/v2"
)

type Service interface {
	GetVariantByID(ctx context.Context, id ulid.ULID) (domain.VariantDTO, error)
	CreateVariant(ctx context.Context, name, desc string, price float64, mainId ulid.ULID) (domain.VariantDTO, error)
	UpdateDataVariant(ctx context.Context, id ulid.ULID, name, desc string, price float64, mainId ulid.ULID) (domain.VariantDTO, error)
	GetVariantsByCursor(ctx context.Context, limit int, cursor string) (res []domain.VariantDetailDTO, length int, nextCursor string, err error)
	DeleteVariant(ctx context.Context, id ulid.ULID) error
}

type service struct {
	repo Repo
	db   *pgxpool.Pool
}

// DeleteVariant implements Service.
func (s *service) DeleteVariant(ctx context.Context, id ulid.ULID) error {
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

func (s *service) GetVariantsByCursor(ctx context.Context, limit int, cursor string) (res []domain.VariantDetailDTO, length int, nextCursor string, err error) {
	prd, nextCursor, err := s.repo.GetByCursor(ctx, limit, cursor)
	if err != nil {
		return []domain.VariantDetailDTO{}, 0, "", err
	}
	dataLen := len(prd)
	if dataLen == 0 {
		return []domain.VariantDetailDTO{}, 0, "", nil
	}
	var data = make([]domain.VariantDetailDTO, dataLen)
	for i := range prd {
		data[i].ID = prd[i].VariantID
		data[i].Name = prd[i].Name
		data[i].Description = prd[i].Description
		data[i].Price = prd[i].Price
		data[i].Image = prd[i].MainProduct.ImagePreview
		data[i].MainProductID = prd[i].MainProduct.ProductID
		data[i].MainProductName = prd[i].MainProduct.Name
	}
	return data, dataLen, nextCursor, nil
}

// CreateVariant implements Service.
func (s *service) UpdateDataVariant(ctx context.Context, id ulid.ULID, name, desc string, price float64, mainId ulid.ULID) (domain.VariantDTO, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return domain.VariantDTO{}, err
	}

	currentPrd, err := s.repo.GetByIDWithTransaction(ctx, tx, id)
	if err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return domain.VariantDTO{}, err
		}
		return domain.VariantDTO{}, err
	}

	currentPrd.Name = name
	currentPrd.Description = desc
	currentPrd.Price = price
	currentPrd.MainProduct.ProductID = mainId

	err = s.repo.EditWithTransaction(ctx, tx, currentPrd)
	if err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return domain.VariantDTO{}, err
		}
		return domain.VariantDTO{}, err
	}

	res := domain.VariantDTO{
		ID:              currentPrd.VariantID,
		Name:            currentPrd.Name,
		Description:     currentPrd.Description,
		Price:           currentPrd.Price,
		Image:           currentPrd.MainProduct.ImagePreview,
		MainProductID:   currentPrd.MainProduct.ProductID,
		MainProductName: currentPrd.MainProduct.Name,
	}

	err = tx.Commit(ctx)
	if err != nil {
		return domain.VariantDTO{}, err
	}
	return res, nil
}

// CreateVariant implements Service.
func (s *service) CreateVariant(ctx context.Context, name, desc string, price float64, mainId ulid.ULID) (domain.VariantDTO, error) {
	newPrd, err := domain.NewVariant(name, desc, price, mainId)
	if err != nil {
		return domain.VariantDTO{}, err
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return domain.VariantDTO{}, err
	}

	err = s.repo.SaveWithTransaction(ctx, tx, &newPrd)
	if err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return domain.VariantDTO{}, err
		}
		return domain.VariantDTO{}, err
	}

	res := domain.VariantDTO{
		ID:              newPrd.VariantID,
		Name:            newPrd.Name,
		Description:     newPrd.Description,
		Price:           newPrd.Price,
		Image:           newPrd.MainProduct.ImagePreview,
		MainProductID:   newPrd.MainProduct.ProductID,
		MainProductName: newPrd.MainProduct.Name,
	}

	err = tx.Commit(ctx)
	if err != nil {
		return domain.VariantDTO{}, err
	}
	return res, nil
}

// GetVariantByID implements Service.
func (s *service) GetVariantByID(ctx context.Context, id ulid.ULID) (domain.VariantDTO, error) {
	prd, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return domain.VariantDTO{}, err
	}
	res := domain.VariantDTO{
		ID:              prd.VariantID,
		Name:            prd.Name,
		Description:     prd.Description,
		Price:           prd.Price,
		Image:           prd.MainProduct.ImagePreview,
		MainProductID:   prd.MainProduct.ProductID,
		MainProductName: prd.MainProduct.Name,
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
