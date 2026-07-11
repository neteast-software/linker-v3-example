package tts

import (
	"context"
	"fmt"

	postgresql "github.com/neteast-software/go-module/db/postgresql/linker"
	grpclinker "github.com/neteast-software/go-module/rpc/grpc/linker"
	linker "github.com/neteast-software/linker/v3"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"

	ttsmodel "linker-v3-example/internal/model/tts"
	_ "linker-v3-example/internal/route/tts" // route 声明随组件进入编译
	ttsrpc "linker-v3-example/internal/rpc/tts"
	ttsservice "linker-v3-example/internal/service/tts"
)

const ID linker.ID = "example/tts"

type Component struct {
	service *ttsservice.Service
}

func NewComponent() *Component {
	return &Component{}
}

func (p *Component) Identity() linker.ID {
	return ID
}

func (p *Component) Dependencies() []linker.Dependency {
	return []linker.Dependency{linker.RequireComponent(postgresql.ID)}
}

func (p *Component) Assets(context.Context, linker.Runtime) ([]linker.Asset, error) {
	return []linker.Asset{
		postgresql.Table(&ttsmodel.Conversion{}, postgresql.Comment("演示 TTS 转写资产")),
		grpclinker.Register(func(server *grpc.Server) {
			ttsrpc.Register(server, p)
		}),
	}, nil
}

func (p *Component) Init(_ context.Context, runtime linker.Runtime) error {
	db, err := postgresql.Require(runtime)
	if err != nil {
		return err
	}
	p.service = ttsservice.New(db)
	return nil
}

func (p *Component) Stop(context.Context) error {
	if p.service == nil {
		return nil
	}
	return p.service.Close()
}

func (p *Component) Transcribe(ctx context.Context, req *wrapperspb.StringValue) (*wrapperspb.StringValue, error) {
	if p.service == nil {
		return nil, fmt.Errorf("tts 组件未初始化")
	}
	return p.service.Transcribe(ctx, req)
}
