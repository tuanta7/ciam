package domain

import (
	"time"

	"github.com/go-jose/go-jose/v4"
)

type Key struct {
	Kid       string    `json:"kid"`
	Alg       string    `json:"alg"`
	Usage     string    `json:"use"`
	Value     any       `json:"key"`
	CreatedBy string    `json:"created_by"`
	UpdatedBy string    `json:"updated_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (k Key) ID() string {
	return k.Kid
}

func (k Key) Algorithm() jose.SignatureAlgorithm {
	return jose.SignatureAlgorithm(k.Alg)
}

func (k Key) Use() string {
	return k.Usage
}

func (k Key) Key() any {
	return k.Value
}
