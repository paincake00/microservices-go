package entity

type Category int32

const (
	CATEGORY_UNSPECIFIED Category = iota
	CATEGORY_ENGINE
	CATEGORY_FUEL
	CATEGORY_PORTHOLE
	CATEGORY_WING
)

// Enum value maps for Category.
var (
	CategoryName = map[int32]string{
		0: "CATEGORY_UNSPECIFIED",
		1: "CATEGORY_ENGINE",
		2: "CATEGORY_FUEL",
		3: "CATEGORY_PORTHOLE",
		4: "CATEGORY_WING",
	}
	CategoryValue = map[string]int32{
		"CATEGORY_UNSPECIFIED": 0,
		"CATEGORY_ENGINE":      1,
		"CATEGORY_FUEL":        2,
		"CATEGORY_PORTHOLE":    3,
		"CATEGORY_WING":        4,
	}
)
