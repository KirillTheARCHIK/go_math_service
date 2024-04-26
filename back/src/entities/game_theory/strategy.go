package game_theory

type Strategy struct {
	Name    string    `json:"name" binding:"required"`
	Profits []float64 `json:"profits" binding:"required"`
}
