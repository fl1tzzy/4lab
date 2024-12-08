package main

import (
	"fmt"
	"sync"
	"time"
)

// Структура ReadersWriters управляет доступом к общему ресурсу для читателей и писателей.
type ReadersWriters struct {
	mu             sync.Mutex  // Мьютекс для обеспечения взаимного исключения.
	cv             *sync.Cond  // Условные переменные для ожидания и сигнализации.
	writers        int         // Счетчик активных писателей.
	readers        int         // Счетчик активных читателей.
	writersPriority bool       // Флаг, указывающий, имеют ли писатели приоритет.
}

// Функция NewReadersWriters создает и возвращает новый экземпляр ReadersWriters.
func NewReadersWriters() *ReadersWriters {
	rw := &ReadersWriters{
		writersPriority: true,  // Инициализация флага приоритета писателей.
	}
	rw.cv = sync.NewCond(&rw.mu)  // Инициализация условной переменной с использованием мьютекса.
	return rw
}

// Метод startWrite блокирует доступ к ресурсу для других писателей и читателей.
func (rw *ReadersWriters) startWrite() {
	rw.mu.Lock()  // Блокировка мьютекса для обеспечения взаимного исключения.
	defer rw.mu.Unlock()  // Разблокировка мьютекса после завершения метода.
	for rw.writers > 0 || (!rw.writersPriority && rw.readers > 0) {
		rw.cv.Wait()  // Ожидание, пока не будет возможности начать запись.
	}
	rw.writers++  // Увеличение счетчика активных писателей.
}

// Метод stopWrite завершает запись и сигнализирует другим горутинам.
func (rw *ReadersWriters) stopWrite() {
	rw.mu.Lock()  // Блокировка мьютекса для обеспечения взаимного исключения.
	defer rw.mu.Unlock()  // Разблокировка мьютекса после завершения метода.
	rw.writers--  // Уменьшение счетчика активных писателей.
	if rw.writers == 0 {
		rw.cv.Broadcast()  // Сигнализирование всем ожидающим горутинам, что ресурс свободен.
	}
}

// Метод startRead блокирует доступ к ресурсу для читателей, если писатель активен.
func (rw *ReadersWriters) startRead() {
	rw.mu.Lock()  // Блокировка мьютекса для обеспечения взаимного исключения.
	defer rw.mu.Unlock()  // Разблокировка мьютекса после завершения метода.
	for rw.writers > 0 {
		rw.cv.Wait()  // Ожидание, пока не будет возможности начать чтение.
	}
	rw.readers++  // Увеличение счетчика активных читателей.
}

// Метод stopRead завершает чтение и сигнализирует другим горутинам.
func (rw *ReadersWriters) stopRead() {
	rw.mu.Lock()  // Блокировка мьютекса для обеспечения взаимного исключения.
	defer rw.mu.Unlock()  // Разблокировка мьютекса после завершения метода.
	rw.readers--  // Уменьшение счетчика активных читателей.
	rw.cv.Broadcast()  // Сигнализирование всем ожидающим горутинам, что ресурс свободен.
}

// Метод write симулирует процесс записи данных.
func (rw *ReadersWriters) write(id, iterations int) {
	for i := 0; i < iterations; i++ {
		rw.startWrite()  // Начало записи.
		fmt.Printf("Writer %d writing\n", id)  // Вывод сообщения о записи.
		time.Sleep(time.Second)  // Задержка на 1 секунду.
		rw.stopWrite()  // Завершение записи.
	}
}

// Метод read симулирует процесс чтения данных.
func (rw *ReadersWriters) read(id, iterations int) {
	for i := 0; i < iterations; i++ {
		rw.startRead()  // Начало чтения.
		fmt.Printf("Reader %d reading\n", id)  // Вывод сообщения о чтении.
		rw.stopRead()  // Завершение чтения.
	}
}

// Функция execute создает горутины для писателей и читателей.
func execute(rw *ReadersWriters) {
	iterations := 1  // Количество итераций для каждого писателя и читателя.
	writeLimit := 2  // Количество писателей.
	readLimit := 4  // Количество читателей.

	if rw.writersPriority {
		fmt.Println("Writer's priority")  // Вывод сообщения о приоритете писателей.
	} else {
		fmt.Println("\nReader priority")  // Вывод сообщения о приоритете читателей.
	}

	var wg sync.WaitGroup  // Создание группы ожидания для синхронизации горутин.

	for i := 0; i < writeLimit; i++ {
		wg.Add(1)  // Добавление горутины в группу ожидания.
		go func(id int) {
			defer wg.Done()  // Уменьшение счетчика группы ожидания после завершения горутины.
			rw.write(id+1, iterations)  // Запуск горутины для записи.
		}(i)
	}

	for i := 0; i < readLimit; i++ {
		wg.Add(1)  // Добавление горутины в группу ожидания.
		go func(id int) {
			defer wg.Done()  // Уменьшение счетчика группы ожидания после завершения горутины.
			rw.read(id+1, iterations)  // Запуск горутины для чтения.
		}(i)
	}

	wg.Wait()  // Ожидание завершения всех горутин в группе.
}

// Основная функция main.
func main() {
	rw := NewReadersWriters()  // Создание экземпляра ReadersWriters с приоритетом писателей.

	execute(rw)  // Выполнение горутин с приоритетом писателей.

	rw.writersPriority = false  // Изменение приоритета на приоритет читателей.
	execute(rw)  // Выполнение горутин с приоритетом читателей.
}