package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"cloud.google.com/go/cloudsqlconn"
	"github.com/Sedayu/client-vendor/handler"
	"github.com/Sedayu/client-vendor/repository"
	"github.com/Sedayu/client-vendor/service"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/labstack/echo/v4"
	"github.com/subosito/gotenv"

	_ "github.com/lib/pq"
)

func main() {
	// Loading environment variables.. Local only
	_ = gotenv.Load()

	var db *sql.DB
	var err error
	// Initiate database connections
	if os.Getenv("APP_ENV") == "local" {
		db, err = sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", os.Getenv("POSTGRES_USERNAME"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT"), os.Getenv("POSTGRES_DB")))
		if err != nil {
			panic(err)
		}
	} else {
		db, err = connectWithConnector()
		if err != nil {
			panic(err)
		}
	}

	defer db.Close()

	// Initiate repo package
	vendorRepository := repository.NewVendors(db)

	// Initiate service package
	vendorService := service.NewVendorsFinderProvider(vendorRepository)

	// Initiate http handler
	vendorHandler := handler.NewVendors(vendorService)

	// Initiate HTTP server and register handler
	e := echo.New()
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK!")
	})
	e.GET("/v1/vendors", vendorHandler.GetVendors)

	// Run HTTP server
	go func(e *echo.Echo) {
		port := fmt.Sprintf(":%s", os.Getenv("PORT"))
		fmt.Println("Lagoe API starting at port", port)
		e.Logger.Fatal(e.Start(port))
	}(e)

	// Listen to close signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-c

	fmt.Println("Closing API connections...")
}

func connectWithConnector() (*sql.DB, error) {
	mustGetenv := func(k string) string {
		v := os.Getenv(k)
		if v == "" {
			log.Fatalf("Fatal Error in connect_connector.go: %s environment variable not set.\n", k)
		}
		return v
	}

	var (
		dbUser                 = mustGetenv("POSTGRES_USERNAME")        // e.g. 'my-db-user'
		dbPwd                  = mustGetenv("POSTGRES_PASSWORD")        // e.g. 'my-db-password'
		dbName                 = mustGetenv("POSTGRES_DB")              // e.g. 'my-database'
		instanceConnectionName = mustGetenv("INSTANCE_CONNECTION_NAME") // e.g. 'project:region:instance'
	)

	dsn := fmt.Sprintf("user=%s password=%s database=%s", dbUser, dbPwd, dbName)
	config, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	var opts []cloudsqlconn.Option
	d, err := cloudsqlconn.NewDialer(context.Background(), opts...)
	if err != nil {
		return nil, err
	}

	// Use the Cloud SQL connector to handle connecting to the instance.
	// This approach does *NOT* require the Cloud SQL proxy.
	config.DialFunc = func(ctx context.Context, network, instance string) (net.Conn, error) {
		return d.Dial(ctx, instanceConnectionName)
	}
	dbURI := stdlib.RegisterConnConfig(config)
	dbPool, err := sql.Open("pgx", dbURI)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}
	return dbPool, nil
}
