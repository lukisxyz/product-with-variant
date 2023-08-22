package domain

import (
	"time"

	"github.com/oklog/ulid/v2"
	"gopkg.in/guregu/null.v4"
)

type Category struct {
	CategoryID  ulid.ULID
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   null.Time
	DeletedAt   null.Time
}

type CategoriesDTO struct {
	ID          ulid.ULID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"desc"`
}

func NewCategory(name, desc string) (Category, error) {
	id := ulid.Make()
	return Category{
		CategoryID:  id,
		Name:        name,
		Description: desc,
		CreatedAt:   time.Now(),
	}, nil
}
