package game_theory

import (
	"cmp"
	"errors"
	"slices"
)

type GameTheoryMatrix struct {
	Outcomes   []Outcome  `json:"outcomes" binding:"required"`
	Strategies []Strategy `json:"strategies" binding:"required"`
}

func (matrix *GameTheoryMatrix) Solve(firstStrategies int) (distribution []StrategyFraction, err error) {
	for _, strategy := range matrix.Strategies {
		if len(strategy.Profits) != len(matrix.Outcomes) {
			return nil, errors.New("количество выгод стратегий должно быть равно количеству исходов")
		}
	}
	if firstStrategies < 1 || firstStrategies > len(matrix.Strategies) {
		return nil, errors.New("количество первых стратегий не может быть меньше 1 и больше количества стратегий")
	}
	//Критерий Байеса
	for _, strategy := range matrix.Strategies {
		var potencialProfit float64 = 0
		for index, profit := range strategy.Profits {
			potencialProfit += profit * matrix.Outcomes[index].Chance
		}
		distribution = append(distribution, StrategyFraction{
			Strategy:        strategy,
			PotencialProfit: potencialProfit,
		})
	}
	slices.SortFunc(distribution, func(a, b StrategyFraction) int {
		return cmp.Compare(a.PotencialProfit, b.PotencialProfit)
	})
	distribution[0].Fraction = 1
	return distribution, nil
}
