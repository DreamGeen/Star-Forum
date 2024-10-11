package tracing

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.9.0"
	trace2 "go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"star/app/constant/settings"
	"star/app/utils/logging"
)

var Tracer trace2.Tracer

func SetTraceProvider(name string) (*trace.TracerProvider, error) {
	client := otlptracehttp.NewClient(
		otlptracehttp.WithEndpoint(settings.Conf.TracingEndPoint),
		otlptracehttp.WithInsecure())
	//导出器
	exporter, err := otlptrace.New(context.Background(), client)
	if err != nil {
		logging.Logger.Error("create exporter error",
			zap.Error(err))
		return nil, err
	}
	//采样器
	var sampler trace.Sampler
	if settings.Conf.OtelState == "disable" {
		sampler = trace.NeverSample()
	} else {
		sampler = trace.TraceIDRatioBased(settings.Conf.OtelSampler)
	}
	//追踪提供者
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(name))),
		trace.WithSampler(sampler))
	//将创建的追踪提供者设置为全局使用
	otel.SetTracerProvider(tp)
	//配置文本映射传播器，处理追踪上下文和附加信息
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	//创建追踪实例
	Tracer = otel.Tracer(name)
	return tp, nil
}
