package runners

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestRunServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name    string
		setup   func() (context.Context, Server)
		wantErr error
	}{
		{
			name: "ListenAndServe returns nil immediately",
			setup: func() (context.Context, Server) {
				ctx := context.Background()
				mockSrv := NewMockServer(ctrl)

				mockSrv.EXPECT().ListenAndServe().Return(nil).Times(1)

				// Shutdown should not be called in this scenario
				return ctx, mockSrv
			},
			wantErr: nil,
		},
		{
			name: "ListenAndServe returns error immediately",
			setup: func() (context.Context, Server) {
				ctx := context.Background()
				mockSrv := NewMockServer(ctrl)

				mockSrv.EXPECT().ListenAndServe().Return(errors.New("listen error")).Times(1)

				return ctx, mockSrv
			},
			wantErr: errors.New("listen error"),
		},
		{
			name: "context canceled triggers Shutdown successfully",
			setup: func() (context.Context, Server) {
				ctx, cancel := context.WithCancel(context.Background())
				mockSrv := NewMockServer(ctrl)

				// ListenAndServe blocks until ctx is canceled
				mockSrv.EXPECT().ListenAndServe().DoAndReturn(func() error {
					<-ctx.Done()
					return nil
				}).Times(1)

				mockSrv.EXPECT().Shutdown(gomock.Any()).Return(nil).Times(1)

				go func() {
					time.Sleep(10 * time.Millisecond)
					cancel()
				}()

				return ctx, mockSrv
			},
			wantErr: context.Canceled,
		},
		{
			name: "context canceled triggers Shutdown with error",
			setup: func() (context.Context, Server) {
				ctx, cancel := context.WithCancel(context.Background())
				mockSrv := NewMockServer(ctrl)

				mockSrv.EXPECT().ListenAndServe().DoAndReturn(func() error {
					<-ctx.Done()
					return nil
				}).Times(1)

				mockSrv.EXPECT().Shutdown(gomock.Any()).Return(errors.New("shutdown error")).Times(1)

				go func() {
					time.Sleep(10 * time.Millisecond)
					cancel()
				}()

				return ctx, mockSrv
			},
			wantErr: errors.New("shutdown error"),
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			ctx, srv := tt.setup()
			err := RunServer(ctx, srv)
			if tt.wantErr == nil {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.EqualError(t, err, tt.wantErr.Error())
			}
		})
	}
}
