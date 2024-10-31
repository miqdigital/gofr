package terminal

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"gofr.dev/pkg/gofr/testutil"
)

func TestProgressBar_SuccessCases(t *testing.T) {
	total := int64(100)

	var out bytes.Buffer
	stream := &output{terminal{isTerminal: true, fd: 1}, &out}
	bar := NewProgressBar(stream, total)

	// Mock terminal size
	bar.tWidth = 120

	// Increment the progress bar
	bar.Incr(10)

	// Verify the output
	expectedOutput := "\r[████░                                             ] 10.000%"
	assert.Equal(t, expectedOutput, out.String())

	// Increment the progress bar to completion
	bar.Incr(total - 10)

	// Verify the completion output
	expectedCompletion := "\r[████░                                             ] 10.000%\r" +
		"[██████████████████████████████████████████████████] 100.000%\n"
	assert.Equal(t, expectedCompletion, out.String())
}

func TestProgressBar_Fail(t *testing.T) {
	out := testutil.StdoutOutputForFunc(func() {
		var out bytes.Buffer
		stream := &output{terminal{isTerminal: true, fd: 1}, &out}
		bar := NewProgressBar(stream, int64(-1))

		assert.Equal(t, bar.total, int64(0))
	})

	assert.Contains(t, out, "error initializing progress bar, total should be > 0")
}

func TestProgressBar_Incr(t *testing.T) {
	var out bytes.Buffer
	stream := &output{terminal{isTerminal: true, fd: 1}, &out}
	bar := NewProgressBar(stream, 100)
	// doing this as while calculating terminal size the code will not
	// be able to determine it's width since we are not attacting an actual
	// terminal for testing
	bar.tWidth = 120

	// Increment the progress by 20
	b := bar.Incr(int64(20))
	assert.True(t, b)
	assert.Equal(t, int64(20), bar.current)

	expectedOut := "\r[█████████░                                        ] 20.000%"
	assert.Equal(t, expectedOut, out.String())

	bar.Incr(int64(100))
	expectedOut = "\r[█████████░                                        ] 20.000%\r" +
		"[██████████████████████████████████████████████████] 100.000%\n"
	assert.Equal(t, expectedOut, out.String())
}

func TestProgressBar_getString(t *testing.T) {
	testCases := []struct {
		desc        string
		current     int64
		tWidth      int
		total       int64
		expectedOut string
	}{
		{
			desc:        "current and total negative",
			current:     -1,
			total:       -1,
			tWidth:      120,
			expectedOut: "",
		},
		{
			desc:        "terminal width < 110",
			current:     20,
			total:       100,
			expectedOut: "20.000%",
		},
		{
			desc:        "0% progress, 50 spaces",
			tWidth:      120,
			current:     0,
			total:       100,
			expectedOut: "[                                                  ] 0.000%",
		},
		{
			desc:        "100% progress",
			tWidth:      120,
			current:     100,
			total:       100,
			expectedOut: "[██████████████████████████████████████████████████] 100.000%",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			p := ProgressBar{
				current: tc.current,
				total:   tc.total,
				tWidth:  tc.tWidth,
			}

			out := p.getString()

			assert.Equal(t, tc.expectedOut, out)
		})
	}
}
