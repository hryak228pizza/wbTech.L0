package producer

import (
	"encoding/json"
	"context"
	"time"
	"bytes"

	kafka "github.com/segmentio/kafka-go"
	"github.com/hryak228pizza/wbTech.L0/internal/generator"
	"github.com/hryak228pizza/wbTech.L0/internal/logger"
	"go.uber.org/zap"
)

const (
	topic              = "new-orders-topic"
	kafkaBrokerAddress = "localhost:9093"
)

// runs a producer for simulate new orders with kafka
func Producer() {

	// init writer
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{kafkaBrokerAddress},
		Topic:   topic,
	})
	defer w.Close()

	// ticker for generator
	c := time.Tick(5 * time.Second)

	// sending msg every 5sec
	for range c {
		order, err := json.Marshal(generator.NewOrder())
		if err != nil { 
			logger.L().Error("serialization failed", 
				zap.String("error", err.Error()),
			)
			return
		}
		sendMsg(w, order)
	}
}

// writes new message into kafka
func sendMsg(w *kafka.Writer, m []byte) {

	// message init
	msg := kafka.Message{
		Value: m,
	}

	err := w.WriteMessages(context.Background(), msg)
	if err != nil {
		logger.L().Error("kafka message write failed", 
			zap.String("message", string(msg.Value)),
			zap.String("error", err.Error()),
		)
	} else {
		// make readable json
		var prettyJSON bytes.Buffer
    	err := json.Indent(&prettyJSON, m, "", "  ")
		if err != nil {
			logger.L().Error("json indent failed",
				zap.String("error", err.Error()),
			)
		}
		logger.L().Info("kafka message send",
			zap.String("message", prettyJSON.String()),
		)
	}
}