package providers

import "go.opentelemetry.io/otel"

var tracer = otel.Tracer("auth-provider-tracer")
