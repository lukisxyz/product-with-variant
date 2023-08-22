package product_attributes

import (
	"context"
	"flukis/product/domain"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oklog/ulid/v2"
)

type Service interface {
	GetAttrById(ctx context.Context, id ulid.ULID) (domain.AttributesDTO, error)
	GetAttrByCursor(ctx context.Context, limit int, cursor string) ([]domain.AttributesDTO, int, string, error)
	DeleteAttr(ctx context.Context, id ulid.ULID) error
	UpdateNameAttr(ctx context.Context, id ulid.ULID, name string) (domain.AttributesDTO, error)
	CreateAttr(ctx context.Context, name string) (domain.AttributesDTO, error)
}

type service struct {
	repo Repo
	db   *pgxpool.Pool
}

// CreateAttr implements Service.
func (s *service) CreateAttr(ctx context.Context, name string) (domain.AttributesDTO, error) {
	newAttr, err := domain.NewAttribute(name)
	if err != nil {
		return domain.AttributesDTO{}, err
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return domain.AttributesDTO{}, err
	}

	err = s.repo.SaveWithTransaction(ctx, tx, &newAttr)
	if err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return domain.AttributesDTO{}, err
		}
		return domain.AttributesDTO{}, err
	}

	res := domain.AttributesDTO{
		ID:   newAttr.AttributeID,
		Name: newAttr.Name,
	}

	err = tx.Commit(ctx)
	if err != nil {
		return domain.AttributesDTO{}, err
	}
	return res, nil
}

// DeleteAttr implements Service.
func (s *service) DeleteAttr(ctx context.Context, id ulid.ULID) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}

	currAttr, err := s.repo.GetByIDWithTransaction(ctx, tx, id)
	if err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return err
		}
		return err
	}

	err = s.repo.DeleteWithTransaction(ctx, tx, currAttr)
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

// GetAttrByCursor implements Service.
func (s *service) GetAttrByCursor(ctx context.Context, limit int, cursor string) (res []domain.AttributesDTO, length int, nextCursor string, err error) {
	attr, nextCursor, err := s.repo.GetByCursor(ctx, limit, cursor)
	if err != nil {
		return []domain.AttributesDTO{}, 0, "", err
	}
	dataLen := len(attr)
	if dataLen == 0 {
		return []domain.AttributesDTO{}, 0, "", nil
	}
	var data = make([]domain.AttributesDTO, dataLen)
	for i := range attr {
		data[i].ID = attr[i].AttributeID
		data[i].Name = attr[i].Name
	}
	return data, dataLen, nextCursor, nil
}

// GetAttrById implements Service.
func (s *service) GetAttrById(ctx context.Context, id ulid.ULID) (domain.AttributesDTO, error) {
	attr, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return domain.AttributesDTO{}, err
	}
	res := domain.AttributesDTO{
		ID:   attr.AttributeID,
		Name: attr.Name,
	}
	return res, nil
}

// UpdateNameAttr implements Service.
func (s *service) UpdateNameAttr(ctx context.Context, id ulid.ULID, name string) (domain.AttributesDTO, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return domain.AttributesDTO{}, err
	}

	currAttr, err := s.repo.GetByIDWithTransaction(ctx, tx, id)
	if err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return domain.AttributesDTO{}, err
		}
		return domain.AttributesDTO{}, err
	}
	currAttr.Name = name

	err = s.repo.EditWithTransaction(ctx, tx, currAttr)
	if err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return domain.AttributesDTO{}, err
		}
		return domain.AttributesDTO{}, err
	}
	res := domain.AttributesDTO{
		ID:   currAttr.AttributeID,
		Name: currAttr.Name,
	}

	err = tx.Commit(ctx)
	if err != nil {
		return domain.AttributesDTO{}, err
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
