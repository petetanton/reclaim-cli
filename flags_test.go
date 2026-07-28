package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// A fixed zone keeps bare-date parsing deterministic and independent of the
// machine's TZ and of DST transitions.
var testZone = time.FixedZone("BST", 60*60)

func testNow() time.Time {
	return time.Date(2026, 7, 28, 10, 30, 0, 0, testZone)
}

func TestParseStartSpec(t *testing.T) {
	now := testNow()

	tests := map[string]struct {
		spec string
		want time.Time
	}{
		"now":                    {NOW, now},
		"one day preset":         {IN_ONE_DAY, time.Date(2026, 7, 29, 10, 30, 0, 0, testZone)},
		"two days preset":        {IN_TWO_DAYS, time.Date(2026, 7, 30, 10, 30, 0, 0, testZone)},
		"one week preset":        {IN_ONE_WEEK, time.Date(2026, 8, 4, 10, 30, 0, 0, testZone)},
		"generalised days":       {"in 3 days", time.Date(2026, 7, 31, 10, 30, 0, 0, testZone)},
		"generalised weeks":      {"in 2 weeks", time.Date(2026, 8, 11, 10, 30, 0, 0, testZone)},
		"generalised hours":      {"in 6 hours", time.Date(2026, 7, 28, 16, 30, 0, 0, testZone)},
		"singular unit":          {"in 1 week", time.Date(2026, 8, 4, 10, 30, 0, 0, testZone)},
		"zero is now":            {"in 0 days", now},
		"mixed case":             {"In 1 Day", time.Date(2026, 7, 29, 10, 30, 0, 0, testZone)},
		"extra whitespace":       {"  in   1   day  ", time.Date(2026, 7, 29, 10, 30, 0, 0, testZone)},
		"bare date starts at 00": {"2026-08-04", time.Date(2026, 8, 4, 0, 0, 0, 0, testZone)},
		"rfc3339 utc":            {"2026-08-04T09:15:00Z", time.Date(2026, 8, 4, 9, 15, 0, 0, time.UTC)},
		"rfc3339 with offset":    {"2026-08-04T09:15:00+01:00", time.Date(2026, 8, 4, 8, 15, 0, 0, time.UTC)},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := parseStartSpec(test.spec, now)
			require.NoError(t, err)
			assert.True(t, test.want.Equal(got), "want %s, got %s", test.want, got)
		})
	}
}

// A task due at midnight would be overdue for the whole of its due date, so a
// bare --due date resolves to the end of that day instead of the start.
func TestParseDueSpec_bareDateIsEndOfDay(t *testing.T) {
	got, err := parseDueSpec("2026-08-04", testNow())
	require.NoError(t, err)

	assert.Equal(t, time.Date(2026, 8, 4, 23, 59, 0, 0, testZone), got)
}

func TestParseDueSpec_acceptsRelativeAndAbsolute(t *testing.T) {
	now := testNow()

	relative, err := parseDueSpec("in 1 week", now)
	require.NoError(t, err)
	assert.Equal(t, time.Date(2026, 8, 4, 10, 30, 0, 0, testZone), relative)

	absolute, err := parseDueSpec("2026-08-14T17:30:00Z", now)
	require.NoError(t, err)
	assert.True(t, time.Date(2026, 8, 14, 17, 30, 0, 0, time.UTC).Equal(absolute))
}

func TestParseTimeSpec_rejectsGarbage(t *testing.T) {
	specs := []string{
		"",
		"tomorrow",
		"next tuesday",
		"in 1 fortnight",
		"in a day",
		"in days",
		"in -1 days",
		"1 day",
		"2026-13-45",
		"04/08/2026",
		"2026-08-04 09:15",
	}

	for _, spec := range specs {
		t.Run(spec, func(t *testing.T) {
			_, err := parseStartSpec(spec, testNow())
			require.Error(t, err)
			assert.Contains(t, err.Error(), "invalid time")
		})
	}
}

func TestValidateIntChoice(t *testing.T) {
	require.NoError(t, validateIntChoice("mins", 30, minsOptions))
	require.NoError(t, validateIntChoice("mins", 240, minsOptions))

	err := validateIntChoice("mins", 37, minsOptions)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid --mins 37")
	assert.Contains(t, err.Error(), "15, 30, 45, 60, 90, 120, 180, 240")
}

func TestValidateStringChoice(t *testing.T) {
	require.NoError(t, validateStringChoice("priority", "P3", priorityOptions))

	// Deliberately case-sensitive: the API only accepts the upper-case form.
	err := validateStringChoice("priority", "p3", priorityOptions)
	require.Error(t, err)
	assert.Contains(t, err.Error(), `invalid --priority "p3"`)
	assert.Contains(t, err.Error(), "P1, P2, P3, P4")
}

// The flags exist to replace the pickers, so a flag value the picker would not
// have offered must be rejected. If these lists ever diverge, --mins 37 would
// reach the API as an odd chunk count.
func TestOptionListsMatchThePrompts(t *testing.T) {
	assert.Equal(t, []int{15, 30, 45, 60, 90, 120, 180, 240}, minsOptions)
	assert.Equal(t, []int{15, 30, 45, 60}, minChunkOptions)
	assert.Equal(t, []string{"P1", "P2", "P3", "P4"}, priorityOptions)
	assert.Equal(t, []string{"now", "in 1 day", "in 2 days", "in 1 week"}, startOptions)

	// Every preset the picker offers must parse, or choosing it interactively
	// would fail after the task had already been created.
	for _, option := range startOptions {
		_, err := parseStartSpec(option, testNow())
		assert.NoError(t, err, "picker offers %q but it does not parse", option)
	}
}

func TestErrNotInteractive_namesTheFlag(t *testing.T) {
	err := errNotInteractive("mins")
	assert.Contains(t, err.Error(), "--mins is required")
	assert.Contains(t, err.Error(), "--non-interactive")
}

func TestAsStrings(t *testing.T) {
	assert.Equal(t, []string{"15", "30", "45", "60"}, asStrings(minChunkOptions))
}
