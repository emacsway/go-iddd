package main

import (
	"context"
	"database/sql"
	"go-iddd/api/grpc/customer"
	"go-iddd/customer/application"
	"go-iddd/customer/domain"
	"go-iddd/customer/ports/secondary/customers"
	"go-iddd/shared"
	"go-iddd/shared/infrastructure/persistance/eventstore"
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	stopSignalChannel chan os.Signal
	logger            *logrus.Logger
	cancelCtx         context.CancelFunc
	grpcClientConn    *grpc.ClientConn
	grpcServer        *grpc.Server
	restServer        *http.Server
	postgresDBConn    *sql.DB
)

func main() {
	buildLogger()
	buildStopSignalChan()
	mustOpenPostgresDBConnection()

	go startGRPC()
	go startHTTP()

	waitForStopSignal()
}

func buildLogger() {
	if logger == nil {
		logger = logrus.New()
		formatter := &logrus.TextFormatter{
			FullTimestamp: true,
		}
		logger.SetFormatter(formatter)
	}
}

func buildStopSignalChan() {
	if stopSignalChannel == nil {
		stopSignalChannel = make(chan os.Signal, 1)
		signal.Notify(stopSignalChannel, os.Interrupt)
	}
}

func mustOpenPostgresDBConnection() {
	var err error

	if postgresDBConn == nil {
		logger.Info("opening Postgres DB connection ...")

		dsn := "postgresql://goiddd:password123@localhost:5432/goiddd_local?sslmode=disable"

		if postgresDBConn, err = sql.Open("postgres", dsn); err != nil {
			logger.Fatalf("failed to create Postgres DB connection: %s", err)
		}

		if err := postgresDBConn.Ping(); err != nil {
			logger.Fatalf("failed to connect to Postgres DB: %s", err)
		}
	}
}

func startGRPC() {
	logger.Info("starting gRPC server ...")

	listener, err := net.Listen("tcp", "localhost:5566")
	if err != nil {
		logger.Errorf("failed to listen: %v", err)
		stopSignalChannel <- os.Interrupt
	}

	grpcServer = grpc.NewServer()
	customerServer := customer.NewCustomerServer(buildCommandHandler())

	customer.RegisterCustomerServer(grpcServer, customerServer)
	reflection.Register(grpcServer)

	logger.Info("gRPC server ready ...")

	if err := grpcServer.Serve(listener); err != nil {
		logger.Errorf("gRPC server failed to serve: %s", err)
		stopSignalChannel <- os.Interrupt
	}
}

func startHTTP() {
	var err error
	var ctx context.Context

	logger.Info("starting REST server ...")

	ctx, cancelCtx = context.WithCancel(context.Background())

	grpcClientConn, err = grpc.Dial("localhost:5566", grpc.WithInsecure())
	if err != nil {
		logger.Errorf("fail to dial: %s", err)
		stopSignalChannel <- os.Interrupt
	}

	rmux := runtime.NewServeMux()
	client := customer.NewCustomerClient(grpcClientConn)

	if err = customer.RegisterCustomerHandlerClient(ctx, rmux, client); err != nil {
		logger.Errorf("failed to register customerHandlerClient: %s", err)
		stopSignalChannel <- os.Interrupt
	}

	// Serve the swagger-ui and swagger file
	mux := http.NewServeMux()
	mux.Handle("/", rmux)

	mux.HandleFunc(
		"/v1/customer/swagger.json",
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "api/grpc/customer/customer.swagger.json")
		},
	)

	restServer = &http.Server{
		Addr:    "localhost:8080",
		Handler: mux,
	}

	logger.Info("REST server ready - serving Swagger file at: http://localhost:8080/v1/customer/swagger.json")

	if err = restServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Errorf("REST server failed to listenAndServe: %s", err)
		stopSignalChannel <- os.Interrupt
	}
}

func buildCommandHandler() shared.CommandHandler {
	es := eventstore.NewPostgresEventStore(postgresDBConn, "eventstore", domain.UnmarshalDomainEvent)
	identityMap := customers.NewIdentityMap()
	repo := customers.NewEventSourcedRepository(es, domain.ReconstituteCustomerFrom, identityMap)
	commandHandler := application.NewCommandHandler(repo, postgresDBConn)

	return commandHandler
}

func waitForStopSignal() {
	s := <-stopSignalChannel

	logger.Infof("received '%s' - stopping services ...", s)

	if cancelCtx != nil {
		logger.Info("canceling context ...")
		cancelCtx()
	}

	if restServer != nil {
		logger.Info("stopping rest server gracefully ...")
		if err := restServer.Shutdown(context.Background()); err != nil {
			logger.Warnf("failed to stop the rest server: %s", err)
		}
	}

	if grpcClientConn != nil {
		logger.Info("closing grpc client connection ...")

		if err := grpcClientConn.Close(); err != nil {
			logger.Warnf("failed to close the grpc client connection: %s", err)
		}
	}

	if grpcServer != nil {
		logger.Info("stopping grpc server gracefully ...")
		grpcServer.GracefulStop()
	}

	if postgresDBConn != nil {
		logger.Info("closing Postgres DB connection ...")
		if err := postgresDBConn.Close(); err != nil {
			logger.Warnf("failed to close the Postgres DB connection: %s", err)
		}
	}

	close(stopSignalChannel)

	logger.Info("all services stopped - exiting")
}
