package main

import (
	"fmt"  // Импортируем пакет для форматированного ввода-вывода
	"sync" // Импортируем пакет для синхронизации горутин
)

var wg sync.WaitGroup // Создаем переменную WaitGroup для ожидания завершения горутин

// Функция MyFunc принимает канал для получения строковых данных
func MyFunc(data chan string) {
	defer wg.Done() // Уменьшаем счетчик WaitGroup при завершении функции
	fmt.Println("Waiting for data...") // Сообщение о том, что функция ожидает данные
	text := <-data // Получаем данные из канала
	fmt.Println(text) // Печатаем полученные данные
}

func main() {
	dataChan := make(chan string) // Создаем канал для передачи строк
	fmt.Println("Run MyFunc goroutine") // Сообщение о запуске горутины
	wg.Add(1) // Увеличиваем счетчик WaitGroup на 1
	go MyFunc(dataChan) // Запускаем MyFunc как горутину
	dataChan <- "hehey" // Отправляем строку "hehey" в канал
	wg.Wait() // Ожидаем завершения всех горутин
}

