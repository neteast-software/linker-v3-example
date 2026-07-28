package tts

import (
	"context"

	postgresql "github.com/neteast-software/go-module/db/postgresql/linker"
	grpclinker "github.com/neteast-software/go-module/rpc/grpc/linker"
	linker "github.com/neteast-software/linker/v3"
	"google.golang.org/grpc"

	tts "linker-v3-example/internal/tts"
	ttsrpc "linker-v3-example/internal/tts/rpc"
)

const ID linker.ID = "example/tts"

type Component struct {
	service *tts.Service
	server  *ttsrpc.Server
}

func New() *Component {
	return &Component{server: ttsrpc.New()}
}

func (p *Component) Identity() linker.ID {
	return ID
}

func (p *Component) Capabilities() linker.Capabilities {
	return linker.Capabilities{linker.Need(postgresql.Key())}
}

func (p *Component) Assets(context.Context, linker.Runtime) ([]linker.Asset, error) {
	return []linker.Asset{
		postgresql.Table(&tts.Conversion{}, postgresql.Comment("演示 TTS 转写资产")),
		grpclinker.Register(func(server *grpc.Server) {
			tts.Register(server, p.server)
		}),
	}, nil
}

func (p *Component) Init(_ context.Context, runtime linker.Runtime) error {
	db, err := postgresql.Require(runtime)
	if err != nil {
		return err
	}
	p.service = tts.New(tts.NewStore(db))
	p.server.Configure(p.service)
	return nil
}

func (p *Component) Stop(context.Context) error {
	if p.service == nil {
		return nil
	}
	return p.service.Close()
}
