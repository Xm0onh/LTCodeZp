package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	Modulus = 127
	Dim     = 11
)

func generateCoefficients() []int {
	coefficients := make([]int, Dim)
	for i := range coefficients {
		coefficients[i] = rand.Intn(Modulus)
	}
	return coefficients
}

func encodeSymbol(coeffs, message []int) int {
	encoded := 0
	for i, coeff := range coeffs {
		encoded += coeff * message[i]
		encoded %= Modulus
	}
	return encoded
}

func solveSystem(matrix [][]int, results []int) ([]int, error) {
	n := len(matrix)
	solution := make([]int, n)

	for i := 0; i < n; i++ {
		if matrix[i][i] == 0 { // Pivot is zero, matrix might be singular
			return nil, fmt.Errorf("matrix is singular")
		}

		invPivot := modInverse(matrix[i][i], Modulus)
		for j := 0; j < n; j++ {
			matrix[i][j] = (matrix[i][j] * invPivot) % Modulus
			matrix[i][j] = (matrix[i][j] + Modulus) % Modulus
		}
		results[i] = (results[i] * invPivot) % Modulus
		results[i] = (results[i] + Modulus) % Modulus

		for j := 0; j < n; j++ {
			if i != j {
				factor := matrix[j][i]
				for k := 0; k < n; k++ {
					matrix[j][k] = (matrix[j][k] - factor*matrix[i][k]) % Modulus
					matrix[j][k] = (matrix[j][k] + Modulus) % Modulus
				}
				results[j] = (results[j] - factor*results[i]) % Modulus
				results[j] = (results[j] + Modulus) % Modulus
			}
		}
	}
	for i := 0; i < n; i++ {
		solution[i] = results[i]
	}
	return solution, nil
}

func modInverse(a, m int) int {
	m0, x0, x1 := m, 0, 1
	if m == 1 {
		return 0
	}
	for a > 1 {
		q := a / m
		m, a = a%m, m
		x0, x1 = x1-q*x0, x0
	}
	if x1 < 0 {
		x1 += m0
	}
	return x1
}

func main() {
	rand.Seed(time.Now().UnixNano())

	message := []int{2, 4, 1, 123, 12, 5, 1, 23, 5, 6, 1}
	fmt.Println("Original Message:", message)

	coeffMatrix := make([][]int, Dim)
	encoded := make([]int, Dim)

	attempts := 0
	maxAttempts := 10 // Limit attempts to prevent infinite loop
	var err error

	for attempts < maxAttempts {
		for i := range coeffMatrix {
			coeffMatrix[i] = generateCoefficients()
			encoded[i] = encodeSymbol(coeffMatrix[i], message)
		}

		_, err = solveSystem(coeffMatrix, encoded)
		if err == nil {
			break
		}
		attempts++
	}

	if err != nil {
		fmt.Println("Failed to decode message after multiple attempts:", err)
		return
	}

	fmt.Println("Encoded Symbols:", encoded)
	fmt.Println("Coefficient Matrix:", coeffMatrix)

	decoded, _ := solveSystem(coeffMatrix, encoded)
	fmt.Println("Decoded Message:", decoded)
}
