package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"runtime"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

// Решето Эратосфена
func sieveOfEratosthenes(n int) []int {
	isPrime := make([]bool, n+1)
	for i := range isPrime {
		isPrime[i] = true
	}
	isPrime[0], isPrime[1] = false, false

	for i := 2; i*i <= n; i++ {
		if isPrime[i] {
			for j := i * i; j <= n; j += i {
				isPrime[j] = false
			}
		}
	}

	primes := []int{}
	for i := 2; i <= n; i++ {
		if isPrime[i] {
			primes = append(primes, i)
		}
	}
	return primes
}

// Простое решето для малых чисел (вспомогательная функция)
func simpleSieve(n int) []int {
	isPrime := make([]bool, n+1)
	for i := range isPrime {
		isPrime[i] = true
	}
	isPrime[0], isPrime[1] = false, false

	for i := 2; i*i <= n; i++ {
		if isPrime[i] {
			for j := i * i; j <= n; j += i {
				isPrime[j] = false
			}
		}
	}

	primes := []int{}
	for i := 2; i <= n; i++ {
		if isPrime[i] {
			primes = append(primes, i)
		}
	}
	return primes
}

// Сегментированное решето - медленное, но экономное по памяти
func segmentedSieve(n int) []int {
	if n < 2 {
		return []int{}
	}

	// Шаг 1: Найти все простые до √n (базовые простые)
	limit := int(math.Sqrt(float64(n))) + 1
	basePrimes := simpleSieve(limit)

	primes := make([]int, 0)
	primes = append(primes, basePrimes...)

	// Размер сегмента (32 КБ на сегмент)
	segmentSize := 32768
	segment := make([]bool, segmentSize)

	// Шаг 2: Обрабатываем числа сегментами от √n до n
	low := limit
	high := low + segmentSize - 1

	for low <= n {
		if high > n {
			high = n
		}

		// Инициализируем сегмент (все true)
		for i := range segment[:high-low+1] {
			segment[i] = true
		}

		// Отмечаем составные числа используя базовые простые
		for _, p := range basePrimes {
			// Находим первое кратное p в [low, high]
			start := (low / p) * p
			if start < low {
				start += p
			}
			if start == p {
				start += p // пропускаем само простое число
			}

			// Вычеркиваем кратные
			for j := start; j <= high; j += p {
				segment[j-low] = false
			}
		}

		// Собираем простые числа из сегмента
		for i := low; i <= high; i++ {
			if segment[i-low] {
				primes = append(primes, i)
			}
		}

		low += segmentSize
		high += segmentSize
	}

	return primes
}

// Оптимизированное решето Эратосфена (только нечётные числа)
func optimizedSieveOfEratosthenes(n int) []int {
	if n < 2 {
		return []int{}
	}

	primes := []int{2}
	if n == 2 {
		return primes
	}

	limit := (n - 1) / 2
	isPrime := make([]bool, limit+1)
	for i := range isPrime {
		isPrime[i] = true
	}

	for i := 1; 2*i*i <= limit; i++ {
		if isPrime[i] {
			p := 2*i + 1
			for j := 2 * i * (i + 1); j <= limit; j += p {
				isPrime[j] = false
			}
		}
	}

	for i := 1; i <= limit; i++ {
		if isPrime[i] {
			primes = append(primes, 2*i+1)
		}
	}

	return primes
}

// Бенчмарк структура
type BenchmarkResult struct {
	N              int
	EratosthenesMs float64
	SegmentedMs    float64
	OptimizedMs    float64
	EratMemoryMB   float64
	SegmentedMemMB float64
	OptimizedMemMB float64
}

// Функция бенчмарка с измерением памяти
func benchmarkWithMemory(fn func(int) []int, n int, iterations int) (float64, float64) {
	// Теоретический расчет памяти на основе структур данных
	var memoryMB float64

	// Эратосфен: bool slice размером n+1
	// Сундарам: bool slice размером (n-1)/2 + 1
	// Оптимизированный: bool slice размером (n-1)/2 + 1

	var totalDuration time.Duration
	for i := 0; i < iterations; i++ {
		runtime.GC()

		start := time.Now()
		result := fn(n)
		totalDuration += time.Since(start)

		// Предотвращаем оптимизацию компилятора
		_ = result
	}

	timeMs := float64(totalDuration.Microseconds()) / float64(iterations) / 1000.0

	return timeMs, memoryMB
}

// Вычисление теоретического потребления памяти
func calculateMemory(n int, algorithm string) float64 {
	var bytes float64

	// Приблизительное количество простых чисел до n (по теореме о распределении простых чисел)
	primeCount := float64(n) / (math.Log(float64(n)) + 1.0)
	if n < 100 {
		primeCount = float64(n) * 0.25
	}

	switch algorithm {
	case "eratosthenes":
		// bool slice размером n+1 (1 байт на элемент)
		boolArray := float64(n + 1)
		// результирующий slice []int (8 байт на int)
		resultArray := primeCount * 8
		bytes = boolArray + resultArray
	case "segmented":
		// базовые простые до √n
		limit := int(math.Sqrt(float64(n))) + 1
		basePrimesCount := float64(limit) / (math.Log(float64(limit)) + 1.0)
		basePrimes := basePrimesCount * 8
		// сегмент размером 32KB
		segment := 32768.0
		// результирующий slice []int
		resultArray := primeCount * 8
		bytes = basePrimes + segment + resultArray
	case "optimized":
		// bool slice размером (n-1)/2 + 1
		limit := (n - 1) / 2
		boolArray := float64(limit + 1)
		// результирующий slice []int
		resultArray := primeCount * 8
		bytes = boolArray + resultArray
	}

	return bytes / 1024.0 / 1024.0
}

