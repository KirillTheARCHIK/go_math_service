package game_theory

type Outcome struct {
	Name   string  `json:"name" binding:"required"`
	Chance float64 `json:"chance" binding:"required"`
}
