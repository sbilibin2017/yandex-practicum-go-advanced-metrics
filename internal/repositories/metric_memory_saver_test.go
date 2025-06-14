package repositories

import (
	"context"
	"sync"
	"testing"

	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestMetricMemorySaverRepository_Send_Save(t *testing.T) {
	repo := &MetricMemorySaverRepository{
		data: make(map[types.MetricID]types.Metrics),
		mu:   sync.RWMutex{},
	}

	delta := int64(42)
	value := 3.14

	tests := []struct {
		name    string
		input   types.Metrics
		wantErr bool
	}{
		{
			name: "save counter metric with delta",
			input: types.Metrics{
				MType: types.Counter,
				ID:    "counter1",
				Delta: &delta,
			},
			wantErr: false,
		},
		{
			name: "save gauge metric with value",
			input: types.Metrics{
				MType: types.Gauge,
				ID:    "gauge1",
				Value: &value,
			},
			wantErr: false,
		},
		{
			name: "save metric with both delta and value",
			input: types.Metrics{
				MType: types.Gauge,
				ID:    "mixed1",
				Delta: &delta,
				Value: &value,
			},
			wantErr: false,
		},
		{
			name: "save metric with no delta or value",
			input: types.Metrics{
				MType: types.Counter,
				ID:    "empty1",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Send(context.Background(), tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// Check that metric is stored
				repo.mu.RLock()
				stored, ok := repo.data[types.MetricID{ID: tt.input.ID, MType: tt.input.MType}]
				repo.mu.RUnlock()
				assert.True(t, ok, "metric should be stored")
				assert.Equal(t, tt.input, stored, "stored metric should match input")
			}
		})
	}

	t.Run("Save returns nil", func(t *testing.T) {
		err := repo.Save(context.Background())
		assert.NoError(t, err)
	})
}