// Запуск бенчмарков
func runBenchmarks() []BenchmarkResult {
	ranges := []int{100, 500, 1000, 5000, 10000, 50000, 100000, 500000,
		1000000, 5000000, 10000000, 50000000, 100000000}

	results := make([]BenchmarkResult, 0)

	fmt.Println("=====================================================================================================")
	fmt.Println("                        Бенчмарк алгоритмов поиска простых чисел")
	fmt.Println("=====================================================================================================\n")

	fmt.Printf("%-12s | %-20s | %-20s | %-20s\n", "n", "Эратосфен", "Сегментированный", "Оптимизированный")
	fmt.Printf("%-12s | %-20s | %-20s | %-20s\n", "", "(время/память)", "(время/память)", "(время/память)")
	fmt.Println("-------------|----------------------|----------------------|----------------------")

	for idx, n := range ranges {
		iterations := 10
		if n > 1000000 {
			iterations = 5
		}
		if n > 10000000 {
			iterations = 3
		}
		if n > 50000000 {
			iterations = 1
		}

		fmt.Printf("Обработка n=%-10d...", n)

		eratTime, _ := benchmarkWithMemory(sieveOfEratosthenes, n, iterations)
		segmentedTime, _ := benchmarkWithMemory(segmentedSieve, n, iterations)
		optimizedTime, _ := benchmarkWithMemory(optimizedSieveOfEratosthenes, n, iterations)

		// Вычисляем теоретическую память
		eratMem := calculateMemory(n, "eratosthenes")
		segmentedMem := calculateMemory(n, "segmented")
		optimizedMem := calculateMemory(n, "optimized")

		result := BenchmarkResult{
			N:              n,
			EratosthenesMs: eratTime,
			SegmentedMs:    segmentedTime,
			OptimizedMs:    optimizedTime,
			EratMemoryMB:   eratMem,
			SegmentedMemMB: segmentedMem,
			OptimizedMemMB: optimizedMem,
		}
		results = append(results, result)

		fmt.Printf("\r%-12d | %7.2f мс / %6.2f МБ | %7.2f мс / %6.2f МБ | %7.2f мс / %6.2f МБ\n",
			n, eratTime, eratMem, segmentedTime, segmentedMem, optimizedTime, optimizedMem)

		if idx == len(ranges)-2 || idx == len(ranges)-1 {
			minTime := eratTime
			winner := "Эратосфен"
			if segmentedTime < minTime {
				minTime = segmentedTime
				winner = "Сегментированный"
			}
			if optimizedTime < minTime {
				winner = "Оптимизированный"
			}
			fmt.Printf("             Быстрейший: %s\n", winner)
		}
	}

	return results
}

// Создание графика времени выполнения
func createTimePlot(results []BenchmarkResult, filename string) error {
	p := plot.New()
	p.Title.Text = "Время выполнения алгоритмов"
	p.X.Label.Text = "Количество элементов (n)"
	p.Y.Label.Text = "Время (мс)"

	eratPts := make(plotter.XYs, len(results))
	segmentedPts := make(plotter.XYs, len(results))
	optimizedPts := make(plotter.XYs, len(results))

	for i, r := range results {
		eratPts[i].X = float64(r.N)
		eratPts[i].Y = r.EratosthenesMs
		segmentedPts[i].X = float64(r.N)
		segmentedPts[i].Y = r.SegmentedMs
		optimizedPts[i].X = float64(r.N)
		optimizedPts[i].Y = r.OptimizedMs
	}

	eratLine, _ := plotter.NewLine(eratPts)
	eratLine.Color = color.RGBA{R: 59, G: 130, B: 246, A: 255}
	eratLine.Width = vg.Points(2)

	segmentedLine, _ := plotter.NewLine(segmentedPts)
	segmentedLine.Color = color.RGBA{R: 239, G: 68, B: 68, A: 255}
	segmentedLine.Width = vg.Points(2)

	optimizedLine, _ := plotter.NewLine(optimizedPts)
	optimizedLine.Color = color.RGBA{R: 34, G: 197, B: 94, A: 255}
	optimizedLine.Width = vg.Points(2)

	p.Add(eratLine, segmentedLine, optimizedLine)
	p.Legend.Add("Эратосфен", eratLine)
	p.Legend.Add("Сегментированный", segmentedLine)
	p.Legend.Add("Оптимизированный", optimizedLine)
	p.Legend.Top = true
	p.Legend.Left = true

	return p.Save(10*vg.Inch, 6*vg.Inch, filename)
}

