package main

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/petetanton/reclaim-cli/pkg/input"
)

// The option lists offered by the interactive prompts. Flag values are validated
// against these same lists so that `--mins 37` fails loudly rather than being
// turned into an odd chunk count.
var (
	minsOptions     = []int{15, 30, 45, 60, 90, 120, 180, 240}
	minChunkOptions = []int{15, 30, 45, 60}
	priorityOptions = []string{"P1", "P2", "P3", "P4"}
	startOptions    = []string{NOW, IN_ONE_DAY, IN_TWO_DAYS, IN_ONE_WEEK}
)

// relativeSpec matches the generalised form of the IN_ONE_DAY-style presets, so
// that "in 3 days" works as well as the four values the picker offers.
var relativeSpec = regexp.MustCompile(`^in (\d+) (hour|hours|day|days|week|weeks)$`)

// errNotInteractive explains why we are refusing to prompt. An agent driving the
// CLI gets this immediately instead of hanging on a prompt it cannot answer.
func errNotInteractive(flag string) error {
	return fmt.Errorf("--%s is required when not running interactively (stdin is not a terminal, or --non-interactive was passed)", flag)
}

func validateIntChoice(flag string, value int, options []int) error {
	if slices.Contains(options, value) {
		return nil
	}

	allowed := make([]string, 0, len(options))
	for _, option := range options {
		allowed = append(allowed, strconv.Itoa(option))
	}
	return fmt.Errorf("invalid --%s %d: must be one of %s", flag, value, strings.Join(allowed, ", "))
}

func validateStringChoice(flag, value string, options []string) error {
	if slices.Contains(options, value) {
		return nil
	}
	return fmt.Errorf("invalid --%s %q: must be one of %s", flag, value, strings.Join(options, ", "))
}

// asStrings renders an int option list for the interactive picker.
func asStrings(options []int) []string {
	out := make([]string, 0, len(options))
	for _, option := range options {
		out = append(out, strconv.Itoa(option))
	}
	return out
}

// resolveIntChoice returns the flag's value when it was given, prompts for it
// when prompting is possible, and otherwise fails fast. This is the precedence
// that keeps the interactive experience unchanged for a human passing no flags.
func resolveIntChoice(c *cli.Context, flag, question string, options []int, interactive bool) (int, error) {
	if c.IsSet(flag) {
		value := c.Int(flag)
		if err := validateIntChoice(flag, value, options); err != nil {
			return 0, err
		}
		return value, nil
	}

	if !interactive {
		return 0, errNotInteractive(flag)
	}

	return strconv.Atoi(input.AskSelect(question, asStrings(options)))
}

func resolveStringChoice(c *cli.Context, flag, question string, options []string, interactive bool) (string, error) {
	if c.IsSet(flag) {
		value := c.String(flag)
		if err := validateStringChoice(flag, value, options); err != nil {
			return "", err
		}
		return value, nil
	}

	if !interactive {
		return "", errNotInteractive(flag)
	}

	return input.AskSelect(question, options), nil
}

// parseTimeSpec resolves a user-supplied time expression relative to now.
// Accepted forms:
//
//	now
//	in <n> hours|days|weeks       e.g. "in 1 day", "in 3 days", "in 2 weeks"
//	YYYY-MM-DD                    that date at bareDateHour:bareDateMin, local time
//	an RFC3339 timestamp          e.g. "2026-08-04T09:15:00Z"
//
// Bare dates carry no time of day, so the caller supplies one: a start defaults
// to the beginning of the day and a due date to the end of it, since a task due
// at midnight would be overdue for the whole of its due date.
func parseTimeSpec(spec string, now time.Time, bareDateHour, bareDateMin int) (time.Time, error) {
	normalised := strings.Join(strings.Fields(strings.ToLower(spec)), " ")

	if normalised == NOW {
		return now, nil
	}

	if match := relativeSpec.FindStringSubmatch(normalised); match != nil {
		count, err := strconv.Atoi(match[1])
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid time %q: %w", spec, err)
		}

		unit := strings.TrimSuffix(match[2], "s")
		switch unit {
		case "hour":
			return now.Add(time.Hour * time.Duration(count)), nil
		case "day":
			return now.AddDate(0, 0, count), nil
		case "week":
			return now.AddDate(0, 0, count*7), nil
		}
	}

	// A bare date is interpreted in the local zone, matching how someone reading
	// "due on the 4th" means their own 4th rather than UTC's.
	if date, err := time.ParseInLocation(time.DateOnly, normalised, now.Location()); err == nil {
		return time.Date(date.Year(), date.Month(), date.Day(), bareDateHour, bareDateMin, 0, 0, now.Location()), nil
	}

	// Use the original spec here: RFC3339 is case-sensitive about its T and Z.
	if timestamp, err := time.Parse(time.RFC3339, strings.TrimSpace(spec)); err == nil {
		return timestamp, nil
	}

	return time.Time{}, fmt.Errorf("invalid time %q: expected %s, a date (2006-01-02), an RFC3339 timestamp, or \"in <n> hours|days|weeks\"", spec, strings.Join(startOptions, ", "))
}

// parseStartSpec resolves a --start or --until value. A bare date starts at
// midnight.
func parseStartSpec(spec string, now time.Time) (time.Time, error) {
	return parseTimeSpec(spec, now, 0, 0)
}

// parseDueSpec resolves a --due value. A bare date is due at the end of that day.
func parseDueSpec(spec string, now time.Time) (time.Time, error) {
	return parseTimeSpec(spec, now, 23, 59)
}
