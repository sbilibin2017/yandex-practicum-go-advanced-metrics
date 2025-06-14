package repositories

import (
	"sync"

	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/types"
)

var (
	mu   = &sync.RWMutex{}
	data = make(map[types.MetricID]types.Metrics)
)
