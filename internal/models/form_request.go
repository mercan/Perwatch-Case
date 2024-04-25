package models

type CreateFormFieldRequest struct {
	FieldType string `json:"field_type"`
	FieldName string `json:"field_name"`

	Values          []string `json:"values"`
	MinValue        int      `json:"min_value"`
	MaxValue        int      `json:"max_value"`
	MinValueDecimal float64  `json:"min_value_decimal"`
	MaxValueDecimal float64  `json:"max_value_decimal"`
	MinLength       int      `json:"min_length"`
	MaxLength       int      `json:"max_length"`

	DefaultString        string  `json:"default_string,omitempty" bson:"default_string,omitempty"`
	DefaultNumber        int     `json:"default_number,omitempty" bson:"default_number,omitempty"`
	DefaultNumberDecimal float64 `json:"default_number_decimal,omitempty" bson:"default_number_decimal,omitempty"`

	Sort int `json:"sort,omitempty" bson:"sort,omitempty"`
}

type FormFieldRequestParams struct {
	FormID  string
	FieldID string
}

type FormNameRequestBody struct {
	FormName string `json:"form_name"`
}

type FormIDRequestBody struct {
	FormID string `json:"form_id"`
}

type GetFormsRequest struct {
	Page int `json:"page"`
}
