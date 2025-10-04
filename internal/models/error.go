package models

type Error struct {
	Error   string `json:"error" example:"INVALID_INPUT"`
	Message string `json:"message" example:"The provided input data is invalid"`
	Details string `json:"details,omitempty" example:"Product ID must be a positive integer"`
}
