package dto

type MultiID struct {
	Ids []string `json:"ids" binding:"required"`
}
