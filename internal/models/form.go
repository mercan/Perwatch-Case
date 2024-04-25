package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Form struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	UserID    primitive.ObjectID `json:"user_id" bson:"user_id"`
	Name      string             `json:"name" bson:"name"`
	IsDeleted bool               `json:"is_deleted" bson:"is_deleted"`

	Fields []Field `json:"fields" bson:"fields,omitempty"`
	Stock  []Stock `json:"stocks" bson:"stocks,omitempty"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

type Stock struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Name      string             `json:"name" bson:"name"`
	Value     interface{}        `json:"value" bson:"value"`
	IsDeleted bool               `json:"is_deleted" bson:"is_deleted"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

type Field struct {
	ID                   primitive.ObjectID `json:"id" bson:"_id"`
	Name                 string             `json:"name" bson:"name"`
	Type                 string             `json:"type" bson:"type"`
	Values               []string           `json:"values,omitempty" bson:"values,omitempty,omitempty"`
	MinValue             int                `json:"min_value,omitempty" bson:"min_value,omitempty"`
	MaxValue             int                `json:"max_value,omitempty" bson:"max_value,omitempty"`
	MinValueDecimal      float64            `json:"min_value_decimal,omitempty" bson:"min_value_decimal,omitempty"`
	MaxValueDecimal      float64            `json:"max_value_decimal,omitempty" bson:"max_value_decimal,omitempty"`
	MinLength            int                `json:"min_length,omitempty" bson:"min_length,omitempty"`
	MaxLength            int                `json:"max_length,omitempty" bson:"max_length,omitempty"`
	DefaultString        string             `json:"default_string,omitempty" bson:"default_string,omitempty"`
	DefaultNumber        int                `json:"default_number,omitempty" bson:"default_number,omitempty"`
	DefaultNumberDecimal float64            `json:"default_number_decimal,omitempty" bson:"default_number_decimal,omitempty"`
	Sort                 int                `json:"sort,omitempty" bson:"sort,omitempty"`
	IsDeleted            bool               `json:"is_deleted,omitempty" bson:"is_deleted"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

func NewForm() Form {
	return Form{
		ID:        primitive.NewObjectID(),
		IsDeleted: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
