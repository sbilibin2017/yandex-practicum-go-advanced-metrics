package services

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	internalErrors "github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/errors"

	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/types"
)

func TestMetricGetService_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGetter := NewMockMetricGetGetter(ctrl)
	svc := NewMetricGetService(mockGetter)

	ctx := context.Background()
	testID := types.MetricID{ID: "metric1", MType: "gauge"}

	t.Run("success returns metric", func(t *testing.T) {
		expectedMetric := &types.Metrics{
			ID:    "metric1",
			MType: "gauge",
			// add other fields as needed
		}

		mockGetter.EXPECT().Get(ctx, testID).Return(expectedMetric, nil)

		result, err := svc.Get(ctx, testID)
		assert.NoError(t, err)
		assert.Equal(t, expectedMetric, result)
	})

	t.Run("getter returns error", func(t *testing.T) {
		mockGetter.EXPECT().Get(ctx, testID).Return(nil, errors.New("some error"))

		result, err := svc.Get(ctx, testID)
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("getter returns nil metric", func(t *testing.T) {
		mockGetter.EXPECT().Get(ctx, testID).Return(nil, nil)

		result, err := svc.Get(ctx, testID)
		assert.ErrorIs(t, err, internalErrors.ErrMetricNotFound)
		assert.Nil(t, result)
	})
}
