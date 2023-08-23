package domain

import (
	"time"

	"github.com/oklog/ulid/v2"
	"gopkg.in/guregu/null.v4"
)

type Variant struct {
	VariantID   ulid.ULID
	MainProduct Product
	Name        string
	Description string
	Price       float64
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   null.Time
}

type VariantDTO struct {
	ID              ulid.ULID `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"desc"`
	Price           float64   `json:"price"`
	Image           []byte    `json:"image"`
	MainProductID   ulid.ULID `json:"main_id"`
	MainProductName string    `json:"main_name"`
}

type VariantDetailDTO struct {
	VariantDTO
	Category  []CategoriesDTO
	Attribute []AttributesDTO
}

func NewVariant(name, desc string, price float64, mainId ulid.ULID) (Variant, error) {
	id := ulid.Make()
	mainProduct := Product{
		ProductID: mainId,
	}
	return Variant{
		VariantID:   id,
		MainProduct: mainProduct,
		Name:        name,
		Description: desc,
		Price:       price,
	}, nil
}
