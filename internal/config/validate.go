package config

import "fmt"

func (c Config) Validate() error {
	if c.HTTP.Addr == "" {
		return fmt.Errorf("HTTP 监听地址不能为空")
	}
	if c.GRPC.Addr == "" {
		return fmt.Errorf("gRPC 监听地址不能为空")
	}
	client := c.TTSClient.Normalize()
	if client.Target == "" {
		return fmt.Errorf("TTS client target 不能为空")
	}
	if c.PostgreSQL.Host == "" {
		return fmt.Errorf("PostgreSQL host 不能为空")
	}
	if c.PostgreSQL.User == "" {
		return fmt.Errorf("PostgreSQL user 不能为空")
	}
	if c.PostgreSQL.Password == "" {
		return fmt.Errorf("PostgreSQL password 不能为空")
	}
	if c.PostgreSQL.DBName == "" {
		return fmt.Errorf("PostgreSQL db_name 不能为空")
	}
	return nil
}
