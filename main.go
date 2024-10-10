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
	"github.com/petetanton/reclaim-cli/pkg/version"
)

const (
	NOW         = "now"
	IN_ONE_DAY  = "in 1 day"
	IN_TWO_DAYS = "in 2 days"
	IN_ONE_WEEK = "in 1 week"
)

func main() {
	client := reclaim.New()
	app := cli.NewApp()
	app.Commands = []*cli.Command{
		{
			Name: "create",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "title",
					Required: true,
				},
			},
			Action: func(c *cli.Context) error {
				title := c.String("title")
				mins := input.AskSelect("how many mins for the task", []string{"15", "30", "45", "60", "75", "90", "105", "120"})
				minsInt, err := strconv.Atoi(mins)
				if err != nil {
					return err
				}
				minChunk := input.AskSelect("what is the min chunk length for the task", []string{"15", "30", "45", "60"})
				minChunkInt, err := strconv.Atoi(minChunk)
				if err != nil {
					return err
				}

				minChunkSize := minChunkInt / 15

				task, err := client.CreateTask(title, minChunkSize, minChunkSize*8, minsInt/15)
				if err != nil {
					return err
				}
				logrus.Printf("task %s created with id %d", title, task.Id)

				delay := input.AskSelect("when would you like to start this task", []string{NOW, IN_ONE_DAY, IN_TWO_DAYS, IN_ONE_WEEK})
				if delay == IN_ONE_DAY {
					return client.SnoozeTask(task.Id, time.Now().Add(time.Hour*24))
				}
				if delay == IN_TWO_DAYS {
					return client.SnoozeTask(task.Id, time.Now().Add(time.Hour*24*2))
				}
				if delay == IN_ONE_WEEK {
					return client.SnoozeTask(task.Id, time.Now().Add(time.Hour*24*7))
				}
				return nil
			},
		},
		{
			Name:        "snooze",
			Description: "Snooze a task so that it is scheduled at a later date",
			Action: func(c *cli.Context) error {
				tasks, err := client.GetTasks()
				if err != nil {
					return err
				}

				var snoozableItems []string
				for _, task := range tasks {
					if task.Status != "COMPLETE" && task.Status != "ARCHIVED" {
						//logrus.Printf("[%s] %s", task.Status, task.Title)
						snoozableItems = append(snoozableItems, fmt.Sprintf("%d, %s", task.Id, task.Title))
					}
				}

				itemToSnooze := input.AskSelect("Which task would you like to snooze", snoozableItems)

				idToSnooze, err := strconv.Atoi(strings.Split(itemToSnooze, ",")[0])
				if err != nil {
					return err
				}

				return client.SnoozeTask(idToSnooze, time.Now().Add(time.Hour*24))
			},
		},
		{
			Name: "version",
			Action: func(c *cli.Context) error {
				logrus.Print(version.Version)
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		logrus.Fatal(err)
	}
}
