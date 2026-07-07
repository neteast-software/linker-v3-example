package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	yaml "github.com/neteast-software/go-module/config/yaml"
)

const Namespace = "linker-v3-example"
const EnvConfig = "LINKER_V3_EXAMPLE_CONFIG"

type Source interface {
	LoadConfig(context.Context, Config) (Config, error)
}

type SourceFunc func(context.Context, Config) (Config, error)

func (p SourceFunc) LoadConfig(ctx context.Context, seed Config) (Config, error) {
	return p(ctx, seed)
}

type LoadOption func(*loadOptions)

type loadOptions struct {
	files   []string
	sources []Source
	env     bool
}

func Load(ctx context.Context, options ...LoadOption) (Config, error) {
	opts := loadOptions{
		files: configFilesFromEnv(),
		env:   true,
	}
	for _, option := range options {
		if option != nil {
			option(&opts)
		}
	}

	cfg := Default()
	if len(opts.files) > 0 {
		next, err := loadYAMLFiles(cfg, opts.files...)
		if err != nil {
			return Config{}, fmt.Errorf("读取本地配置失败: %w", err)
		}
		cfg = next
	}
	for _, source := range opts.sources {
		if source == nil {
			continue
		}
		next, err := source.LoadConfig(ctx, cfg)
		if err != nil {
			return Config{}, fmt.Errorf("读取配置中心失败: %w", err)
		}
		cfg = next
	}
	if opts.env {
		ApplyEnv(&cfg)
	}
	return cfg, nil
}

func WithFiles(paths ...string) LoadOption {
	return func(opts *loadOptions) {
		opts.files = append(opts.files, cleanPaths(paths)...)
	}
}

func WithSource(sources ...Source) LoadOption {
	return func(opts *loadOptions) {
		opts.sources = append(opts.sources, sources...)
	}
}

func WithoutEnv() LoadOption {
	return func(opts *loadOptions) {
		opts.env = false
	}
}

func loadYAMLFiles(seed Config, paths ...string) (Config, error) {
	setting, err := yaml.LoadFiles(paths...)
	if err != nil {
		return Config{}, err
	}
	data, ok := setting.Lookup(Namespace)
	if !ok {
		return seed, nil
	}
	if err = json.Unmarshal(data, &seed); err != nil {
		return Config{}, fmt.Errorf("%s: %w", Namespace, err)
	}
	return seed, nil
}

func configFilesFromEnv() []string {
	return cleanPaths(strings.Split(os.Getenv(EnvConfig), ","))
}

func cleanPaths(paths []string) []string {
	ret := make([]string, 0, len(paths))
	for _, path := range paths {
		path = strings.TrimSpace(path)
		if path != "" {
			ret = append(ret, path)
		}
	}
	return ret
}
