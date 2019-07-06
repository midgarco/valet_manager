package pagination

// Paging determines how to page the data
type Paging struct {
	Offset  int
	Limit   int
	OrderBy []Order
}

// Order contains the field and sort direction
type Order struct {
	Field     string
	Direction string
}
