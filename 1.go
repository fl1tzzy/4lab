package main

import (
 "fmt"
 "math/rand"
 "sync"
 "time"
)

// MeasureExecutionTime - измерение времени выполнения
func MeasureExecutionTime(name string, fn func()) {
 start := time.Now()
 fn()
 elapsed := time.Since(start)
 fmt.Printf("%s execution time: %s\n", name, elapsed)
}

// === Генерация случайных символов в гонке ===
func generateASCIIInRace(numThreads int) {
 var wg sync.WaitGroup
 raceOutput := make([]string, numThreads)
 wg.Add(numThreads)

 for i := 0; i < numThreads; i++ {
  go func(id int) {
   defer wg.Done()
   rand.Seed(time.Now().UnixNano() + int64(id))
   randomChar := string(rune(rand.Intn(95) + 32)) // Генерация случайного символа из таблицы ASCII (32-126)
   fmt.Printf("Goroutine %d generated: %s\n", id, randomChar)
   raceOutput[id] = randomChar
  }(i)
 }

 wg.Wait()
 fmt.Println("Race Output:", raceOutput)
}

// === Пример использования Mutex ===
func mutexExample() {
 var mu sync.Mutex
 counter := 0

 var wg sync.WaitGroup
 wg.Add(10)

 for i := 0; i < 10; i++ {
  go func(id int) {
   defer wg.Done()
   mu.Lock()
   counter++
   time.Sleep(10 * time.Millisecond) // Имитация работы
   mu.Unlock()
  }(i)
 }

 wg.Wait()
 fmt.Println("Final Counter (Mutex):", counter)
}

// === Пример использования Семафора ===
func semaphoreExample() {
 var sem = make(chan struct{}, 3) // Семафор с 3 слотами
 var wg sync.WaitGroup

 for i := 0; i < 10; i++ {
  wg.Add(1)
  go func(id int) {
   defer wg.Done()
   sem <- struct{}{} // Захват семафора
   fmt.Printf("Goroutine %d started\n", id)
   time.Sleep(50 * time.Millisecond) // Имитация работы
   fmt.Printf("Goroutine %d finished\n", id)
   <-sem // Освобождение семафора
  }(i)
 }

 wg.Wait()
}

// === Пример использования Барьера ===
func barrierExample() {
 const numGoroutines = 5
 var barrier sync.WaitGroup
 barrier.Add(numGoroutines)

 for i := 0; i < numGoroutines; i++ {
  go func(id int) {
   defer barrier.Done()
   fmt.Printf("Goroutine %d reached the barrier\n", id)
   time.Sleep(20 * time.Millisecond) // Имитация работы
  }(i)
 }

 barrier.Wait()
 fmt.Println("All goroutines have reached the barrier")
}

// === Пример использования Монтитора ===
type Monitor struct {
 mu    sync.Mutex
 value int
}

func (m *Monitor) Increment() {
 m.mu.Lock()
 defer m.mu.Unlock()
 m.value++
 time.Sleep(10 * time.Millisecond) // Имитация работы
 fmt.Println("Incremented value to:", m.value)
}

func monitorExample() {
 monitor := Monitor{}
 var wg sync.WaitGroup
 wg.Add(10)

 for i := 0; i < 10; i++ {
  go func() {
   defer wg.Done()
   monitor.Increment()
  }()
 }

 wg.Wait()
 fmt.Println("Final Value (Monitor):", monitor.value)
}

// === Основная функция ===
func main() {
 fmt.Println("=== Анализ производительности ===")
 MeasureExecutionTime("Race with Random ASCII", func() { generateASCIIInRace(10) })
 MeasureExecutionTime("Mutex", mutexExample)
 MeasureExecutionTime("Semaphore", semaphoreExample)
 MeasureExecutionTime("Barrier", barrierExample)
 MeasureExecutionTime("Monitor", monitorExample)
}