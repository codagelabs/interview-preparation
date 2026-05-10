package models

type SearchRequest struct {
	SearchText string `form:"search_text"  binding:"required" example:"search_text"`
}

type SearchSuggestionRequest struct {
	SearchText string `form:"search_text"  binding:"required" example:"search_text"`
}

type SearchSuggestionResponse struct {
	Tittle           string `json:"tittle"`
	ShortDescription string `json:"short_description"`
	ImageUrl         string `json:"image_url"`
	DestinationUrl   string `json:"destination_url"`
}
