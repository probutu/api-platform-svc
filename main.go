package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"probutu/api-platform-svc/model"
	"probutu/api-platform-svc/repository"
	"probutu/api-platform-svc/service"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	_ "github.com/go-playground/validator/v10"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	app              = gin.New()
	env              = viper.GetString("env")
	serviceName      = mustEnv("SERVICE_NAME", "apiplatformsvc")
	serviceVersion   = mustEnv("SERVICE_VERSION", "v1.0.0")
	serviceNamespace = mustEnv("SERVICE_NAMESPACE", "example")
	tracer           = otel.Tracer(serviceName)
	otelendpoint     = viper.GetString("otel.endpoint")
	resources        *resource.Resource
	db               *gorm.DB
)

func mustEnv(env, d string) string {
	v, ok := os.LookupEnv(env)
	if ok {
		return v
	}

	return d
}

func initResource() {
	ctx := context.Background()

	extraResources, err := resource.New(ctx,
		resource.WithOS(),
		resource.WithProcess(),
		resource.WithContainer(),
		resource.WithHost(),
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion(serviceVersion),
			semconv.ServiceNamespace(serviceNamespace),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	res, err := resource.Merge(
		resource.Default(),
		extraResources,
	)
	if err != nil {
		log.Fatal(err)
	}

	resources = res
}

// Initializes an OTLP exporter, and configures the corresponding trace and
// metric providers.
func initTraceProvider() (func(context.Context) error, error) {
	ctx := context.Background()

	// If the OpenTelemetry Collector is running on a local cluster (minikube or
	// microk8s), it should be accessible through the NodePort service at the
	// `localhost:30080` endpoint. Otherwise, replace `localhost` with the
	// endpoint of your cluster. If you run the app inside k8s, then you can
	// probably connect directly to the service through dns.
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Set up a trace exporter
	traceClient := otlptracehttp.NewClient(
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithEndpoint(otelendpoint),
	)

	tracerExp, err := otlptrace.New(ctx, traceClient)
	if err != nil {
		zap.Error(err)
		return nil, err
	}

	// Register the trace exporter with a TracerProvider, using a batch
	// span processor to aggregate spans before export.
	bsp := sdktrace.NewBatchSpanProcessor(tracerExp)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithResource(resources),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)

	// set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{}, propagation.Baggage{},
	))

	// Shutdown will flush any remaining spans and shut down the exporter.
	return tracerProvider.Shutdown, nil
}

func initMetricProvider() (func(context.Context) error, error) {
	ctx := context.Background()

	// Set up a metrics exporter
	metricClient, err := otlpmetrichttp.New(ctx,
		otlpmetrichttp.WithInsecure(),
		otlpmetrichttp.WithEndpoint(otelendpoint),
	)
	if err != nil {
		return nil, err
	}

	mp := metric.NewMeterProvider(
		metric.WithReader(
			metric.NewPeriodicReader(metricClient),
		),
	)
	defer func() {
		if err := mp.Shutdown(ctx); err != nil {
			panic(err)
		}
	}()
	otel.SetMeterProvider(mp)

	return mp.Shutdown, nil
}

func initPostgre() (err error) {
	dsn := postgres.Open(
		fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
			viper.GetString("database.host"),
			viper.GetString("database.user"),
			viper.GetString("database.password"),
			viper.GetString("database.db"),
			viper.GetString("database.port"),
			viper.GetString("database.sslMode"),
			viper.GetString("database.timezone"),
		))

	db, err = gorm.Open(dsn)
	if err != nil {
		zap.Error(err)
		return
	}

	if viper.GetString("env") == "development" {
		// db = db.Debug()
	}

	db.Use(otelgorm.NewPlugin())

	return
}

func initViper() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("failed ReadInConfig: %v\n", err)
	}

	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
	viper.WatchConfig()
}

func init() {
	initViper()
	initPostgre()
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	shutdown, err := initTraceProvider()
	if err != nil {
		log.Fatal(err.Error())
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatal(err.Error())
		}
	}()

	// shutdown, err = initMetricProvider()
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }
	// defer func() {
	// 	if err := shutdown(ctx); err != nil {
	// 		log.Fatal(err.Error())
	// 	}
	// }()

	if err := db.AutoMigrate(
		&model.User{},
		&model.Workspace{},
		&model.WorkspaceUser{},
		&model.WorkspaceEnvironment{},
		&model.Environment{},
		&model.Collection{},
		&model.Folder{},
		&model.Request{},
		&model.Header{},
		&model.Response{},
	); err != nil {
		log.Fatal(err.Error())
	}

	app.Use(func(c *gin.Context) {
		ctx := c.Request.Context()
		header := c.Request.Header

		carrier := propagation.HeaderCarrier{}
		carrier.Set("Traceparent", header.Get("traceparent"))

		propgator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
		propgator.Inject(ctx, carrier)
		ctx, span := tracer.Start(propgator.Extract(ctx, carrier), c.Request.RequestURI)
		defer span.End()

		span.SetAttributes(
			attribute.String("user.id", header.Get("x-user-id")),
			attribute.String("user.roles", header.Get("x-user-roles")),
		)

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	})

	run()
}

func run() {

	v1 := app.Group("/v1")
	{
		users := v1.Group("/users")
		{
			users.POST("")
			users.GET("")
			user := users.Group("/:userId")
			{
				user.GET("")
				user.PUT("")
				user.DELETE("")
			}
		}

		workspaceRepository := repository.NewWorkspaceRepository(db)
		workspaceService := service.NewWorkspaceService(
			workspaceRepository,
		)
		workspaces := v1.Group("/workspaces")
		{
			workspaces.GET("", workspaceService.HandleList)
			workspaces.POST("", workspaceService.HandleCreate)
			workspace := workspaces.Group("/:workspaceId")
			{
				workspace.PUT("", workspaceService.HandleUpdate)
				workspace.DELETE("")
			}
		}

		collectionRepository := repository.NewCollectionRepository(db)
		collectionService := service.NewCollectionService(
			collectionRepository,
		)
		collections := v1.Group("/collections")
		{
			collections.GET("", collectionService.HandleList)
			collections.POST("", collectionService.HandleCreate)
			collection := collections.Group("/:collectionId")
			{
				collection.GET("")
				collection.PUT("")
				collection.DELETE("")
			}
		}

		folderRepository := repository.NewFolderRepository(db)
		folderService := service.NewFolderService(
			folderRepository,
		)
		folders := v1.Group("/folders")
		{
			folders.GET("", folderService.HandleList)
			folders.POST("", folderService.HandleCreate)
			folder := folders.Group("/:folderId")
			{
				folder.PUT("")
				folder.DELETE("")
			}
		}

		requestRepository := repository.NewRequestRepository(db)
		requestService := service.NewRequestService(
			requestRepository,
		)

		requests := v1.Group("/requests")
		{
			requests.GET("", requestService.HandleList)
			requests.POST("", requestService.HandleCreate)
			request := requests.Group("/:requestId")
			{
				request.GET("")
				request.PUT("")
				request.DELETE("")
			}
		}
	}

	m := http.NewServeMux()
	m.Handle("/metrics", promhttp.Handler())
	m.Handle("/", app)

	h := http.Server{Addr: mustEnv("PORT", ":8080"), Handler: m}
	if err := h.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

	select {}
}
