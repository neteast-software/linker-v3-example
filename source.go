package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	env "github.com/neteast-software/go-module/config/env/linker"
	yaml "github.com/neteast-software/go-module/config/yaml/linker"
	nacos "github.com/neteast-software/go-module/registry/nacos/linker"
	linker "github.com/neteast-software/linker/v3"

	gateway "linker-v3-example/internal/gateway/linker"
)

const configPathEnv = "LINKER_V3_EXAMPLE_CONFIG"
const configOverridePrefix = "APP_"
const gatewayConfigPath = "config/gateway.example.yaml"
const gatewayNacosConfigPath = "config/gateway.nacos.example.yaml"
const gatewayRoutesPath = "config/gateway.routes.yaml"

func configSources(extra ...linker.Source) ([]linker.Source, error) {
	sources := []linker.Source{yaml.File(configPaths()...)}
	source, err := nacosSource()
	if err != nil {
		return nil, err
	}
	if source != nil {
		sources = append(sources, source)
	}
	sources = append(sources, extra...)
	sources = append(sources, env.Prefix(configOverridePrefix))
	return sources, nil
}

func gatewaySources(nacosProfile bool) ([]linker.Source, error) {
	path := gatewayConfigPath
	if nacosProfile {
		path = gatewayNacosConfigPath
	}
	sources := []linker.Source{yaml.File(path)}
	if nacosProfile {
		source, err := nacosSource()
		if err != nil {
			return nil, err
		}
		if source != nil {
			sources = append(sources, source)
		}
	} else {
		sources = append(sources, gateway.Routes(gatewayRoutesPath))
	}
	sources = append(sources, env.Prefix(configOverridePrefix))
	return sources, nil
}

func configPaths() []string {
	paths := splitValues(os.Getenv(configPathEnv))
	if len(paths) == 0 {
		return []string{"config/app.example.yaml"}
	}
	return paths
}

func nacosSource() (linker.Source, error) {
	dataID := strings.TrimSpace(os.Getenv("LINKER_V3_EXAMPLE_NACOS_DATA_ID"))
	if dataID == "" {
		return nil, nil
	}
	client := nacos.DefaultClientConfig()
	client.Host = strings.TrimSpace(os.Getenv("LINKER_V3_EXAMPLE_NACOS_HOST"))
	client.Username = strings.TrimSpace(os.Getenv("LINKER_V3_EXAMPLE_NACOS_USERNAME"))
	client.Password = os.Getenv("LINKER_V3_EXAMPLE_NACOS_PASSWORD")
	client.NamespaceID = strings.TrimSpace(os.Getenv("LINKER_V3_EXAMPLE_NACOS_NAMESPACE"))
	if value := strings.TrimSpace(os.Getenv("LINKER_V3_EXAMPLE_NACOS_PORT")); value != "" {
		port, err := strconv.ParseUint(value, 10, 16)
		if err != nil || port == 0 {
			return nil, fmt.Errorf("LINKER_V3_EXAMPLE_NACOS_PORT 必须是 1 到 65535 的端口")
		}
		client.Port = port
	}
	return nacos.Config(
		dataID,
		nacos.Group(strings.TrimSpace(os.Getenv("LINKER_V3_EXAMPLE_NACOS_GROUP"))),
		nacos.Bootstrap(client),
	), nil
}

func splitValues(value string) []string {
	ret := make([]string, 0, strings.Count(value, ",")+1)
	for part := range strings.SplitSeq(value, ",") {
		if part = strings.TrimSpace(part); part != "" {
			ret = append(ret, part)
		}
	}
	return ret
}
