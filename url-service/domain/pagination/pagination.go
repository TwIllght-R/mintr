package pagination

type Pagination struct {
	Page     int    // current page number
	PageSize int    // number of items per page
	Sort     string // sorting order (e.g., "asc" or "desc")
	SortBy   string // field to sort by (e.g., "created_at")
}

type PaginatedResponse[T any] struct {
	Items    []T    // the current page of results
	Page     int    // current page number
	PageSize int    // number of items per page
	Sort     string // sorting order (e.g., "asc" or "desc")
	SortBy   string // field to sort by (e.g., "created_at")
	Total    int64  // total number of items available
}

type CursorPagination struct {
	Limit  int     // number of items to return
	Cursor *string // nil or empty = start from the beginning
}

type PaginatedCursorResponse[T any] struct {
	Items      []T     // the current page of results
	NextCursor *string // nil means “no more pages”
}
