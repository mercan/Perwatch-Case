package models

type StockRequestParams struct {
	FormID  string `json:"form_id"`
	StockID string `json:"stock_id,omitempty"`
}

type CreateStockRequest struct {
	Fields []struct {
		Name  string      `json:"name"`
		Value interface{} `json:"value"`
	} `json:"fields"`
}

type UpdateStockValueRequest struct {
	Value interface{} `json:"value"`
}
