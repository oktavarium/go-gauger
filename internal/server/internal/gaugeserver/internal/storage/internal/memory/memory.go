package memory

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/oktavarium/go-gauger/internal/server/internal/gaugeserver/internal/storage/internal/memory/archive"
	"github.com/oktavarium/go-gauger/internal/server/internal/logger"
	"go.uber.org/zap"
)

type storage struct {
	gauge   map[string]float64
	counter map[string]int64
	archive archive.FileArchive
	sync    bool
	mx      sync.RWMutex
}

func NewStorage(
	filename string,
	restore bool,
	timeout time.Duration,
) (*storage, error) {
	s := &storage{
		gauge:   map[string]float64{},
		counter: map[string]int64{},
		archive: archive.NewFileArchive(filename),
		sync:    timeout == 0,
	}

	if restore {
		err := s.restore()
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				return nil,
					fmt.Errorf("failed to restore data from file: %w", err)
			}
		}
	}

	if !s.sync {
		go func() {
			ticker := time.NewTicker(timeout)
			for range ticker.C {
				if err := s.save(); err != nil {
					logger.Logger().Error("error",
						zap.String("func", "NewStorage"),
						zap.Error(err),
					)
				}
			}
		}()
	}

	return s, nil
}

func (s *storage) Ping(context.Context) error {
	return nil
}
