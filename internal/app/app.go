package app

import (
	"context"
	"delivery/internal/config"
	"delivery/internal/pkg/closer"
	"log"
	"sync"
)

const (
	serviceName = "delivery"
)

type App struct {
	serviceProvider *serviceProvider
	configPath      string
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
	}

	wg := sync.WaitGroup{}
	wg.Add(len(runActions))

	for _, runAction := range runActions {
		currentRunAction := runAction
		go func() {
			defer wg.Done()

			err := currentRunAction.action()
			if err != nil {
				log.Fatalf(currentRunAction.errMsg)
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

func (a *App) runGRPCServer() error {
	// TODO: когда будем добавлять grpc сервер - реализовать
	return nil
}
