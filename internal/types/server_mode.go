package types

import (
	"fmt"
)

type ServerMode string

const (
	ServerModeDev  ServerMode = "dev"
	ServerModeProd ServerMode = "prod"
	ServerModeTest ServerMode = "test"
)

func (m ServerMode) String() string {
	return string(m)
}

func ParseServerMode(mode string) (ServerMode, error) {
	switch mode {
	case "dev":
		return ServerModeDev, nil
	case "prod":
		return ServerModeProd, nil
	case "test":
		return ServerModeTest, nil
	}
	return ServerModeDev, fmt.Errorf("invalid server mode: %s", mode)
}
