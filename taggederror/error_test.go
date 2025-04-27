package taggederror

import (
	"errors"
	"fmt"
	"testing"
)

func Test_X(t *testing.T) {
	e1 := errors.New("1")
	e2 := fmt.Errorf("%w - %w", errors.New("2"), e1)
	println(errors.Is(e2, e1))
}

func Test_Is(t *testing.T) {
	root := NewRoot(errors.New("Root"), "Root", 1000)

	level1Part := errors.New("Level1")
	level1 := root.Wrap(New(level1Part, "Level1"))

	level2Part := errors.New("Level2")
	level2 := level1.Wrap(level2Part)

	tests := []struct {
		name    string
		err     error
		target  error
		isError bool
	}{
		{
			name:    "equal #1",
			err:     root,
			target:  root,
			isError: true,
		},
		{
			name:    "equal #2",
			err:     level1,
			target:  level1,
			isError: true,
		},
		{
			name:    "equal #3",
			err:     level2,
			target:  level2,
			isError: true,
		},
		{
			name:    "wrapper #1",
			err:     level1,
			target:  level1Part,
			isError: true,
		},
		{
			name:    "wrapper #2",
			err:     level2,
			target:  level2Part,
			isError: true,
		},
		{
			name:    "wrapped #1",
			err:     level1Part,
			target:  level1,
			isError: false,
		},
		{
			name:    "wrapped #2",
			err:     level2Part,
			target:  level2,
			isError: false,
		},
		{
			name:    "as parent #1",
			err:     level1,
			target:  root,
			isError: true,
		},
		{
			name:    "as parent #2",
			err:     level2,
			target:  level1,
			isError: true,
		},
		{
			name:    "as grand parent",
			err:     level2,
			target:  root,
			isError: true,
		},
		{
			name:    "as child #1",
			err:     root,
			target:  level1,
			isError: false,
		},
		{
			name:    "as child #2",
			err:     level1,
			target:  level2,
			isError: false,
		},
		{
			name:    "as child #3",
			err:     root,
			target:  level2,
			isError: false,
		},
		{
			name:    "not related",
			err:     root,
			target:  errors.New("NotRoot"),
			isError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := Is(tc.err, tc.target); got != tc.isError {
				t.Errorf("IsError() = %v, want %v", got, tc.isError)
			}
		})
	}
}
