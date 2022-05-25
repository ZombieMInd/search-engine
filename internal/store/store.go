package store

import (
	"github.com/ZombieMInd/go-logger/internal/logger"
)

type Store interface {
	Log() logger.LogRepository
}
