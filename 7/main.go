package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

// Предварительная обработка текста
func preprocessText(text string) string {
	text = strings.ToLower(text)
	var result strings.Builder
	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsSpace(r) {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// Простое шифрование (шифр Цезаря)
func simpleCipher(text string, key int) string {
	var result strings.Builder

	for _, r := range text {
		if r >= 'a' && r <= 'z' {
			shifted := ((int(r)-'a'+key)%26 + 26) % 26
			result.WriteRune(rune('a' + shifted))
		} else if r >= 'а' && r <= 'я' {
			shifted := ((int(r)-'а'+key)%32 + 32) % 32
			result.WriteRune(rune('а' + shifted))
		} else if r == 'ё' {
			baseChar := 'е'
			shifted := ((int(baseChar)-'а'+key)%32 + 32) % 32
			shiftedChar := rune('а' + shifted)
			if shiftedChar == 'е' {
				result.WriteRune('ё')
			} else {
				result.WriteRune(shiftedChar)
			}
		} else {
			result.WriteRune(r)
		}
	}

	return result.String()
}

// Простое дешифрование
func simpleDecipher(text string, key int) string {
	return simpleCipher(text, -key)
}

// Получить позицию буквы в алфавите
func getLetterPosition(r rune) int {
	if r >= 'a' && r <= 'z' {
		return int(r - 'a')
	} else if r >= 'а' && r <= 'я' {
		return int(r - 'а')
	} else if r == 'ё' {
		return int('е' - 'а')
	}
	return 0
}

// Сдвинуть букву на указанное количество позиций
func shiftLetter(r rune, shift int) rune {
	if r >= 'a' && r <= 'z' {
		shifted := ((int(r)-'a'+shift)%26 + 26) % 26
		return rune('a' + shifted)
	} else if r >= 'а' && r <= 'я' {
		shifted := ((int(r)-'а'+shift)%32 + 32) % 32
		return rune('а' + shifted)
	} else if r == 'ё' {
		shifted := ((int('е')-'а'+shift)%32 + 32) % 32
		result := rune('а' + shifted)
		if result == 'е' {
			return 'ё'
		}
		return result
	}
	return r
}

// Подготовка ключа для Виженера
func prepareVigenereKey(key string) string {
	key = strings.ToLower(key)
	var filtered strings.Builder
	for _, r := range key {
		if unicode.IsLetter(r) {
			filtered.WriteRune(r)
		}
	}
	result := filtered.String()
	if len(result) == 0 {
		return "ключ"
	}
	return result
}

// Сложное шифрование (шифр Виженера)
func complexCipher(text string, key string) string {
	key = prepareVigenereKey(key)

	var result strings.Builder
	keyIndex := 0

	for _, r := range text {
		if unicode.IsLetter(r) {
			// Получаем символ ключа
			keyChar := rune(key[keyIndex%len(key)])
			// Получаем сдвиг
			shift := getLetterPosition(keyChar)
			// Применяем сдвиг
			shifted := shiftLetter(r, shift)
			result.WriteRune(shifted)
			keyIndex++
		} else {
			result.WriteRune(r)
		}
	}

	return result.String()
}

// Сложное дешифрование
func complexDecipher(text string, key string) string {
	key = prepareVigenereKey(key)

	var result strings.Builder
	keyIndex := 0

	for _, r := range text {
		if unicode.IsLetter(r) {
			// Получаем символ ключа
			keyChar := rune(key[keyIndex%len(key)])
			// Получаем сдвиг (отрицательный для дешифрования)
			shift := -getLetterPosition(keyChar)
			// Применяем сдвиг
			shifted := shiftLetter(r, shift)
			result.WriteRune(shifted)
			keyIndex++
		} else {
			result.WriteRune(r)
		}
	}

	return result.String()
}

func main() {
	inputFile := "input.txt"

	data, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Printf("Ошибка чтения файла %s: %v\n", inputFile, err)
		fmt.Println("Убедитесь, что файл input.txt находится в той же папке, что и программа.")
		return
	}
	originalText := string(data)

	fmt.Println("Файл input.txt успешно прочитан!")
	fmt.Printf("Размер текста: %d символов\n\n", len(originalText))

	processedText := preprocessText(originalText)
	fmt.Println("Текст обработан (приведен к нижнему регистру, удалены знаки препинания)")

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("\nВведите ключ для простого шифрования (целое число, например 3): ")
	scanner.Scan()
	var simpleKey int
	fmt.Sscanf(scanner.Text(), "%d", &simpleKey)

	fmt.Print("Введите ключ для сложного шифрования (слово, например 'секрет'): ")
	scanner.Scan()
	complexKeyInput := scanner.Text()

	if complexKeyInput == "" {
		fmt.Println("Ключ не может быть пустым! Использую 'ключ' по умолчанию.")
		complexKeyInput = "ключ"
	}

	// Подготавливаем ключ один раз и проверяем
	complexKey := prepareVigenereKey(complexKeyInput)
	if complexKey != complexKeyInput && complexKey != strings.ToLower(complexKeyInput) {
		fmt.Printf("ПРЕДУПРЕЖДЕНИЕ: Ключ '%s' содержит недопустимые символы!\n", complexKeyInput)
		fmt.Printf("Используется ключ: '%s'\n", complexKey)
	}

	fmt.Println("\nВыполняется шифрование...")

	simpleCiphered := simpleCipher(processedText, simpleKey)
	simpleDeciphered := simpleDecipher(simpleCiphered, simpleKey)

	complexCiphered := complexCipher(processedText, complexKey)
	complexDeciphered := complexDecipher(complexCiphered, complexKey)

	// Проверка: выводим первые 50 символов
	fmt.Printf("\nПроверка шифрования (первые 50 символов):\n")
	if len(processedText) > 50 {
		fmt.Printf("Исходный:      %s...\n", processedText[:50])
		fmt.Printf("Зашифрованный: %s...\n", complexCiphered[:50])
	}

	outputFile := "output.txt"
	f, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("Ошибка создания файла: %v\n", err)
		return
	}
	defer f.Close()

	writer := bufio.NewWriter(f)

	fmt.Fprintln(writer, "=== ИСХОДНЫЙ ТЕКСТ ===")
	fmt.Fprintln(writer, originalText)
	fmt.Fprintln(writer, "\n=== ОБРАБОТАННЫЙ ТЕКСТ (нижний регистр, без знаков препинания) ===")
	fmt.Fprintln(writer, processedText)
	fmt.Fprintln(writer, "\n=== ПРОСТОЕ ШИФРОВАНИЕ (шифр Цезаря, ключ:", simpleKey, ") ===")
	fmt.Fprintln(writer, simpleCiphered)
	fmt.Fprintln(writer, "\n=== ПРОСТОЕ ДЕШИФРОВАНИЕ ===")
	fmt.Fprintln(writer, simpleDeciphered)
	fmt.Fprintln(writer, "\n=== СЛОЖНОЕ ШИФРОВАНИЕ (шифр Виженера, ключ:", complexKey, ") ===")
	fmt.Fprintln(writer, complexCiphered)
	fmt.Fprintln(writer, "\n=== СЛОЖНОЕ ДЕШИФРОВАНИЕ ===")
	fmt.Fprintln(writer, complexDeciphered)

	fmt.Fprintln(writer, "\n=== ПРОВЕРКА КОРРЕКТНОСТИ ===")
	fmt.Fprintf(writer, "Простое дешифрование совпадает с обработанным: %v\n", simpleDeciphered == processedText)
	fmt.Fprintf(writer, "Сложное дешифрование совпадает с обработанным: %v\n", complexDeciphered == processedText)

	writer.Flush()

	fmt.Printf("\n Все операции выполнены успешно!")
	fmt.Printf("\n Результаты записаны в файл %s\n\n", outputFile)
	fmt.Println("Проверка корректности:")
	fmt.Printf("   Простое дешифрование: %v\n", simpleDeciphered == processedText)
	fmt.Printf("   Сложное дешифрование: %v\n", complexDeciphered == processedText)
}
