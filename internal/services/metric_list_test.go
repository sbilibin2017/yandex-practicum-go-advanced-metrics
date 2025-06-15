package services

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/types"
)

func TestMetricListService_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLister := NewMockMetricListLister(ctrl)
	svc := NewMetricListService(mockLister)

	ctx := context.Background()

	tests := []struct {
		name       string
		mockReturn []types.Metrics
		mockErr    error
		want       []types.Metrics
		wantErr    bool
	}{
		{
			name: "success with metrics",
			mockReturn: []types.Metrics{
				{ID: "m1", MType: "gauge"},
				{ID: "m2", MType: "counter"},
			},
			mockErr: nil,
			want: []types.Metrics{
				{ID: "m1", MType: "gauge"},
				{ID: "m2", MType: "counter"},
			},
			wantErr: false,
		},
		{
			name:       "success with empty list",
			mockReturn: []types.Metrics{},
			mockErr:    nil,
			want:       []types.Metrics{},
			wantErr:    false,
		},
		{
			name:       "error from lister",
			mockReturn: nil,
			mockErr:    errors.New("some error"),
			want:       nil,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLister.EXPECT().List(ctx).Return(tt.mockReturn, tt.mockErr)

			got, err := svc.List(ctx)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
