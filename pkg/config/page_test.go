package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPage_EveryDuration(t *testing.T) {
	// Arrange
	testCases := []struct {
		In       string
		Expected time.Duration
	}{
		{
			In:       "",
			Expected: DefaultPageEvery,
		}, {
			In:       "not a duration",
			Expected: DefaultPageEvery,
		}, {
			In:       "0.1s",
			Expected: MinimumPageEvery,
		}, {
			In:       "3s",
			Expected: MinimumPageEvery,
		}, {
			In:       "30s",
			Expected: 30 * time.Second,
		},
	}

	p := &Page{}
	for _, testCase := range testCases {
		p.Every = testCase.In

		// Act
		result := p.EveryDuration()

		// Assert
		require.Equal(t, testCase.Expected, result)
	}
}
