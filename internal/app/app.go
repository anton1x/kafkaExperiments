package app

import (
	"context"
	"github.com/IBM/sarama"
	"kafkaExperiments/internal/infrastructure/eventBus"
	"kafkaExperiments/internal/order/delivery/http"
	"kafkaExperiments/internal/order/repo"
	"kafkaExperiments/internal/order/usecase"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
	logger *slog.Logger
}

func New() *App {
	logger := slog.Default()
	return &App{logger: logger}
}

func (app *App) Run(ctx context.Context) error {
	ordersRepo := &repo.OrdersMockRepo{}

	// Создание продюсера Kafka
	producer, err := sarama.NewSyncProducer([]string{"localhost:29092", "localhost:39092"}, nil)
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	defer producer.Close()

	consumer, err := sarama.NewConsumer([]string{"localhost:29092", "localhost:39092"}, nil)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	defer consumer.Close()

	bus := eventBus.New(app.logger, producer, consumer, "events")

	ctl := usecase.NewOrderUseCase(app.logger, ordersRepo, bus)

	srv := http.NewServer(ctl)

	go func() {
		if err := srv.Run(ctx); err != nil {
			app.logger.InfoContext(ctx, "server error", err)
		}
	}()

	go func() {
		app.logger.Log(ctx, slog.LevelInfo, "listener started")
		busCh, err := bus.Listen(ctx)
		if err != nil {
			app.logger.Error("error listen", err)
			return
		}

		for event := range busCh {
			app.logger.Log(ctx, slog.LevelInfo, "event fired", "event", event)
		}

		app.logger.Info("end of listener")
	}()

	app.logger.InfoContext(ctx, "app started")

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	app.logger.InfoContext(ctx, "app stopped")

	return nil
}
