package main

import (
 "fmt"
 "sync"
 "time"
)

// Продукт
type Product struct {
 Name          string
 Weight        int     // Вес в граммах
 Calories      float64 // Калории на 100 грамм
 Proteins      float64 // Белки на 100 грамм
 Fats          float64 // Жиры на 100 грамм
 Carbohydrates float64 // Углеводы на 100 грамм
}

// Функция фильтрации продуктов без многопоточности
func filterProductsSequential(products []Product, maxCalories float64, maxCarbs float64) []Product {
 var result []Product
 for _, product := range products {
  totalCalories := (product.Calories * float64(product.Weight)) / 100
  if totalCalories <= maxCalories && product.Carbohydrates < maxCarbs {
   result = append(result, product)
  }
 }
 return result
}

// Функция фильтрации продуктов с многопоточностью
func filterProductsParallel(products []Product, maxCalories float64, maxCarbs float64) []Product {
 var result []Product
 var wg sync.WaitGroup
 var mu sync.Mutex

 numThreads := 4 // Количество потоков
 chunkSize := (len(products) + numThreads - 1) / numThreads

 for i := 0; i < len(products); i += chunkSize {
  end := i + chunkSize
  if end > len(products) {
   end = len(products)
  }

  wg.Add(1)
  go func(chunk []Product) {
   defer wg.Done()
   var localResult []Product
   for _, product := range chunk {
    totalCalories := (product.Calories * float64(product.Weight)) / 100
    if totalCalories <= maxCalories && product.Carbohydrates < maxCarbs {
     localResult = append(localResult, product)
    }
   }

   mu.Lock()
   result = append(result, localResult...)
   mu.Unlock()
  }(products[i:end])
 }

 wg.Wait()
 return result
}

func main() {
 // Данные продуктов
 products := []Product{
  {"Яблоко", 150, 52, 0.3, 0.2, 14},
  {"Банан", 120, 89, 1.1, 0.3, 22},
  {"Гречка", 500, 330, 12.6, 3.3, 62},
  {"Молоко", 1000, 42, 3.4, 1.0, 5},
  {"Курица", 600, 239, 27, 14, 0},
 }

 // Условия фильтрации
 maxCalories := 200.0
 maxCarbs := 15.0

 // Фильтрация без многопоточности
 start := time.Now()
 sequentialResult := filterProductsSequential(products, maxCalories, maxCarbs)
 sequentialDuration := time.Since(start)

 // Фильтрация с многопоточностью
 start = time.Now()
 parallelResult := filterProductsParallel(products, maxCalories, maxCarbs)
 parallelDuration := time.Since(start)

 // Вывод результатов
 fmt.Println("=== Фильтрация продуктов ===")
 fmt.Println("Без многопоточности:")
 for _, product := range sequentialResult {
  fmt.Printf("- %s (Вес: %d г, Калорий: %.2f на 100 г, Углеводов: %.2f)\n", product.Name, product.Weight, product.Calories, product.Carbohydrates)
 }
 fmt.Printf("Время выполнения: %s\n\n", sequentialDuration)

 fmt.Println("С многопоточностью:")
 for _, product := range parallelResult {
  fmt.Printf("- %s (Вес: %d г, Калорий: %.2f на 100 г, Углеводов: %.2f)\n", product.Name, product.Weight, product.Calories, product.Carbohydrates)
 }
 fmt.Printf("Время выполнения: %s\n", parallelDuration)
}
