package entity

import (
	"slices"

	"github.com/paincake00/microservices-go/inventory/internal/entity/enum"
)

type ListPartsFilter struct {
	Uuids                 map[string]struct{}
	Names                 map[string]struct{}
	Categories            map[enum.Category]struct{}
	ManufacturerCountries map[string]struct{}
	Tags                  map[string]struct{}
}

func NewListPartsFilter(
	uuids []string,
	names []string,
	categories []enum.Category,
	manufacturerCountries []string,
	tags []string,
) ListPartsFilter {
	return ListPartsFilter{
		Uuids:                 createSet(uuids),
		Names:                 createSet(names),
		Categories:            createSet(categories),
		ManufacturerCountries: createSet(manufacturerCountries),
		Tags:                  createSet(tags),
	}
}

func (f *ListPartsFilter) IsEmpty() bool {
	return len(f.Uuids) == 0 && len(f.Names) == 0 && len(f.Categories) == 0 && len(f.ManufacturerCountries) == 0 && len(f.Tags) == 0
}

func (f *ListPartsFilter) Match(part *Part) bool {
	if _, ok := f.Uuids[part.Uuid]; len(f.Uuids) != 0 && !ok {
		return false
	}
	if _, ok := f.Names[part.Name]; len(f.Names) != 0 && !ok {
		return false
	}
	if _, ok := f.Categories[part.Category]; len(f.Categories) != 0 && !ok {
		return false
	}
	if _, ok := f.ManufacturerCountries[part.Manufacturer.Country]; len(f.ManufacturerCountries) != 0 && !ok {
		return false
	}
	if len(f.Tags) != 0 && !slices.ContainsFunc(
		part.Tags, func(s string) bool {
			_, ok := f.Tags[s]
			return ok
		},
	) {
		return false
	}
	return true
}

func createSet[T comparable](items []T) map[T]struct{} {
	var set map[T]struct{}
	if len(items) != 0 {
		set = make(map[T]struct{}, len(items))
		for _, item := range items {
			set[item] = struct{}{}
		}
	}
	return set
}