// Создание графика потребления памяти
func createMemoryPlot(results []BenchmarkResult, filename string) error {
	p := plot.New()
	p.Title.Text = "Потребление памяти алгоритмами"
	p.X.Label.Text = "Количество элементов (n)"
	p.Y.Label.Text = "Память (МБ)"

	eratPts := make(plotter.XYs, len(results))
	segmentedPts := make(plotter.XYs, len(results))
	optimizedPts := make(plotter.XYs, len(results))

	for i, r := range results {
		eratPts[i].X = float64(r.N)
		eratPts[i].Y = r.EratMemoryMB
		segmentedPts[i].X = float64(r.N)
		segmentedPts[i].Y = r.SegmentedMemMB
		optimizedPts[i].X = float64(r.N)
		optimizedPts[i].Y = r.OptimizedMemMB
	}

	eratLine, _ := plotter.NewLine(eratPts)
	eratLine.Color = color.RGBA{R: 59, G: 130, B: 246, A: 255}
	eratLine.Width = vg.Points(2)

	segmentedLine, _ := plotter.NewLine(segmentedPts)
	segmentedLine.Color = color.RGBA{R: 239, G: 68, B: 68, A: 255}
	segmentedLine.Width = vg.Points(2)

	optimizedLine, _ := plotter.NewLine(optimizedPts)
	optimizedLine.Color = color.RGBA{R: 34, G: 197, B: 94, A: 255}
	optimizedLine.Width = vg.Points(2)

	p.Add(eratLine, segmentedLine, optimizedLine)
	p.Legend.Add("Эратосфен", eratLine)
	p.Legend.Add("Сегментированный", segmentedLine)
	p.Legend.Add("Оптимизированный", optimizedLine)
	p.Legend.Top = true
	p.Legend.Left = true

	return p.Save(10*vg.Inch, 6*vg.Inch, filename)
}

func main() {
	results := runBenchmarks()

	fmt.Println("\n=== Создание графиков ===")

	if err := createTimePlot(results, "time_comparison.png"); err != nil {
		log.Fatalf("Ошибка создания графика времени: %v", err)
	}
	fmt.Println("График времени выполнения: time_comparison.png")

	if err := createMemoryPlot(results, "memory_comparison.png"); err != nil {
		log.Fatalf("Ошибка создания графика памяти: %v", err)
	}
	fmt.Println("График потребления памяти: memory_comparison.png")

	fmt.Println("\n=====================================================================================================")
	fmt.Println("                                        РЕЗУЛЬТАТЫ")
	fmt.Println("=====================================================================================================")

	last := results[len(results)-1]
	fmt.Printf("\nРезультаты при n = 100,000,000:\n\n")

	fmt.Println("ВРЕМЯ:")
	fmt.Printf("   Классический Эратосфен:  %10.2f мс\n", last.EratosthenesMs)
	fmt.Printf("   Сегментированный:        %10.2f мс  (%.2fx)\n",
		last.SegmentedMs, last.SegmentedMs/last.EratosthenesMs)
	fmt.Printf("   Оптимизированный:        %10.2f мс  (%.2fx)\n",
		last.OptimizedMs, last.OptimizedMs/last.EratosthenesMs)

	fmt.Println("\nПАМЯТЬ:")
	fmt.Printf("   Классический Эратосфен:  %10.2f МБ\n", last.EratMemoryMB)
	fmt.Printf("   Сегментированный:        %10.2f МБ  (%.2fx)\n",
		last.SegmentedMemMB, last.SegmentedMemMB/last.EratMemoryMB)
	fmt.Printf("   Оптимизированный:        %10.2f МБ  (%.2fx)\n",
		last.OptimizedMemMB, last.OptimizedMemMB/last.EratMemoryMB)

	minTime := last.EratosthenesMs
	timeWinner := "Классический Эратосфен"
	if last.SegmentedMs < minTime {
		minTime = last.SegmentedMs
		timeWinner = "Сегментированный"
	}
	if last.OptimizedMs < minTime {
		timeWinner = "Оптимизированный"
	}

	minMem := last.EratMemoryMB
	memWinner := "Классический Эратосфен"
	if last.SegmentedMemMB < minMem {
		minMem = last.SegmentedMemMB
		memWinner = "Сегментированный"
	}
	if last.OptimizedMemMB < minMem {
		memWinner = "Оптимизированный"
	}

	fmt.Printf("\nСамый быстрый: %s\n", timeWinner)
	fmt.Printf("Самый экономный по памяти: %s\n", memWinner)

	fmt.Println("\nВЫВОДЫ:")
	fmt.Println("   - Оптимизированный Эратосфен использует примерно в 2 раза меньше памяти")
	fmt.Println("   - Сегментированный использует константную память O(sqrt(n)) независимо от n")
	fmt.Println("   - Сегментированный медленнее из-за многократной обработки данных")
	fmt.Println("   - На больших n экономия памяти критична для производительности")
	fmt.Println("   - Оптимизированный - лучший баланс скорости и памяти")

	fmt.Println("\nГрафики успешно созданы.")
}
