package dto

type (
	ViewPagination struct {
		Limit  int64 `json:"limit"`
		Offset int64 `json:"offset"`
		Total  int64 `json:"total"`
	}
)
