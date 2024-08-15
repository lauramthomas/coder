package terraform

import (
	"bufio"
	"bytes"
	_ "embed"
	"slices"
	"testing"

	"github.com/cespare/xxhash"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/txtar"
	"google.golang.org/protobuf/encoding/protojson"
	protobuf "google.golang.org/protobuf/proto"

	"github.com/coder/coder/v2/coderd/database"
	"github.com/coder/coder/v2/provisionersdk/proto"
)

var (
	//go:embed testdata/timings-aggregation/simple.txtar
	inputSimple []byte
	//go:embed testdata/timings-aggregation/init.txtar
	inputInit []byte
	//go:embed testdata/timings-aggregation/error.txtar
	inputError []byte
	//go:embed testdata/timings-aggregation/complete.txtar
	inputComplete []byte
)

func TestAggregation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input []byte
	}{
		{
			name:  "init",
			input: inputInit,
		},
		{
			name:  "simple",
			input: inputSimple,
		},
		{
			name:  "error",
			input: inputError,
		},
		{
			name:  "complete",
			input: inputComplete,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// txtar is a text-based archive format used in the stdlib for simple and elegant tests.
			//
			// We ALWAYS expect that the archive contains two or more "files":
			// 	 1. JSON logs generated by a terraform execution, one per line, *one file per stage*
			//   N. Expected resulting timings in JSON form, one per line
			arc := txtar.Parse(tc.input)
			require.GreaterOrEqual(t, len(arc.Files), 2)

			t.Logf("%s: %s", t.Name(), arc.Comment)

			var actualTimings []*proto.Timing
			expectedTimings := arc.Files[len(arc.Files)-1]

			for i := 0; i < len(arc.Files)-1; i++ {
				file := arc.Files[i]
				stage := database.ProvisionerJobTimingStage(file.Name)
				require.Truef(t, stage.Valid(), "%q is not a valid stage name; acceptable values: %v",
					file.Name, database.AllProvisionerJobTimingStageValues())

				agg := newTimingAggregator(stage)
				extractAllSpans(t, file.Data, agg)
				actualTimings = append(actualTimings, agg.aggregate()...)
			}

			stableSortTimings(t, actualTimings) // To reduce flakiness.
			require.True(t, timingsAreEqual(t, expectedTimings.Data, actualTimings))
		})
	}
}

func timingsAreEqual(t *testing.T, input []byte, actual []*proto.Timing) bool {
	t.Helper()

	// Parse the input into *proto.Timing structs.
	var expected []*proto.Timing
	scanner := bufio.NewScanner(bytes.NewBuffer(input))
	for scanner.Scan() {
		line := scanner.Bytes()

		var msg proto.Timing
		require.NoError(t, protojson.Unmarshal(line, &msg))

		expected = append(expected, &msg)
	}
	require.NoError(t, scanner.Err())

	// Shortcut check.
	if len(expected)+len(actual) == 0 {
		t.Logf("both timings are empty")
		return true
	}

	// Shortcut check.
	if len(expected) != len(actual) {
		t.Logf("timings lengths are not equal: %d != %d", len(expected), len(actual))
		printExpectation(t, actual)
		return false
	}

	// Compare each element; both are expected to be sorted in a stable manner.
	for i := 0; i < len(expected); i++ {
		ex := expected[i]
		ac := actual[i]
		if !protobuf.Equal(ex, ac) {
			t.Logf("timings are not equivalent: %q != %q", ex.String(), ac.String())
			printExpectation(t, actual)
			return false
		}
	}

	return true
}

func extractAllSpans(t *testing.T, input []byte, aggregator *timingAggregator) {
	t.Helper()

	scanner := bufio.NewScanner(bytes.NewBuffer(input))
	for scanner.Scan() {
		line := scanner.Bytes()
		log := parseTerraformLogLine(line)
		if log == nil {
			continue
		}

		ts, span, err := extractTimingSpan(log)
		if err != nil {
			// t.Logf("%s: failed span extraction on line: %q", err, line)
			continue
		}

		require.NotZerof(t, ts, "failed on line: %q", line)
		require.NotNilf(t, span, "failed on line: %q", line)

		aggregator.ingest(ts, span)
	}

	require.NoError(t, scanner.Err())
}

func printExpectation(t *testing.T, actual []*proto.Timing) {
	t.Helper()

	t.Log("expected:")
	for _, a := range actual {
		printTiming(t, a)
	}
}

func printTiming(t *testing.T, timing *proto.Timing) {
	t.Helper()

	marshaler := protojson.MarshalOptions{
		Multiline: false, // Ensure it's set to false for single-line JSON
		Indent:    "",    // No indentation
	}

	out, err := marshaler.Marshal(timing)
	assert.NoError(t, err)
	t.Logf("%s", out)
}

func stableSortTimings(t *testing.T, timings []*proto.Timing) {
	slices.SortStableFunc(timings, func(a, b *proto.Timing) int {
		if a == nil || b == nil || a.Start == nil || b.Start == nil {
			return 0
		}

		if a.Start.AsTime().Equal(b.Start.AsTime()) {
			// Special case: when start times are equal, we need to keep the ordering stable, so we hash both entries
			// and sort based on that (since end times could be equal too, in principle).
			ah := xxhash.Sum64String(a.String())
			bh := xxhash.Sum64String(b.String())

			if ah == bh {
				// WTF.
				t.Logf("identical timings detected!")
				printTiming(t, a)
				printTiming(t, b)
				return 0
			}

			if ah < bh {
				return -1
			}

			return 1
		}

		if a.Start.AsTime().Before(b.Start.AsTime()) {
			return -1
		}

		return 1
	})
}
