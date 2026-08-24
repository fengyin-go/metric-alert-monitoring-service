// Package service 实现业务逻辑层。
package service

import (
	"monitoring/internal/config"
	"monitoring/internal/store"
	"monitoring/pkg/logger"
)

// Service 业务服务聚合。
type Service struct {
	store store.Store
	log   *logger.Logger
	cfg   *config.Config
}

// New 创建业务服务。
func New(st store.Store, log *logger.Logger, cfg *config.Config) *Service {
	return &Service{store: st, log: log, cfg: cfg}
}
