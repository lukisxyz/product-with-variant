package category

import (
	"context"
	"flukis/product/domain"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oklog/ulid/v2"
)

type Service interface {
	GetCategoryById(ctx context.Context, id ulid.ULID) (domain.CategoriesDTO, error)
	GetCategoryByCursor(ctx context.Context, limit int, cursor string) ([]domain.CategoriesDTO, int, string, error)
	DeleteCategory(ctx context.Context, id ulid.ULID) error
	UpdateCategory(ctx context.Context, id ulid.ULID, name, desc string) (domain.CategoriesDTO, error)
	CreateCategory(ctx context.Context, name, desc string) (domain.CategoriesDTO, error)
}

type service struct {
	repo Repo
	db   *pgxpool.Pool
}

// Createcat implements Service.
func (s *service) CreateCategory(ctx context.Context, name, desc string) (domain.CategoriesDTO, error) {
	newcat, err := domain.NewCategory(name, desc)
	if err != nil {
		return domain.CategoriesDTO{}, err
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return domain.CategoriesDTO{}, err
	}

	err = s.repo.SaveWithTransaction(ctx, tx, &newcat)
	if err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return domain.CategoriesDTO{}, err
		}
		return domain.CategoriesDTO{}, err
	}

	res := domain.CategoriesDTO{
		ID:          newcat.CategoryID,
		Name:        newcat.Name,
		Description: newcat.Description,
	}

	err = tx.Commit(ctx)
	if err != nil {
		return domain.CategoriesDTO{}, err
	}
	return res, nil
}

// Deletecat implements Service.
func (s *service) DeleteCategory(ctx context.Context, id ulid.ULID) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}

	currcat, err := s.repo.GetByIDWithTransaction(ctx, tx, id)
	if err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return err
		}
		return err
	}

	err = s.repo.DeleteWithTransaction(ctx, tx, currcat)
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

// GetcatByCursor implements Service.
func (s *service) GetCategoryByCursor(ctx context.Context, limit int, cursor string) (res []domain.CategoriesDTO, length int, nextCursor string, err error) {
	category, nextCursor, err := s.repo.GetByCursor(ctx, limit, cursor)
	if err != nil {
		return []domain.CategoriesDTO{}, 0, "", err
	}
	dataLen := len(category)
	if dataLen == 0 {
		return []domain.CategoriesDTO{}, 0, "", nil
	}
	var data = make([]domain.CategoriesDTO, dataLen)
	for i := range category {
		data[i].ID = category[i].CategoryID
		data[i].Name = category[i].Name
		data[i].Description = category[i].Description
	}
	return data, dataLen, nextCursor, nil
}

// GetcatById implements Service.
func (s *service) GetCategoryById(ctx context.Context, id ulid.ULID) (domain.CategoriesDTO, error) {
	cat, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return domain.CategoriesDTO{}, err
	}
	res := domain.CategoriesDTO{
		ID:          cat.CategoryID,
		Name:        cat.Name,
		Description: cat.Description,
	}
	return res, nil
}

// UpdateNamecat implements Service.
func (s *service) UpdateCategory(ctx context.Context, id ulid.ULID, name, desc string) (domain.CategoriesDTO, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return domain.CategoriesDTO{}, err
	}

	currcat, err := s.repo.GetByIDWithTransaction(ctx, tx, id)
	if err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return domain.CategoriesDTO{}, err
		}
		return domain.CategoriesDTO{}, err
	}
	currcat.Name = name
	currcat.Description = desc

	err = s.repo.EditWithTransaction(ctx, tx, currcat)
	if err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return domain.CategoriesDTO{}, err
		}
		return domain.CategoriesDTO{}, err
	}
	res := domain.CategoriesDTO{
		ID:          currcat.CategoryID,
		Name:        currcat.Name,
		Description: currcat.Description,
	}

	err = tx.Commit(ctx)
	if err != nil {
		return domain.CategoriesDTO{}, err
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
