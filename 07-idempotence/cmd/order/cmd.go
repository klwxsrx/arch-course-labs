package main

import (
	"context"
	"github.com/klwxsrx/arch-course-labs/idempotence/data/mysql/order"
	"github.com/klwxsrx/arch-course-labs/idempotence/pkg/common/app/log"
	loggerImpl "github.com/klwxsrx/arch-course-labs/idempotence/pkg/common/infra/logger"
	commonMysql "github.com/klwxsrx/arch-course-labs/idempotence/pkg/common/infra/mysql"
	"github.com/klwxsrx/arch-course-labs/idempotence/pkg/order/app/service"
	"github.com/klwxsrx/arch-course-labs/idempotence/pkg/order/infra/mysql"
	"github.com/klwxsrx/arch-course-labs/idempotence/pkg/order/infra/transport"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger := loggerImpl.New()

	config, err := parseConfig()
	if err != nil {
		logger.WithError(err).Fatal("failed to parse config")
	}

	db, client, err := getDatabaseClient(config, logger)
	if err != nil {
		logger.WithError(err).Fatal("failed to setup db connection")
	}
	defer db.Close()

	migration, err := commonMysql.NewMigration(client, logger, order.MysqlMigrations)
	if err != nil {
		logger.WithError(err).Fatal("failed to setup db migration")
	}
	err = migration.Migrate()
	if err != nil {
		logger.WithError(err).Fatal("failed to execute db migration")
	}

	unitOfWork := mysql.NewUnitOfWork(client)
	orderService := service.NewOrderService(unitOfWork, logger)

	server, err := startServer(orderService, logger)
	if err != nil {
		logger.WithError(err).Fatal("failed to start server")
	}
	logger.Info("app is ready")

	listenOSKillSignals()
	_ = server.Shutdown(context.Background())
}

func getDatabaseClient(config *config, logger log.Logger) (commonMysql.Connection, commonMysql.TransactionalClient, error) {
	db, err := commonMysql.NewConnection(commonMysql.Config{DSN: commonMysql.Dsn{
		User:     config.DBUser,
		Password: config.DBPassword,
		Host:     config.DBHost,
		Port:     config.DBPort,
		Database: config.DBName,
	}}, logger)
	if err != nil {
		return nil, nil, err
	}
	client, err := db.Client()
	if err != nil {
		db.Close()
		return nil, nil, err
	}
	return db, client, nil
}

func startServer(orderService *service.OrderService, logger log.Logger) (*http.Server, error) {
	handler, err := transport.NewHTTPHandler(orderService, logger)
	if err != nil {
		return nil, err
	}
	srv := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.WithError(err).Fatal("unable to start the server")
		}
	}()
	return srv, nil
}

func listenOSKillSignals() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	<-ch
}
