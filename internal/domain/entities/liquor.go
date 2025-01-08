package entities

type (
	Liquor struct {
		LiquorID    string `json:"liquor_id"`
		GTIN        string `json:"gtin"`
		Name        string `json:"name"`
		Brand       string `json:"brand"`
		Description string `json:"description"`
		ImageURL    string `json:"image_url"`
		CategoryID  string `json:"category_id"`
		Content     string `json:"content"`
		Country     string `json:"country"`
	}
)
