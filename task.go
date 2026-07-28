package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/petetanton/reclaim-cli/pkg/input"
	"github.com/petetanton/reclaim-cli/pkg/reclaim"
)

// createTask resolves every parameter from its flag when given and falls back to
// the interactive picker otherwise, so passing no flags behaves exactly as it
// always has while an agent can supply everything up front.
func createTask(client *reclaim.Client, c *cli.Context) error {
	now := time.Now()
	interactive := !c.Bool("non-interactive") && input.IsInteractive()

	mins, err := resolveIntChoice(c, "mins", "how many mins for the task", minsOptions, interactive)
	if err != nil {
		return err
	}

	minChunk, err := resolveIntChoice(c, "min-chunk", "what is the min chunk length for the task", minChunkOptions, interactive)
	if err != nil {
		return err
	}

	priority, err := resolveStringChoice(c, "priority", "what is the priority of the task", priorityOptions, interactive)
	if err != nil {
		return err
	}

	// Resolve both times before creating anything: a malformed --start or --due
	// should not leave a task behind that never got its schedule applied.
	var due *time.Time
	if c.IsSet("due") {
		parsed, err := parseDueSpec(c.String("due"), now)
		if err != nil {
			return err
		}
		due = &parsed
	}

	startGivenAsFlag := c.IsSet("start")
	startSpec := c.String("start")
	if startGivenAsFlag {
		if _, err := parseStartSpec(startSpec, now); err != nil {
			return err
		}
	} else if !interactive {
		return errNotInteractive("start")
	}

	title := c.String("title")
	minChunkSize := minChunk / 15

	task, err := client.CreateTask(&reclaim.CreateTaskRequest{
		Title:              title,
		Notes:              c.String("notes"),
		MinChunkSize:       minChunkSize,
		MaxChunkSize:       minChunkSize * 8,
		TimeChunksRequired: mins / 15,
		Priority:           reclaim.TaskPriority(priority),
		Due:                due,
	})
	if err != nil {
		return err
	}
	logrus.Infof("task %s created with id %d", title, task.Id)

	// The start prompt deliberately stays after the create call, which is where it
	// has always been in the interactive flow.
	if !startGivenAsFlag {
		startSpec = input.AskSelect("when would you like to start this task", startOptions)
	}

	startAt, err := parseStartSpec(startSpec, now)
	if err != nil {
		return err
	}

	// "now" and anything already in the past mean "do not snooze", which covers the
	// NOW preset without having to compare strings.
	if startAt.After(now) {
		task, err = client.SnoozeTask(task.Id, startAt)
		if err != nil {
			return err
		}
	}

	if c.Bool("json") {
		return writeTaskJSON(os.Stdout, task)
	}

	return nil
}

// snoozeTask picks its target from --id, from --title, or from the interactive
// list, in that order.
func snoozeTask(client *reclaim.Client, c *cli.Context) error {
	now := time.Now()
	interactive := !c.Bool("non-interactive") && input.IsInteractive()

	snoozeUntil, err := parseStartSpec(c.String("until"), now)
	if err != nil {
		return err
	}

	id, err := resolveTaskToSnooze(client, c, interactive)
	if err != nil {
		return err
	}

	task, err := client.SnoozeTask(id, snoozeUntil)
	if err != nil {
		return err
	}
	logrus.Infof("task %d snoozed until %s", task.Id, snoozeUntil.Format(time.RFC3339))

	if c.Bool("json") {
		return writeTaskJSON(os.Stdout, task)
	}

	return nil
}

func resolveTaskToSnooze(client *reclaim.Client, c *cli.Context, interactive bool) (int, error) {
	if c.IsSet("id") {
		return c.Int("id"), nil
	}

	if !c.IsSet("title") && !interactive {
		return 0, fmt.Errorf("--id or --title is required when not running interactively (stdin is not a terminal, or --non-interactive was passed)")
	}

	tasks, err := client.GetTasks([]string{})
	if err != nil {
		return 0, err
	}

	var snoozable []*reclaim.Task
	for _, task := range tasks {
		if task.Status != "COMPLETE" && task.Status != "ARCHIVED" {
			snoozable = append(snoozable, task)
		}
	}

	if c.IsSet("title") {
		return matchTaskByTitle(snoozable, c.String("title"))
	}

	var snoozableItems []string
	for _, task := range snoozable {
		snoozableItems = append(snoozableItems, fmt.Sprintf("%d, %s", task.Id, task.Title))
	}

	itemToSnooze := input.AskSelect("Which task would you like to snooze", snoozableItems)

	return strconv.Atoi(strings.Split(itemToSnooze, ",")[0])
}

// matchTaskByTitle requires exactly one match, listing the candidates otherwise,
// so an ambiguous --title never silently snoozes the wrong task.
func matchTaskByTitle(tasks []*reclaim.Task, title string) (int, error) {
	needle := strings.ToLower(title)

	var matches []*reclaim.Task
	for _, task := range tasks {
		if strings.Contains(strings.ToLower(task.Title), needle) {
			matches = append(matches, task)
		}
	}

	switch len(matches) {
	case 0:
		return 0, fmt.Errorf("no open task has a title containing %q", title)
	case 1:
		return matches[0].Id, nil
	}

	candidates := make([]string, 0, len(matches))
	for _, task := range matches {
		candidates = append(candidates, fmt.Sprintf("  %d  %s", task.Id, task.Title))
	}
	return 0, fmt.Errorf("%d open tasks have a title containing %q, pass --id for one of:\n%s", len(matches), title, strings.Join(candidates, "\n"))
}
