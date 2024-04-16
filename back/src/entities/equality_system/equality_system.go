package equality_system

import (
	"errors"
	"fmt"
)

func Solve(A [][]float64, B []float64) (X []float64, err error) {
	var n int = len(A)
	for _, row := range A {
		if len(row) != n {
			return nil, errors.New("число неизвестных в каждом уравнении должно быть равно числу неизвестных в системе")

		}
	}
	if len(B) != n {
		return nil, errors.New("число свободных членов должно быть равно числу неизвестных")
	}
	PrintSystem(A, B)
	//Прямой ход
	for col := 0; col < n-1; col++ {
		for row := col + 1; row < n; row++ {
			if A[row][col] == 0 {
				continue
			}
			var coef = A[row][col] / A[col][col]
			// fmt.Println(coef)
			for c := col; c < n; c++ {
				A[row][c] -= A[col][c] * coef
			}
			B[row] -= B[col] * coef
		}
	}
	PrintSystem(A, B)
	//Обратный ход
	for col := n - 1; col >= 1; col-- {
		for row := col - 1; row >= 0; row-- {
			if A[row][col] == 0 {
				continue
			}
			var coef = A[row][col] / A[col][col]
			// fmt.Println(coef)
			A[row][col] -= A[col][col] * coef
			B[row] -= B[col] * coef
		}
	}
	PrintSystem(A, B)
	//Решение
	for col := 0; col < n; col++ {
		X = append(X, B[col]/A[col][col])
	}
	return X, nil
}

func PrintSystem(A [][]float64, B []float64) {
	var n int = len(A)
	for row := 0; row < n; row++ {
		for col := 0; col < n; col++ {
			fmt.Printf("%6.2f ", A[row][col])
		}
		fmt.Printf("| %6.2f\n", B[row])
	}
	fmt.Println()
}
