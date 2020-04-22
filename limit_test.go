package limit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMemLimitStart(t *testing.T) {
	t.Run("when memory limit is reached", func(t *testing.T) {
		memLimit := &MemLimit{
			ctx:       context.Background(),
			maxMemory: 0,
			onLimit: func(msg string) {
				require.Contains(t, msg, "memory limit exceeded")
			},
		}

		err := memLimit.Execute(func(ctx context.Context) {
			<-ctx.Done()
		})

		require.Error(t, err, "limits reached")
	})

	t.Run("when parent context has been canceled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		memLimit := &MemLimit{
			ctx:       ctx,
			maxMemory: 1 << 30,
			onLimit: func(msg string) {
				require.Equal(t, "context done", msg)
			},
		}

		err := memLimit.Execute(func(ctx context.Context) {
			t.Fatalf("this should not be executed")
		})

		require.Error(t, err, "limits reached")
	})

	t.Run("when limits have not been reached", func(t *testing.T) {
		memLimit := &MemLimit{ctx: context.Background(), maxMemory: 1 << 30}

		err := memLimit.Execute(func(ctx context.Context) {
			require.NoError(t, ctx.Err())
		})

		require.NoError(t, err)
	})
}
