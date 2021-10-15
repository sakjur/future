package future_test

import (
	"context"
	"fmt"
	"github.com/sakjur/future"
	"time"
)

func ExampleNew() {
	f := future.New(context.Background(), func(ctx context.Context) (string, error) {
		return "Hello, world", nil
	})
	fmt.Println(f.Wait())
	// Output: Hello, world <nil>
}

func ExampleFuture_Cancel() {
	f := future.New(context.Background(), func(ctx context.Context) (int, error) {
		time.Sleep(20 * time.Millisecond)

		if ctx.Err() != nil {
			return 0, ctx.Err()
		}

		return 42, nil
	})
	f.Cancel()
	fmt.Println(f.Wait())
	// Output: 0 context canceled
}

func ExampleFuture_MustWait() {
	f := future.New(context.Background(), func(ctx context.Context) (int, error) {
		return 42, nil
	})

	// We know f never returns an error, ok to use Must.
	fmt.Println(22 + f.MustWait())
	// Output: 64
}

func ExampleFuture_MustWait2() {
	f := future.New(context.Background(), func(ctx context.Context) (int, error) {
		return 0, fmt.Errorf("whoops")
	})
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("panicking:", err)
		}
	}()

	// This panics!
	fmt.Println(22 + f.MustWait())
	// Output: panicking: whoops
}
