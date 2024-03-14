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

func TestPage_ValidateNormalize_NoError(t *testing.T) {
	// Arrange
	testCases := []struct {
		Regexp        map[string]string
		ErrorContains string
	}{
		{
			Regexp: map[string]string{
				`[[:blank:]]+`:   " ",
				` +`:             " ",
				`[[:blank:]]+\n`: "",
			},
		}, {
			Regexp: map[string]string{
				"\\j": "",
			},
			ErrorContains: "invalid escape sequence",
		}, {
			Regexp: map[string]string{
				`[[:blank:]]+`: " ",
				"\\j":          "",
			},
			ErrorContains: "invalid escape sequence",
		}, {
			Regexp: map[string]string{
				`[[:blank:]]+`:   " ",
				"\\j":            "",
				`[[:blank:]]+\n`: "",
			},
			ErrorContains: "invalid escape sequence",
		},
	}

	p := &Page{}
	for i, testCase := range testCases {
		t.Logf("Testing case %d", i)
		p.Normalize = testCase.Regexp
		p.normalizeRegex = nil

		// Act
		err := p.ValidateNormalize()

		// Assert
		if testCase.ErrorContains == "" {
			require.NoError(t, err)
			require.NotNil(t, p.normalizeRegex)
			require.Equal(t, len(testCase.Regexp), len(p.normalizeRegex))
			for r, to := range p.normalizeRegex {
				require.Contains(t, testCase.Regexp, r.String())
				require.Equal(t, testCase.Regexp[r.String()], to)
			}
		} else {
			require.ErrorContains(t, err, testCase.ErrorContains)
			require.Nil(t, p.normalizeRegex)
		}
	}
}

func TestPage_NormalizeString(t *testing.T) {
	// Arrange
	p := &Page{
		Normalize: map[string]string{
			`[[:blank:]]+`: "x",
			` +`:           "y",
		},
	}
	testString := "               "
	expectedString := "xyxyx"
	err := p.ValidateNormalize()
	require.NoError(t, err)

	// Act
	resultString := p.NormalizeString(testString)

	// Assert
	require.Equal(t, expectedString, resultString)
}
