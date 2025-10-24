package model

type Product struct {
	Model

	Name          string  `bson:"name" json:"name" validate:"required"`
	ProductId     string  `bson:"productId" json:"productId"`
	Price         float64 `bson:"price" json:"price" validate:"required"`
	StockQuantity int64   `bson:"stockQuantity" json:"stockQuantity" validate:"required"`
	Weight        float64 `bson:"weight" json:"weight" validate:"required"`
}
