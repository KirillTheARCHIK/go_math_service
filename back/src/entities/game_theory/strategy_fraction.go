package game_theory

type StrategyFraction struct {
	Strategy        Strategy `json:"strategy" binding:"required"`
	PotencialProfit float64  `json:"potencialProfit" binding:"required"`
	Fraction        float64  `json:"fraction" binding:"required"`
}
