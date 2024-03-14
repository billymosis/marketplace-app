package model

import (
	"net/http"
	"strconv"
)

type Product struct {
	Id             uint     `json:"id"`
	Name           string   `json:"name"`
	Price          uint     `json:"price"`
	ImageUrl       string   `json:"imageUrl"`
	Stock          uint     `json:"stock"`
	Condition      string   `json:"condition"`
	Tags           []string `json:"tags"`
	IsPurchaseable bool     `json:"isPurchaseable"`
	UserId         uint     `json:"userId"`
}

type ProductQueryParams struct {
	UserOnly       *bool
	Limit          *int
	Offset         *int
	Tags           []*string
	Condition      *string
	ShowEmptyStock *bool
	MaxPrice       *int
	MinPrice       *int
	SortBy         *string
	OrderBy        *string
	Search         *string
}

func ParseProductQueryParams(r *http.Request) (*ProductQueryParams, error) {
	return nil, nil
	query := r.URL.Query()
	var params ProductQueryParams

	userOnlyStr := query.Get("userOnly")
	if userOnlyStr != "" {
		userOnly, err := strconv.ParseBool(userOnlyStr)
		if err != nil {
			return nil, err
		}
		params.UserOnly = &userOnly
	}

	limitStr := query.Get("limit")
	if limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			return nil, err
		}
		params.Limit = &limit
	}

	offsetStr := query.Get("offset")
	if offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			return nil, err
		}
		params.Offset = &offset
	}

	tags := query["tags"]
	if len(tags) > 0 {
		params.Tags = make([]*string, len(tags))
		for i, tag := range tags {
			params.Tags[i] = &tag
		}
	}

	condition := query.Get("condition")
	if condition != "" {
		params.Condition = &condition
	}

	showEmptyStockStr := query.Get("showEmptyStock")
	if showEmptyStockStr != "" {
		showEmptyStock, err := strconv.ParseBool(showEmptyStockStr)
		if err != nil {
			return nil, err
		}
		params.ShowEmptyStock = &showEmptyStock
	}

	maxPriceStr := query.Get("maxPrice")
	if maxPriceStr != "" {
		maxPrice, err := strconv.Atoi(maxPriceStr)
		if err != nil {
			return nil, err
		}
		params.MaxPrice = &maxPrice
	}

	minPriceStr := query.Get("minPrice")
	if minPriceStr != "" {
		minPrice, err := strconv.Atoi(minPriceStr)
		if err != nil {
			return nil, err
		}
		params.MinPrice = &minPrice
	}

	sortBy := query.Get("sortBy")
	if sortBy != "" {
		params.SortBy = &sortBy
	}

	orderBy := query.Get("orderBy")
	if orderBy != "" {
		params.OrderBy = &orderBy
	}

	search := query.Get("search")
	if search != "" {
		params.Search = &search
	}

	return &params, nil
}
