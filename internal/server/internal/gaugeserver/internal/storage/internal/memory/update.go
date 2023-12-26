package memory

import (
	"context"
	"fmt"
)

func (s *storage) UpdateCounter(
	ctx context.Context,
	name string,
	val int64,
) (int64, error) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.counter[name] += val
	if s.sync {
		err := s.save()
		if err != nil {
			return 0, fmt.Errorf("failed to update counter: %w", err)
		}
	}
	return s.counter[name], nil
}
