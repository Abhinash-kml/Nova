package social

import "go.opentelemetry.io/otel"

var tracer = otel.Tracer("social-module")
