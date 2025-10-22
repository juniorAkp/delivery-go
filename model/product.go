package model

type Product struct {
	Model

	Name          string  `bson:"name" json:"name"`
	Price         float64 `bson:"price" json:"price"`
	StockQuantity int64   `bson:"stockQuantity" json:"stockQuantity"`
	Weight        float64 `bson:"weight" json:"weight"`
}
