package economy

import "go.opentelemetry.io/otel"

var tracer = otel.Tracer("currency-tracer")
