package user

import (
	"fmt"

	configcore "github.com/neteast-software/go-module/config"
	linker "github.com/neteast-software/linker/v3"
)

const Namespace linker.Namespace = "example/user"

type Config struct {
	TokenKey string `json:"token_key" yaml:"token_key"`
}

func (p Config) Validate() error {
	if len(p.TokenKey) < 32 {
		return fmt.Errorf("example/user token_key 至少需要 32 个字符")
	}
	return nil
}

func decodeConfig(content []byte) (Config, error) {
	var config Config
	if err := configcore.Decode(content, &config); err != nil {
		return Config{}, err
	}
	if err := config.Validate(); err != nil {
		return Config{}, err
	}
	return config, nil
}
