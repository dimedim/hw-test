package main

// import (
// 	"fmt"
// )

// func main() {
// 	fmt.Println("sender consumer")
// }

func main() {
	HelloExample()

	//TODO:
	/*
			Задачи рассыльщика:

		Прочитать тот же rabbitmq-блок конфига (только url и queue).

		Подключиться к RabbitMQ и убедиться, что очередь существует.

		Вызвать Consume(queue) и в цикле читать Delivery.

		Для каждого сообщения распарсить в Notification и просто вывести в STDOUT или лог:

		При SIGINT/TERM корректно закрыть соединение.
	*/

}
