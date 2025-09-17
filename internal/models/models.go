package models

// Material представляет изоляционный материал
type Material struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	PricePerM2  float64 `json:"price_per_m2"`
	Lambda      float64 `json:"lambda"`
	Description string  `json:"description"`
	ImageURL    string  `json:"image_url"`
}

// CartItem представляет элемент в заявке
type CartItem struct {
	MaterialID int     `json:"material_id"`
	Area       float64 `json:"area"`
}

// Cart представляет заявку с расчетами
type Cart struct {
	ID           int        `json:"id"`
	Items        []CartItem `json:"items"`
	TotalSavings float64    `json:"total_savings"`
}
