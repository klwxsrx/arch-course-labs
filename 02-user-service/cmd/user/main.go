package main

import (
	"context"
	"github.com/klwxsrx/arch-course-labs/user-service/pkg/common/app/log"
	loggerImpl "github.com/klwxsrx/arch-course-labs/user-service/pkg/common/infrastructure/logger"
	commonMysql "github.com/klwxsrx/arch-course-labs/user-service/pkg/common/infrastructure/mysql"
	"github.com/klwxsrx/arch-course-labs/user-service/pkg/user/app"
	"github.com/klwxsrx/arch-course-labs/user-service/pkg/user/infrastructure/mysql"
	"github.com/klwxsrx/arch-course-labs/user-service/pkg/user/infrastructure/transport"
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

	userService := app.NewUserService(mysql.NewUserRepository(client))

	server := startServer(userService, logger)
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

func startServer(userService app.UserService, logger log.Logger) *http.Server {
	srv := &http.Server{
		Addr:    ":8080",
		Handler: transport.NewHTTPHandler(userService, logger),
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.WithError(err).Fatal("unable to start the server")
		}
	}()
	return srv
}

func listenOSKillSignals() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	<-ch
}
