package main

import (
	"context"
	"log"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
)

func main() {
	exporter, err := jaeger.New(jaeger.WithAgentEndpoint())
	if err != nil {
		log.Fatalf("failed to get exporter; %s", err)
	}
	traceProvider := trace.NewTracerProvider(
		trace.WithSyncer(exporter),
		trace.WithResource(resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceNameKey.String("turbo-funicular"),
			semconv.ServiceVersionKey.String("1.0"))),
	)
	otel.SetTracerProvider(traceProvider)
	defer func() {
		if err := traceProvider.Shutdown(context.Background()); err != nil {
			log.Printf("failed to shutdown trace provider; %s", err)
		}
	}()

	e := echo.New()
	e.Use(otelecho.Middleware("turbo-funicular"))
	e.GET("/users/:id", handleUsers)
	e.Start(":8000")
}

func handleUsers(ctx echo.Context) error {
	userID := ctx.Param("id")
	name := getUser(ctx.Request().Context(), userID)
	result := map[string]string{"id": userID, "name": name}
	return ctx.JSON(http.StatusOK, result)
}

func getUser(ctx context.Context, id string) string {
	tr := otel.Tracer("turbo-funicular")
	_, span := tr.Start(ctx, "getUser")
	defer span.End()
	switch id {
	case "1":
		return "alpha"
	case "2":
		return "bravo"
	default:
		return "unknown"
	}
}
