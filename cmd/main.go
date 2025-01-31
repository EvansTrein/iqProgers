package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/EvansTrein/iqProgers/internal/app"
	"github.com/EvansTrein/iqProgers/internal/config"
	"github.com/EvansTrein/iqProgers/pkg/logs"
)

func main() {
	var conf *config.Config
	var log *slog.Logger

	conf = config.MustLoad()
	log = logs.InitLog(conf.Env)

	application := app.New(conf, log)

	go func() {
		application.MustStart()
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-done
	if err := application.Stop(); err != nil {
		log.Error("an error occurred when stopping the application", "error", err)
	}
}

// Тестовое задание:
// Сделать REST API для финансовых операций с тремя ручками
// Пополнение баланса пользователя
// Перевод денег от одного пользователя к другому
// Просмотр 10 последних операций пользователя
// У пользователя есть баланс, а также список транзакций. Не забывать использовать SQL транзакции!!!
// Авторизацию делать не нужно!
// Технологии: go, gin, pgx, postgreSQL, docker. Для миграций использовать goose, выполнять конфигурацию через ENV. К решению также приложить файл .env.example
// Сделать запуск приложения через docker-compose при помощи Makefile. При вызове команды make run должны подняться контейнеры, выполниться миграции и запуститься сервер на порту 8080
// Для сложных запросов можно использовать query-builder. Также нужно написать unit-тесты для сервисного слоя приложения
// в приложении использовать подход clean architecture
// handler -> service -> repository
// Можно писать с уточняющими вопросами по задаче
// После выполнения прислать ссылку на github с выполненным тестовым заданием
// В случае успешного выполнения запланируем собеседование
