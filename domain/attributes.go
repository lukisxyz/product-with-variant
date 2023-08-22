package domain

import (
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
	"gopkg.in/guregu/null.v4"
)

type Attribute struct {
	AttributeID ulid.ULID
	Name        string
	CreatedAt   time.Time
	UpdatedAt   null.Time
	DeletedAt   null.Time
}

type AttributesDTO struct {
	ID   ulid.ULID `json:"id"`
	Name string    `json:"name"`
}

func NewAttribute(name string) (Attribute, error) {
	id := ulid.Make()
	return Attribute{
		AttributeID: id,
		Name:        strings.ToLower(name),
		CreatedAt:   time.Now(),
	}, nil
}
