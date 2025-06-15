package repositories

import (
	"context"
	"testing"

	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/engines"
	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestMetricMemoryListRepository_List(t *testing.T) {
	storage := engines.NewMemoryStorage[types.MetricID, types.Metrics]()
	repo := NewMetricMemoryListRepository(storage)

	tests := []struct {
		name          string
		setupData     []types.Metrics
		expectedOrder []string // expected sorted metric IDs
	}{
		{
			name:          "empty storage returns empty slice",
			setupData:     []types.Metrics{},
			expectedOrder: []string{},
		},
		{
			name: "single metric returns it",
			setupData: []types.Metrics{
				{ID: "metric1", MType: "gauge"},
			},
			expectedOrder: []string{"metric1"},
		},
		{
			name: "multiple metrics returned sorted by ID",
			setupData: []types.Metrics{
				{ID: "metricB", MType: "gauge"},
				{ID: "metricA", MType: "counter"},
				{ID: "metricC", MType: "gauge"},
			},
			expectedOrder: []string{"metricA", "metricB", "metricC"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset storage data with locking
			storage.Mu.Lock()
			storage.Data = make(map[types.MetricID]types.Metrics)
			for _, metric := range tt.setupData {
				key := types.MetricID{ID: metric.ID, MType: metric.MType}
				storage.Data[key] = metric
			}
			storage.Mu.Unlock()

			result, err := repo.List(context.Background())
			assert.NoError(t, err)

			assert.Len(t, result, len(tt.expectedOrder))
			for i, expectedID := range tt.expectedOrder {
				assert.Equal(t, expectedID, result[i].ID)
			}
		})
	}
}
