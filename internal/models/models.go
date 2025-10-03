package models

type Material struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	PricePerM2  float64 `json:"price_per_m2"`
	Lambda      float64 `json:"lambda"`
	Description string  `json:"description"`
	ImageURL    string  `json:"image_url"`
}

type CartMaterial struct {
	MaterialID int     `json:"material_id"`
	Area       float64 `json:"area"`
}

type ApplicationMaterials struct {
	ID           int            `json:"id"`
	Items        []CartMaterial `json:"items"`
	TotalSavings float64        `json:"total_savings"`
	IndoorTemp   float64        `json:"indoor_temp"`
	OutdoorTemp  float64        `json:"outdoor_temp"`
	TotalArea    float64        `json:"total_area"`
}
