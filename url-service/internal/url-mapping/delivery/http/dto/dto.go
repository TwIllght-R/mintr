package dto

import (
	"url-service/domain/entities"
	"url-service/domain/pagination"
)

type GenerateShortCodeRequest struct {
	OriginalURL string `json:"original_url"`
	Title       string `json:"title,omitempty"`
}

func (d *GenerateShortCodeRequest) ToEntity() entities.URLMapping {
	return entities.URLMapping{
		OriginalURL: d.OriginalURL,
		Title:       d.Title,
	}
}

type GetAllMappingURLRequest struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Sort     string `form:"sort"`
	SortBy   string `form:"sort_by"`
}

func (d *GetAllMappingURLRequest) ToEntity() pagination.Pagination {
	return pagination.Pagination{
		Page:     d.Page,
		PageSize: d.PageSize,
		Sort:     d.Sort,
		SortBy:   d.SortBy,
	}
}

type UpdateURLMappingRequest struct {
	OriginalURL string `json:"original_url,omitempty"`
	Title       string `json:"title,omitempty"`
}

func (d *UpdateURLMappingRequest) ToEntity() entities.URLMapping {
	return entities.URLMapping{
		OriginalURL: d.OriginalURL,
		Title:       d.Title,
	}
}

type URLMappingResponse struct {
	UUID        string `json:"uuid"`
	Title       string `json:"title"`
	ShortCode   string `json:"short_code"`
	OriginalURL string `json:"original_url"`
}

func ToResponse(in entities.URLMapping) URLMappingResponse {
	return URLMappingResponse{
		UUID:        in.UUID,
		Title:       in.Title,
		ShortCode:   in.ShortCode,
		OriginalURL: in.OriginalURL,
	}
}

type URLMappingPaginatedResponse struct {
	Items []URLMappingResponse `json:"items"`
	Meta  struct {
		Page     int    `json:"page"`
		PageSize int    `json:"page_size"`
		Sort     string `json:"sort"`
		SortBy   string `json:"sort_by"`
		Total    int64  `json:"total"`
	} `json:"meta"`
}

func ToResponsePaginated(in []entities.URLMapping, page, pageSize int, sort, sortBy string, total int64) URLMappingPaginatedResponse {
	items := make([]URLMappingResponse, len(in))
	for i, item := range in {
		items[i] = ToResponse(item)
	}
	return URLMappingPaginatedResponse{
		Items: items,
		Meta: struct {
			Page     int    `json:"page"`
			PageSize int    `json:"page_size"`
			Sort     string `json:"sort"`
			SortBy   string `json:"sort_by"`
			Total    int64  `json:"total"`
		}{
			Page:     page,
			PageSize: pageSize,
			Sort:     sort,
			SortBy:   sortBy,
			Total:    total,
		},
	}
}
