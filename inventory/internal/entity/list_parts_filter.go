package entity

import (
	"github.com/paincake00/microservices-go/inventory/internal/entity/enum"
)

type ListPartsFilter struct {
	Uuids                 []string
	Names                 []string
	Categories            []enum.Category
	ManufacturerCountries []string
	Tags                  []string
}

func NewListPartsFilter(
	uuids []string,
	names []string,
	categories []enum.Category,
	manufacturerCountries []string,
	tags []string,
) ListPartsFilter {
	return ListPartsFilter{
		Uuids:                 uuids,
		Names:                 names,
		Categories:            categories,
		ManufacturerCountries: manufacturerCountries,
		Tags:                  tags,
	}
}
