package producer

import (
	"encoding/json"
	"context"
	"fmt"
	"log"
	"time"

	kafka "github.com/segmentio/kafka-go"
	gen "github.com/hryak228pizza/wbTech.L0/generator"
)

const (
	topic              = "new-orders-topic"
	kafkaBrokerAddress = "localhost:9093"
)

func Producer() {
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{kafkaBrokerAddress},
		Topic:   topic,
	})
	defer w.Close()

	order, err := json.Marshal(gen.NewOrder())
	if err != nil { 
		log.Fatal(err)
		return
	}

	sendMsg(w, string(order))

	fmt.Println("Producer finished")
}

func sendMsg(w *kafka.Writer, m string) {
	msg := kafka.Message{
		Value: []byte(m),
	}

	err := w.WriteMessages(context.Background(), msg)
	if err != nil {
		log.Printf("ошибка отправки сообщения в Кафку '%s': %v\n", msg.Value, err)
	} else {
		fmt.Printf("Отправленное в Кафку сообщение: %s\n", msg.Value)
	}
	time.Sleep(200 * time.Millisecond)
}