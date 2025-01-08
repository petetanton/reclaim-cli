package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
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
				logrus.Infof("task %s created with id %d", title, task.Id)

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
			Name:        "dedupe",
			Description: "Deduplicate tasks with the same name (usually tasks that were created via automation)",
			Action: func(c *cli.Context) error {
				tasks, err := client.GetTasks()
				if err != nil {
					return err
				}

				logrus.Infof("found %d tasks", len(tasks))

				taskMap := make(map[string][]*reclaim.Task)
				for _, task := range tasks {
					if task.Status != "COMPLETE" && task.Status != "ARCHIVED" {
						taskMap[task.Title] = append(taskMap[task.Title], task)
					}
				}

				wg := sync.WaitGroup{}
				for title, dupeTasks := range taskMap {
					if len(dupeTasks) > 1 {
						wg.Add(1)

						go dedupe(client, title, dupeTasks, &wg)
					}
				}

				wg.Wait()
				return nil
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

func dedupe(client *reclaim.Client, title string, dupeTasks []*reclaim.Task, wg *sync.WaitGroup) {
	defer wg.Done()
	logrus.Infof("deduping %d %s", dupeTasks[0].Id, title)
	chunksRemaining := 0
	chunksRequired := 0
	for i := 1; i < len(dupeTasks); i++ {
		chunksRemaining += dupeTasks[i].TimeChunksRemaining
		chunksRequired += dupeTasks[i].TimeChunksRequired
		err := client.DeleteTask(dupeTasks[i].Id)
		if err != nil {
			logrus.Error(err)
			return
		}
	}

	mainTask := dupeTasks[0]
	mainTask.TimeChunksRemaining += chunksRemaining
	mainTask.TimeChunksRequired += chunksRequired
	updatedTask, err := client.UpdateTask(mainTask)
	if err != nil {
		logrus.Error(err)
		return
	}

	logrus.Infof("task %s updated with %d chunks remaining", title, updatedTask.TimeChunksRemaining)
}
