package runners

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewRunContext(t *testing.T) {
	ctx, cancel := NewRunContext(context.Background())
	defer cancel()

	// The returned context should not be nil
	assert.NotNil(t, ctx)

	// The cancel function should not be nil
	assert.NotNil(t, cancel)

	doneCh := make(chan struct{})
	go func() {
		<-ctx.Done()
		close(doneCh)
	}()

	// Cancel the context
	cancel()

	select {
	case <-doneCh:
		// success: context was cancelled
	case <-time.After(time.Second):
		t.Fatal("context was not cancelled after cancel function called")
	}
}
