package tts

import (
	"github.com/neteast-software/go-module/observe/metrics"
	metricgrpc "github.com/neteast-software/go-module/observe/metrics/rpc/grpc"
	tracegrpc "github.com/neteast-software/go-module/observe/tracing/rpc/grpc"
	grpclinker "github.com/neteast-software/go-module/rpc/grpc/linker"
	linker "github.com/neteast-software/linker/v3"
	"google.golang.org/grpc"

	ttsrpc "linker-v3-example/internal/rpc/tts"
)

const ID linker.ID = "rpc/client/tts"

type Client = ttsrpc.Client

type Option func(*providerOptions)

type providerOptions struct {
	recorder metrics.Recorder
	labels   []metrics.LabelValue
}

func New(conn grpc.ClientConnInterface) Client {
	return ttsrpc.NewClient(conn)
}

func ClientKey() linker.CapabilityKey[Client] {
	return grpclinker.ClientKey[Client](ID)
}

func WithMetricRecorder(recorder metrics.Recorder) Option {
	return func(p *providerOptions) {
		p.recorder = recorder
	}
}

func WithMetricLabels(labels ...metrics.LabelValue) Option {
	return func(p *providerOptions) {
		p.labels = append(p.labels, labels...)
	}
}

func Provider(opts ...Option) linker.Component {
	options := providerOptions{}
	for _, opt := range opts {
		if opt != nil {
			opt(&options)
		}
	}
	interceptors := []grpc.UnaryClientInterceptor{tracegrpc.UnaryClient()}
	if options.recorder != nil {
		interceptors = append(interceptors, metricgrpc.UnaryClient(
			options.recorder,
			metricgrpc.WithConstLabels(options.labels...),
		))
	}
	return grpclinker.NewClientProvider[Client](
		ID,
		New,
		grpclinker.WithDialOptions[Client](grpc.WithChainUnaryInterceptor(interceptors...)),
	)
}
