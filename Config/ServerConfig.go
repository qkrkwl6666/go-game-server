package Config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type RedisConfig struct {
	IP   string `json:"IP"`
	Port int    `json:"Port"`
	Pwd  string `json:"Pwd"`
}

type MySQLConfig struct {
	IP   string `json:"IP"`
	Port int    `json:"Port"`
	ID   string `json:"ID"`
	PW   string `json:"PW"`
}

type ServerConfig struct {
	Redis    RedisConfig `json:"REDIS"`
	MasterDB MySQLConfig `json:"MasterDB"`
}

func (config *ServerConfig) ServerConfigLoad() error {

	// 작업 디렉토리 경로

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("getwd error: %v", err)
	}
	fmt.Println("working dir:", wd)

	cfgPath := filepath.Join(wd, "Config", "ServerConfig.json")

	b, err := os.ReadFile(cfgPath)
	if err != nil {
		return fmt.Errorf("read config: %w", err)
	}

	dec := json.NewDecoder(bytes.NewReader(b))
	dec.DisallowUnknownFields() // JSON에 정의되지 않은 필드 있으면 에러
	if err := dec.Decode(&config); err != nil {
		return fmt.Errorf("parse config: %w", err)
	}

	return nil
}
