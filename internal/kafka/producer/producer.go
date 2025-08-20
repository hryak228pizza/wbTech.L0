package producer

import (
	"encoding/json"
	"context"
	"fmt"
	"log"
	"time"
	"bytes"

	kafka "github.com/segmentio/kafka-go"
	gen "github.com/hryak228pizza/wbTech.L0/internal/generator"
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

	c := time.Tick(5 * time.Second)

	for range c {
		order, err := json.Marshal(gen.NewOrder())
		if err != nil { 
			log.Fatal(err)
			return
		}

		sendMsg(w, order)
	}	

	// fmt.Println("Producer finished")
}

func sendMsg(w *kafka.Writer, m []byte) {
	msg := kafka.Message{
		Value: m,
	}

	err := w.WriteMessages(context.Background(), msg)
	if err != nil {
		log.Printf("ошибка отправки сообщения в Кафку '%s': %v\n", msg.Value, err)
	} else {
		var prettyJSON bytes.Buffer
    	err := json.Indent(&prettyJSON, m, "", "  ")
		if err != nil {
			log.Print(err)
		} 
		//fmt.Printf("Отправленное в Кафку сообщение: %s\n", msg.Value)
		fmt.Printf("Отправленное в Кафку сообщение: ", prettyJSON.String())
	}
	time.Sleep(200 * time.Millisecond)
}