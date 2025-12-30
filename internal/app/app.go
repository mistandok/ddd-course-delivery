package app

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	httpmiddleware "delivery/internal/adapters/in/http/middleware"
	"delivery/internal/config"
	"delivery/internal/core/domain/model/event"
	"delivery/internal/generated/servers"
	"delivery/internal/pkg/closer"
	"delivery/internal/pkg/outbox/mappers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/robfig/cron/v3"

	"github.com/mehdihadeli/go-mediatr"
)

const (
	serviceName = "delivery"
)

type App struct {
	serviceProvider *serviceProvider
	configPath      string
	httpServer      *http.Server
	cronScheduler   *cron.Cron
}

func NewApp(ctx context.Context, configPath string) (*App, error) {
	a := &App{configPath: configPath}

	if err := a.initDeps(ctx); err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	runActions := []struct {
		action func() error
		errMsg string
	}{
		{action: a.runGRPCServer, errMsg: "ошибка при запуске GRPC сервера"},
		{action: a.runHttpServer, errMsg: "ошибка при запуске HTTP сервера"},
		{action: a.runCronScheduler, errMsg: "ошибка при запуске Cron планировщика"},
		{action: a.runKafkaConsumerGroup, errMsg: "ошибка при запуске Kafka consumer group"},
	}

	wg := sync.WaitGroup{}
	wg.Add(len(runActions))

	for _, runAction := range runActions {
		currentRunAction := runAction
		go func() {
			defer wg.Done()

			err := currentRunAction.action()
			if err != nil {
				log.Fatalf("%s", currentRunAction.errMsg)
			}
		}()
	}

	wg.Wait()

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	initDepFunctions := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.registerEvents,
		a.initMediator,
		a.initHttpServer,
		a.initCronScheduler,
	}

	for _, f := range initDepFunctions {
		if err := f(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	err := config.Load(a.configPath)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	// TODO: когда будем добавлять grpc сервер - реализовать
	return nil
}

func (a *App) initHttpServer(ctx context.Context) error {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(httpmiddleware.ErrorHandlingMiddleware())

	servers.RegisterHandlers(e, a.serviceProvider.HttpHandlers())

	httpConfig := a.serviceProvider.HttpConfig()
	a.httpServer = &http.Server{
		Addr:         httpConfig.Address(),
		Handler:      e,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	closer.Add(func() error {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return a.httpServer.Shutdown(shutdownCtx)
	})

	return nil
}

func (a *App) registerEvents(ctx context.Context) error {
	if err := a.serviceProvider.EventRegistry().RegisterDomainEvent(&event.OrderCreated{},
		mappers.NewOrderCreatedToJSONMapper(),
		mappers.NewOrderCreatedFromJSONMapper(),
	); err != nil {
		return err
	}

	if err := a.serviceProvider.EventRegistry().RegisterDomainEvent(&event.OrderCompleted{},
		mappers.NewOrderCompletedToJSONMapper(),
		mappers.NewOrderCompletedFromJSONMapper(),
	); err != nil {
		return err
	}

	return nil
}

func (a *App) initMediator(_ context.Context) error {
	if err := mediatr.RegisterNotificationHandler[*event.OrderCreated](a.serviceProvider.OrderCreatedHandler()); err != nil {
		return err
	}

	if err := mediatr.RegisterNotificationHandler[*event.OrderCompleted](a.serviceProvider.OrderCompletedHandler()); err != nil {
		return err
	}

	return nil
}

func (a *App) runGRPCServer() error {
	// TODO: когда будем добавлять grpc сервер - реализовать
	return nil
}

func (a *App) initCronScheduler(ctx context.Context) error {
	a.cronScheduler = cron.New()

	_, err := a.cronScheduler.AddJob("@every 1s", a.serviceProvider.AssignOrdersJob())
	if err != nil {
		return err
	}

	_, err = a.cronScheduler.AddJob("@every 1s", a.serviceProvider.MoveCouriersJob())
	if err != nil {
		return err
	}

	closer.Add(func() error {
		ctx := a.cronScheduler.Stop()
		<-ctx.Done()
		return nil
	})

	return nil
}

func (a *App) runCronScheduler() error {
	log.Printf("Starting Cron scheduler")
	a.cronScheduler.Start()
	select {} // Block forever
}

func (a *App) runKafkaConsumerGroup() error {
	log.Printf("Starting Kafka consumer group")
	return a.serviceProvider.BasketConfirmedConsumerGroup().Consume()
}

func (a *App) runHttpServer() error {
	log.Printf("Starting HTTP server on %s", a.httpServer.Addr)
	return a.httpServer.ListenAndServe()
}
