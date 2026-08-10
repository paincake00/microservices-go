package enum

// Category for describing part kind
type Category string

const (
	CategoryUnspecified Category = "CATEGORY_UNSPECIFIED"
	CategoryEngine      Category = "CATEGORY_ENGINE"
	CategoryFuel        Category = "CATEGORY_FUEL"
	CategoryPorthole    Category = "CATEGORY_PORTHOLE"
	CategoryWing        Category = "CATEGORY_WING"
)

// CategoryAll - all values of Category Enum.
var (
	CategoryAll = []string{
		"CATEGORY_UNSPECIFIED",
		"CATEGORY_ENGINE",
		"CATEGORY_FUEL",
		"CATEGORY_PORTHOLE",
		"CATEGORY_WING",
	}
)
