package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mark3labs/mcp-go/server"
	"soundofgothic.pl/backend/internal/config"
	"soundofgothic.pl/backend/internal/mcp"
	"soundofgothic.pl/backend/internal/mcp/resources"
	"soundofgothic.pl/backend/internal/mcp/tools"
	"soundofgothic.pl/backend/internal/postgres"
	"soundofgothic.pl/backend/internal/postgres/migrations"
)

func run() int {
	var configPath string
	flag.StringVar(&configPath, "config", "", "Path to config file")
	flag.StringVar(&configPath, "c", "", "Path to config file (shorthand)")
	flag.Parse()

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	cfg, err := config.LoadConfigFromFile(configPath)
	if err != nil {
		slog.Error("failed to load config", "error", err)
		return 1
	}

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))

	repositories, err := postgres.NewPostgresRepositories(postgres.WithAuth(
		postgres.DBAuth{
			Host:     cfg.Postgres.Host,
			Port:     cfg.Postgres.Port,
			Username: cfg.Postgres.User,
			Password: cfg.Postgres.Password,
			Name:     cfg.Postgres.Database,
		},
	))
	if err != nil {
		slog.Error("failed to connect to Postgres", "error", err)
		return 1
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := migrations.Migrate(ctx, repositories); err != nil {
		slog.Error("failed to run migrations", "error", err)
		return 1
	}

	mcpServer := mcp.NewMCPServer(
		mcp.WithResources(
			resources.NewActorResource(repositories),
		),
		mcp.WithTools(
			tools.NewGothicToolsPack(repositories)...,
		),
	)

	switch cfg.MCP.Address {
	case "stdio":
		err = server.ServeStdio(mcpServer.MCPServer, server.WithStdioContextFunc(func(ctx context.Context) context.Context {
			ctx, _ = signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
			return ctx
		}))
		if err != nil && !errors.Is(err, context.Canceled) {
			slog.Error("failed to start MCP server on stdio", "error", err)
			return 1
		}
	default:
		log.Printf("Starting MCP server on %s\n", cfg.MCP.Address)
		sseServer := server.NewSSEServer(mcpServer.MCPServer,
			server.WithStaticBasePath("/mcp"),
		)
		if err := sseServer.Start(cfg.MCP.Address); err != nil {
			slog.Error("failed to start MCP server", "address", cfg.MCP.Address, "error", err)
			return 1
		}
	}

	return 0
}

func main() {
	os.Exit(run())
}
