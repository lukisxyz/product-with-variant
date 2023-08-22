package domain

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type Product struct {
	ProductID    ulid.ULID
	Name         string
	Description  string
	Price        float64
	ImagePreview []byte
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    time.Time
}
type ProductDTO struct {
	ID          ulid.ULID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"desc"`
	Price       float64   `json:"price"`
	Image       []byte    `json:"image"`
}

type ProductDetailDTO struct {
	ProductDTO
	Category  []CategoriesDTO
	Attribute []AttributesDTO
}

func NewProduct(name, desc string, price float64) (Product, error) {
	id := ulid.Make()
	return Product{
		ProductID:   id,
		Name:        name,
		Description: desc,
		Price:       price,
	}, nil
}
