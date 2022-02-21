package main

import (
	"context"
	"github.com/klwxsrx/arch-course-labs/auth-service/pkg/common/app/log"
	loggerImpl "github.com/klwxsrx/arch-course-labs/auth-service/pkg/common/infrastructure/logger"
	"github.com/klwxsrx/arch-course-labs/auth-service/pkg/userprofile/app"
	"github.com/klwxsrx/arch-course-labs/auth-service/pkg/userprofile/infrastructure/storage"
	"github.com/klwxsrx/arch-course-labs/auth-service/pkg/userprofile/infrastructure/transport"
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

	userProfileRepo := storage.NewUserProfileRepository()
	userProfileService := app.NewUserProfileService(userProfileRepo)

	server, err := startServer(userProfileService, logger)
	if err != nil {
		logger.WithError(err).Fatal("failed to start server")
	}
	logger.Info("app is ready")

	listenOSKillSignals()
	_ = server.Shutdown(context.Background())
}

func startServer(userProfileService *app.UserProfileService, logger log.Logger) (*http.Server, error) {
	handler, err := transport.NewHTTPHandler(userProfileService, logger)
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
