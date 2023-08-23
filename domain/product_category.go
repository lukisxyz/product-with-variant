package domain

import (
	"time"

	"github.com/oklog/ulid/v2"
	"gopkg.in/guregu/null.v4"
)

type ProductCategory struct {
	ProductCategoryID ulid.ULID
	Product           Product
	Category          Category
	CreatedAt         time.Time
	UpdatedAt         null.Time
	DeletedAt         null.Time
}

func NewRelationProductCategory(pId, cId ulid.ULID) (ProductCategory, error) {
	id := ulid.Make()
	prd := Product{
		ProductID: pId,
	}
	cat := Category{
		CategoryID: cId,
	}
	return ProductCategory{
		ProductCategoryID: id,
		Product:           prd,
		Category:          cat,
		CreatedAt:         time.Now(),
	}, nil
}
