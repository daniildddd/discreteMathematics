package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
)

func main() {
	var m int
	fmt.Print("Введите количество проверочных битов (m): ")
	fmt.Scan(&m)

	// k может быть равно только (2^m - 1) - m, где m - количество проверочных битов
	k := int(math.Pow(2, float64(m)) - 1 - float64(m))

	infoBits := generateCombination(k)
	table := generateHammingTable(m)
	encoded := hammingEncode(infoBits, table)

	// Внести ошибку
	errorPos := rand.Intn(len(encoded))
	encodedWithError := make([]int, len(encoded))
	copy(encodedWithError, encoded)
	encodedWithError[errorPos] = 1 - encodedWithError[errorPos]

	// Найти синдром и исправить ошибку
	syndrome := findSyndrome(encodedWithError, table)
	errorPosition := findErrorPosition(syndrome, table)
	if errorPosition != -1 {
		encodedWithError[errorPosition] = 1 - encodedWithError[errorPosition]
	}

	// Создаем или перезаписываем файл
	file, err := os.Create("hamming_result.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writeLine := func(format string, a ...any) {
		fmt.Fprintf(file, format+"\n", a...)
	}

	writeLine("Сгенерированная информационная комбинация: %v", infoBits)
	writeLine("Таблица Хемминга (%dx%d):", len(table), len(table[0]))
	for i := 0; i < len(table); i++ {
		writeLine("P%d: %v", i+1, table[i])
	}
	writeLine("Закодированная комбинация: %v", encoded)
	writeLine("Внесена ошибка в позицию %d: %v", errorPos+1, encodedWithError)
	writeLine("Синдром: %v", syndrome)
	if errorPosition != -1 {
		writeLine("Ошибка найдена в позиции %d", errorPosition+1)
		writeLine("Исправленная комбинация: %v", encodedWithError)
	} else {
		writeLine("Ошибок не обнаружено")
	}

	if compareVectors(encoded, encodedWithError) {
		writeLine("Результат совпадает с исходной комбинацией")
	} else {
		writeLine("Результат НЕ совпадает с исходной комбинацией")
	}

	fmt.Println("Результат записан в файл hamming_result.txt")
}

func generateCombination(size int) []int {
	combination := make([]int, size)
	for i := 0; i < size; i++ {
		combination[i] = rand.Intn(2)
	}
	return combination
}

// Создает таблицу проверочных битов Хэмминга
func generateHammingTable(m int) [][]int {
	n := 1<<m - 1
	table := make([][]int, m)
	for i := 0; i < m; i++ {
		table[i] = make([]int, n)
		pattern := 1 << i
		for j := 0; j < n; j++ {
			if ((j + 1) & pattern) != 0 {
				table[i][j] = 1
			}
		}
	}
	return table
}

func hammingEncode(infoBits []int, table [][]int) []int {
	m := len(table)
	n := len(table[0])
	k := len(infoBits)

	if k != n-m {
		panic("Несоответствие размеров")
	}

	encoded := make([]int, n)
	infoIndex := 0
	for i := 0; i < n; i++ {
		if !isPowerOfTwo(i + 1) {
			encoded[i] = infoBits[infoIndex]
			infoIndex++
		}
	}

	for i := 0; i < m; i++ {
		parity := 0
		for j := 0; j < n; j++ {
			if table[i][j] == 1 {
				parity ^= encoded[j]
			}
		}
		pos := (1 << i) - 1
		if pos < n {
			encoded[pos] = parity
		}
	}

	return encoded
}

func findSyndrome(encoded []int, table [][]int) []int {
	m := len(table)
	syndrome := make([]int, m)
	for i := 0; i < m; i++ {
		parity := 0
		for j := 0; j < len(encoded); j++ {
			if table[i][j] == 1 {
				parity ^= encoded[j]
			}
		}
		syndrome[i] = parity
	}
	return syndrome
}

func findErrorPosition(syndrome []int, table [][]int) int {
	if isZero(syndrome) {
		return -1
	}
	for j := 0; j < len(table[0]); j++ {
		match := true
		for i := 0; i < len(syndrome); i++ {
			if table[i][j] != syndrome[i] {
				match = false
				break
			}
		}
		if match {
			return j
		}
	}
	return -1
}

func isPowerOfTwo(n int) bool {
	return n > 0 && (n&(n-1)) == 0
}

func isZero(vector []int) bool {
	for _, v := range vector {
		if v != 0 {
			return false
		}
	}
	return true
}

func compareVectors(v1, v2 []int) bool {
	if len(v1) != len(v2) {
		return false
	}
	for i := 0; i < len(v1); i++ {
		if v1[i] != v2[i] {
			return false
		}
	}
	return true
}
