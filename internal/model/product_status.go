package model

type ProductStatus string

const (
	ProductStatusActive   ProductStatus = "ACTIVE"
	ProductStatusInactive ProductStatus = "INACTIVE"
)

func (ps ProductStatus) IsValid() bool {
	return ps == ProductStatusActive || ps == ProductStatusInactive
}
