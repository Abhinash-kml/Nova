package observability

import (
	"context"
	"errors"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.uber.org/zap"
)

var (
	tProvider *trace.TracerProvider = nil
	mProvider *metric.MeterProvider = nil
	lProvider *log.LoggerProvider   = nil
)

// setupOTelSDK bootstraps the OpenTelemetry pipeline.
// If it does not return an error, make sure to call shutdown for proper cleanup.
func SetupOTelSDK(ctx context.Context) (func(context.Context) error, error) {
	var shutdownFuncs []func(context.Context) error
	var err error

	// shutdown calls cleanup functions registered via shutdownFuncs.
	// The errors from the calls are joined.
	// Each registered cleanup will be invoked once.
	shutdown := func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}

	// handleErr calls shutdown for cleanup and makes sure that all errors are returned.
	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}

	// Set up propagator.
	prop := newPropagator()
	otel.SetTextMapPropagator(prop)

	// Creaste a new resource
	resource, err := newResource(ctx, "nova-server", "1.0.0", "meow")
	if err != nil {
		handleErr(err)
		return shutdown, err
	}

	// Set up trace provider.
	tracerProvider, err := newTracerProvider(resource)
	if err != nil {
		handleErr(err)
		return shutdown, err
	}
	tProvider = tracerProvider

	shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)
	otel.SetTracerProvider(tracerProvider)

	//Set up meter provider.
	meterProvider, err := newMeterProvider(resource)
	if err != nil {
		handleErr(err)
		return shutdown, err
	}
	mProvider = meterProvider

	shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
	otel.SetMeterProvider(meterProvider)

	// Set up logger provider.
	loggerProvider, err := newLoggerProvider(resource)
	if err != nil {
		handleErr(err)
		return shutdown, err
	}
	lProvider = loggerProvider

	shutdownFuncs = append(shutdownFuncs, loggerProvider.Shutdown)
	global.SetLoggerProvider(loggerProvider)

	return shutdown, err
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func newResource(ctx context.Context, name, version, environment string) (*resource.Resource, error) {
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(name),
			semconv.ServiceVersionKey.String(version),
			semconv.DeploymentEnvironmentNameKey.String(environment),
		))
	if err != nil {
		zap.L().Error("Failed to create semconv resource", zap.Error(err))
		return nil, err
	}

	return res, nil
}

func newTracerProvider(res *resource.Resource) (*trace.TracerProvider, error) {
	// traceExporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	// if err != nil {
	// 	return nil, err
	// }

	traceExporter, err := otlptracegrpc.New(context.Background(),
		otlptracegrpc.WithEndpoint("localhost:4317"),
		otlptracegrpc.WithInsecure())
	if err != nil {
		zap.L().Fatal("Failed to created Otel trace exporter", zap.Error(err))
	}

	tracerProvider := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter,
			trace.WithBatchTimeout(time.Second)), // Default is 5s. Set to 1s for demonstrative purposes.
		trace.WithResource(res),
	)

	return tracerProvider, nil
}

func newMeterProvider(res *resource.Resource) (*metric.MeterProvider, error) {
	// metricExporter, err := stdoutmetric.New(stdoutmetric.WithPrettyPrint())
	// if err != nil {
	// 	return nil, err
	// }

	metricExporter, err := otlpmetricgrpc.New(context.Background(),
		otlpmetricgrpc.WithEndpoint("localhost:4317"),
		otlpmetricgrpc.WithInsecure())
	if err != nil {
		zap.L().Fatal("Failed to create Otel metric exporter")
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(metricExporter,
			metric.WithInterval(3*time.Second))),
		metric.WithResource(res), // Default is 1m. Set to 3s for demonstrative purposes.
	)

	return meterProvider, nil
}

func newLoggerProvider(res *resource.Resource) (*log.LoggerProvider, error) {
	// logExporter, err := stdoutlog.New(stdoutlog.WithPrettyPrint())
	// if err != nil {
	// 	return nil, err
	// }

	logExporter, err := otlploggrpc.New(context.Background(),
		otlploggrpc.WithEndpoint("localhost:4317"),
		otlploggrpc.WithInsecure())
	if err != nil {
		zap.L().Fatal("Failed to create Otel log exporter")
	}

	loggerProvider := log.NewLoggerProvider(
		log.WithProcessor(log.NewBatchProcessor(logExporter)),
		log.WithResource(res),
	)

	return loggerProvider, nil
}

func TraceProvider() *trace.TracerProvider {
	return tProvider
}

func MetricsProvider() *metric.MeterProvider {
	return mProvider
}

func LoggerProvider() *log.LoggerProvider {
	return lProvider
}
