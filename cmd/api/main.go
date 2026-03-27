package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	callbackapp "github.com/nathang0147/Request-Service-API/internal/callback"
	"github.com/nathang0147/Request-Service-API/internal/persistence/postgres"
	platformconfig "github.com/nathang0147/Request-Service-API/internal/platform/config"
	platformlogger "github.com/nathang0147/Request-Service-API/internal/platform/logger"
	platformrouter "github.com/nathang0147/Request-Service-API/internal/platform/router"
	providerresolver "github.com/nathang0147/Request-Service-API/internal/provider"
	waltprovider "github.com/nathang0147/Request-Service-API/internal/provider/walt"
	verificationapp "github.com/nathang0147/Request-Service-API/internal/verification"
	"go.uber.org/zap"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := run(ctx); err != nil {
		slog.New(slog.NewTextHandler(os.Stderr, nil)).Error("service exited", "error", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	cfg, err := platformconfig.Load()
	if err != nil {
		return err
	}

	logger, err := platformlogger.New(cfg.LogLevel)
	if err != nil {
		return err
	}
	defer func() {
		_ = logger.Sync()
	}()

	pool, err := postgres.Open(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()

	waltClient := waltprovider.NewClient(
		cfg.WaltVerifierBaseURL,
		cfg.WaltVerifierMode,
		cfg.WaltBearerToken,
		cfg.WaltVCPolicyWebhookURL,
		cfg.CallbackBaseURL,
		cfg.PublicBaseURL,
		cfg.PublicRedirectTemplate,
		cfg.CallbackAuthSecret,
		http.DefaultClient,
	)
	walt := waltprovider.New(waltClient)
	resolver := providerresolver.NewResolver(walt)
	verificationStore := postgres.NewVerificationStore(pool)
	verificationService := verificationapp.NewServiceWithResolver(cfg.DefaultProvider, verificationStore, resolver)
	verificationRoutes := verificationapp.NewHandler(verificationService)
	callbackStore := postgres.NewCallbackStore(pool)
	callbackService := callbackapp.NewService(callbackStore, walt)
	callbackRoutes := callbackapp.NewHandler(callbackService, callbackapp.NewStaticAuthenticator(cfg.CallbackAuthSecret))

	server := newServer(cfg.Port, newHandler(logger, verificationRoutes, callbackRoutes))

	serverErrors := make(chan error, 1)
	go func() {
		logger.Info("starting request service api", zap.String("addr", server.Addr))
		serverErrors <- server.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		return shutdownServer(server)
	case err := <-serverErrors:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
}

func newHandler(logger *zap.Logger, registrars ...platformrouter.RouteRegistrar) http.Handler {
	return platformrouter.New(logger, registrars...)
}

func newServer(port string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:              ":" + port,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}
}

func shutdownServer(server *http.Server) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}
