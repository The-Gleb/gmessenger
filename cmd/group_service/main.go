package main

import (
	"context"
	"github.com/The-Gleb/gmessenger/group_service/pkg/client/postgresql"
	"github.com/The-Gleb/gmessenger/internal/group_service/adapter/db/postgres"
	"github.com/The-Gleb/gmessenger/internal/group_service/config"
	serverapi "github.com/The-Gleb/gmessenger/internal/group_service/controller/grpc"
	"github.com/The-Gleb/gmessenger/internal/group_service/domain/service"
	"github.com/The-Gleb/gmessenger/internal/logger"
	"github.com/The-Gleb/gmessenger/pkg/proto/group"
	"google.golang.org/grpc"
	"log"
	"log/slog"
	"net"
)

func main() {

	cfg := config.MustBuild("config.yml")

	logger.Initialize("debug")
	slog.Info("here is config", "config", cfg)

	client, err := postgresql.NewClient(context.Background(), 2, cfg.DB)
	if err != nil {
		log.Fatal(err)
	}

	groupStorage := postgres.NewGroupStorage(client)
	messageStorage := postgres.NewMessageStorage(client)

	groupService := service.NewGroupService(groupStorage)
	messageService := service.NewMessageService(messageStorage)

	groupServerAPI := serverapi.NewServerAPI(messageService, groupService)

	l, err := net.Listen("tcp", cfg.ListenAddress)
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()

	group.RegisterGroupServer(grpcServer, groupServerAPI)

	slog.Info("starting group server", "config", cfg)
	if err := grpcServer.Serve(l); err != nil {
		log.Fatal(err)
	}
}
