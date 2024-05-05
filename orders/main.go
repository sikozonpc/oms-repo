package main

import (
	"context"
	"fmt"
	"net"
	"time"

	_ "github.com/joho/godotenv/autoload"
	common "github.com/sikozonpc/commons"
	"github.com/sikozonpc/commons/broker"
	"github.com/sikozonpc/commons/discovery"
	"github.com/sikozonpc/commons/discovery/consul"
	"github.com/sikozonpc/omsv2-orders/gateway"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	serviceName = "orders"
	grpcAddr    = common.EnvString("GRPC_ADDR", "localhost:2000")
	consulAddr  = common.EnvString("CONSUL_ADDR", "localhost:8500")
	amqpUser    = common.EnvString("RABBITMQ_USER", "guest")
	amqpPass    = common.EnvString("RABBITMQ_PASS", "guest")
	amqpHost    = common.EnvString("RABBITMQ_HOST", "localhost")
	amqpPort    = common.EnvString("RABBITMQ_PORT", "5672")
	mongoUser   = common.EnvString("MONGO_DB_USER", "root")
	mongoPass   = common.EnvString("MONGO_DB_PASS", "example")
	mongoAddr   = common.EnvString("MONGO_DB_HOST", "localhost:27017")
	jaegerAddr  = common.EnvString("JAEGER_ADDR", "localhost:4318")
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	zap.ReplaceGlobals(logger)

	if err := common.SetGlobalTracer(context.TODO(), serviceName, jaegerAddr); err != nil {
		logger.Fatal("could set global tracer", zap.Error(err))
	}

	registry, err := consul.NewRegistry(consulAddr, serviceName)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, grpcAddr); err != nil {
		panic(err)
	}

	go func() {
		for {
			if err := registry.HealthCheck(instanceID, serviceName); err != nil {
				logger.Error("Failed to health check", zap.Error(err))
			}
			time.Sleep(time.Second * 1)
		}
	}()

	defer registry.Deregister(ctx, instanceID, serviceName)

	ch, close := broker.Connect(amqpUser, amqpPass, amqpHost, amqpPort)
	defer func() {
		close()
		ch.Close()
	}()

	// mongo db conn
	uri := fmt.Sprintf("mongodb://%s:%s@%s", mongoUser, mongoPass, mongoAddr)
	mongoClient, err := connectToMongoDB(uri)
	if err != nil {
		logger.Fatal("failed to connect to mongo db", zap.Error(err))
	}

	grpcServer := grpc.NewServer()

	l, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		logger.Fatal("failed to listen", zap.Error(err))
	}
	defer l.Close()

	gateway := gateway.NewGateway(registry)

	store := NewStore(mongoClient)
	svc := NewService(store, gateway)
	svcWithTelemetry := NewTelemetryMiddleware(svc)
	svcWithLogging := NewLoggingMiddleware(svcWithTelemetry)

	NewGRPCHandler(grpcServer, svcWithLogging, ch)

	consumer := NewConsumer(svcWithLogging)
	go consumer.Listen(ch)

	logger.Info("Starting HTTP server", zap.String("port", grpcAddr))

	if err := grpcServer.Serve(l); err != nil {
		logger.Fatal("failed to serve", zap.Error(err))
	}
}

func connectToMongoDB(uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	return client, err
}
