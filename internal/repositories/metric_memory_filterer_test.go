package repositories

import (
	"context"
	"sync"
	"testing"

	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestMetricMemoryFiltererRepository_Filter(t *testing.T) {
	// Prepare sample metrics
	delta1, delta2 := int64(10), int64(20)
	value1, value2 := 1.1, 2.2

	metric1 := types.Metrics{MType: types.Counter, ID: "m1", Delta: &delta1}
	metric2 := types.Metrics{MType: types.Gauge, ID: "m2", Value: &value1}
	metric3 := types.Metrics{MType: types.Gauge, ID: "m3", Value: &value2}
	metric4 := types.Metrics{MType: types.Counter, ID: "m4", Delta: &delta2}

	repo := &MetricMemoryFiltererRepository{
		data: map[types.MetricID]types.Metrics{
			{ID: metric1.ID, MType: metric1.MType}: metric1,
			{ID: metric2.ID, MType: metric2.MType}: metric2,
			{ID: metric3.ID, MType: metric3.MType}: metric3,
			{ID: metric4.ID, MType: metric4.MType}: metric4,
		},
		mu: sync.RWMutex{},
	}

	tests := []struct {
		name     string
		ids      []types.MetricID
		expected map[types.MetricID]types.Metrics
	}{
		{
			name:     "empty ids returns empty map",
			ids:      []types.MetricID{},
			expected: map[types.MetricID]types.Metrics{},
		},
		{
			name: "single existing id",
			ids:  []types.MetricID{{ID: metric1.ID, MType: metric1.MType}},
			expected: map[types.MetricID]types.Metrics{
				{ID: metric1.ID, MType: metric1.MType}: metric1,
			},
		},
		{
			name: "multiple existing ids",
			ids: []types.MetricID{
				{ID: metric1.ID, MType: metric1.MType},
				{ID: metric3.ID, MType: metric3.MType},
			},
			expected: map[types.MetricID]types.Metrics{
				{ID: metric1.ID, MType: metric1.MType}: metric1,
				{ID: metric3.ID, MType: metric3.MType}: metric3,
			},
		},
		{
			name: "mix of existing and non-existing ids",
			ids: []types.MetricID{
				{ID: metric2.ID, MType: metric2.MType},
				{ID: "nonexistent", MType: types.Counter},
			},
			expected: map[types.MetricID]types.Metrics{
				{ID: metric2.ID, MType: metric2.MType}: metric2,
			},
		},
		{
			name: "all non-existing ids",
			ids: []types.MetricID{
				{ID: "x1", MType: types.Counter},
				{ID: "x2", MType: types.Gauge},
			},
			expected: map[types.MetricID]types.Metrics{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.Filter(context.Background(), tt.ids)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}
