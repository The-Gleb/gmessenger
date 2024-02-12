package main

import (
	"context"
	"log/slog"
	"net"

	"github.com/The-Gleb/gmessenger/app/pkg/proto/go/group"
	"github.com/The-Gleb/gmessenger/group_service/internal/adapter/db/postgres"
	"github.com/The-Gleb/gmessenger/group_service/internal/config"
	serverapi "github.com/The-Gleb/gmessenger/group_service/internal/controller/grpc"
	"github.com/The-Gleb/gmessenger/group_service/internal/domain/service"
	"github.com/The-Gleb/gmessenger/group_service/internal/logger"
	"github.com/The-Gleb/gmessenger/group_service/pkg/client/postgresql"
	"google.golang.org/grpc"
)

func main() {

	cfg := config.MustBuild("config.yml")

	logger.Initialize("debug")
	slog.Info("here is config", "config", cfg)

	client, err := postgresql.NewClient(context.Background(), 2, cfg.DB)
	if err != nil {
		panic(err)
	}

	groupStorage := postgres.NewGroupStorage(client)
	messageStorage := postgres.NewMessageStorage(client)

	groupService := service.NewGroupService(groupStorage)
	messageService := service.NewMessageService(messageStorage)

	groupServerAPI := serverapi.NewServerAPI(messageService, groupService)

	l, err := net.Listen("tcp", "localhost:8081")
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()

	group.RegisterGroupServer(grpcServer, groupServerAPI)

	slog.Info("lister address", l.Addr())

	if err := grpcServer.Serve(l); err != nil {
		panic(err)
	}
}
