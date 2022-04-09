package plouf_modules

import (
	"os"
	"strings"

	"github.com/quantumsheep/plouf"
)

type ConfigService struct {
	plouf.Service

	config map[string]string
}

func (s *ConfigService) Init(self plouf.IInjectable) error {
	s.config = make(map[string]string)
	return nil
}

func (s *ConfigService) Get(key string) string {
	if value, ok := s.config[key]; ok {
		return value
	}

	return os.Getenv(strings.Replace(key, ".", "_", -1))
}

func (s *ConfigService) GetOr(key string, fallback string) string {
	if value := s.Get(key); value != "" {
		return value
	}

	return fallback
}

func (s *ConfigService) Set(key string, value string) {
	s.config[key] = value
}
