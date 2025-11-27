package openai

import (
	"errors"
	"testing"
)

func TestShouldRetry(t *testing.T) {
	cs := []struct {
		e error
		w bool
	}{
		{errors.New("stream error"), true},
	}

	for i, c := range cs {
		a := shouldRetry(c.e)
		if a != c.w {
			t.Errorf("#%d shouldRetry(%v) = %v, want %v", i, c.e, a, c.w)
		}
	}
}
