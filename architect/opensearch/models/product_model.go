package models

type Product struct {
	ProductID       string  `json:"product_id"`
	ProductName     string  `json:"product_name"`
	ProductCategory string  `json:"product_category"`
	ProductPrise    int     `json:"product_prise"`
	IsAvailable     bool    `json:"is_available"`
	Suggest         Suggest `json:"suggest"`
}

type Suggest struct {
	Input  []string `json:"input"`
	Weight int      `json:"weight"`
}
