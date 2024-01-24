package eventBus

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/hashicorp/go-uuid"
	"log/slog"
	"time"
)

type EventBus struct {
	PubChannel sarama.SyncProducer
	Receiver   sarama.Consumer
	EventTopic string
	logger     *slog.Logger
}

type EventMessage struct {
	ID     string      `json:"ID"`
	Event  string      `json:"event"`
	Params interface{} `json:"params"`
}

func New(logger *slog.Logger, producer sarama.SyncProducer, receiver sarama.Consumer, topic string) *EventBus {
	return &EventBus{
		PubChannel: producer,
		EventTopic: topic,
		logger:     logger,
		Receiver:   receiver,
	}
}

func (eb *EventBus) toProducerMessage(msg EventMessage) (id string, res *sarama.ProducerMessage, err error) {
	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return "", nil, fmt.Errorf("error generating uuid: %w", err)
	}

	jsonBytes, err := json.Marshal(msg)
	if err != nil {
		return "", nil, fmt.Errorf("error marshaling: %w", err)
	}

	res = &sarama.ProducerMessage{
		Topic:     eb.EventTopic,
		Key:       sarama.StringEncoder(uuid),
		Value:     sarama.ByteEncoder(jsonBytes),
		Timestamp: time.Now().Add(30 * time.Second),
	}

	return uuid, res, nil
}

func (eb *EventBus) Publish(ctx context.Context, msg EventMessage) error {
	uuid, convertedMsg, err := eb.toProducerMessage(msg)
	if err != nil {
		fmt.Errorf("error publish msg: %w", err)
	}

	pt, offset, err := eb.PubChannel.SendMessage(convertedMsg)
	if err != nil {
		return fmt.Errorf("error publish msg: %w", err)
	}

	eb.logger.InfoContext(ctx, "event was succesfully published", "eventInfo", map[string]any{
		"pt":     pt,
		"offset": offset,
		"uuid":   uuid,
	})

	return nil
}

func (eb *EventBus) Listen(ctx context.Context) (<-chan EventMessage, error) {
	ch := make(chan EventMessage)

	pts, err := eb.Receiver.Partitions("events")
	if err != nil {
		return nil, err
	}

	for _, pt := range pts {
		pt := pt
		go func() {
			cons, err := eb.Receiver.ConsumePartition("events", pt, sarama.OffsetNewest)
			if err != nil {
				return
			}

			for {
				select {
				case msg := <-cons.Messages():
					decoded := EventMessage{}
					err := json.Unmarshal(msg.Value, &decoded)
					if err == nil {
						ch <- decoded
					}
					break
				case <-ctx.Done():
					close(ch)
					return
				}
			}
		}()
	}

	return ch, nil
}
