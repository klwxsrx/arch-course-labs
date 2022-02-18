package main

import (
	"context"
	"github.com/klwxsrx/arch-course-labs/auth-service/pkg/auth/app"
	"github.com/klwxsrx/arch-course-labs/auth-service/pkg/auth/infrastructure/auth"
	"github.com/klwxsrx/arch-course-labs/auth-service/pkg/auth/infrastructure/encoding"
	"github.com/klwxsrx/arch-course-labs/auth-service/pkg/auth/infrastructure/storage"
	"github.com/klwxsrx/arch-course-labs/auth-service/pkg/auth/infrastructure/transport"
	"github.com/klwxsrx/arch-course-labs/auth-service/pkg/common/app/log"
	loggerImpl "github.com/klwxsrx/arch-course-labs/auth-service/pkg/common/infrastructure/logger"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	logger := loggerImpl.New()

	userRepo := storage.NewUserRepository()
	passwordEncoder := encoding.NewPasswordEncoder()
	userService := app.NewUserService(userRepo, passwordEncoder)
	sessionService := auth.NewSessionService(storage.NewSessionStorage(), userRepo, passwordEncoder)

	server, err := startServer(userService, sessionService, logger)
	if err != nil {
		logger.WithError(err).Fatal("failed to start server")
	}
	logger.Info("app is ready")

	listenOSKillSignals()
	_ = server.Shutdown(context.Background())
}

func startServer(userService *app.UserService, sessionService *auth.SessionService, logger log.Logger) (*http.Server, error) {
	handler, err := transport.NewHTTPHandler(userService, sessionService, logger)
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
