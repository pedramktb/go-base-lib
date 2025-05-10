package lifecycle

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Lifecycle(t *testing.T) {
	ctx, cancel, shutdown := Context(31 * time.Second)
	done1, err := RegisterCloser(ctx)
	if err != nil {
		t.Fatal(err)
	}
	done2, err := RegisterCloser(ctx)
	if err != nil {
		t.Fatal(err)
	}
	cancel()
	go func() {
		<-ctx.Done()
		done1(errors.New("test error"))
	}()
	go func() {
		<-ctx.Done()
		done2(nil)
	}()
	errors := []error{}
	for err := range shutdown {
		if err != nil {
			errors = append(errors, err)
		}
	}
	assert.Equal(t, 1, len(errors))

}
