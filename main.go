package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"gitlab.com/gitlab-org/api/client-go"

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
				mins := input.AskSelect("how many mins for the task", []string{"15", "30", "45", "60", "90", "120", "180", "240"})
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
				tasks, err := client.GetTasks([]string{})
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
				tasks, err := client.GetTasks([]string{})
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

				tasks, err = client.GetTasks([]string{})
				if err != nil {
					return err
				}

				for _, task := range tasks {
					err = removeGitlabTaskIfClosed(client, task)
					if err != nil {
						return err
					}
				}
				return nil
			},
		},
		{
			Name:        "meeting",
			Description: "create a meeting",
			Action: func(c *cli.Context) error {

				links, err := client.GetScheduleLinks()
				if err != nil {
					return err
				}

				linkMap := make(map[string]*reclaim.ScheduleLink)
				linkTitles := []string{}
				for _, link := range links {
					linkMap[link.Title] = link
					linkTitles = append(linkTitles, link.Title)
				}
				linkName := input.AskSelect("Which scheduling link should we use", linkTitles)

				meetingTime, err := client.GetNextMeetingTime(linkMap[linkName].Id)
				if err != nil {
					return err
				}

				inviteeName := input.AskString("What is the invitee name")
				inviteEmail := input.AskString("What is the invitee email")
				meetingTitle := input.AskString("What is the meeting title")

				createdMeeting, err := client.CreateMeeting(inviteeName, inviteEmail, meetingTitle, meetingTime, linkMap[linkName].Id)
				if err != nil {
					return err
				}
				logrus.Infof("meeting created: %s", createdMeeting.ConferenceData.JoinUrl)
				return nil
			},
		},
		{
			Name:        "archive",
			Description: "Archive complete tasks",
			Action: func(c *cli.Context) error {
				noOfWorkers := 10
				maxCount := c.Uint("max-count")
				if maxCount == 0 {
					maxCount = 10
				}
				autoArchiveAge := c.Uint("auto-archive-age")
				if autoArchiveAge == 0 {
					autoArchiveAge = 365
				}
				tasks, err := client.GetTasks([]string{"COMPLETE"})
				if err != nil {
					return err
				}

				logrus.Infof("found %d tasks. will attempt to archive %d of them", len(tasks), maxCount)
				wg := sync.WaitGroup{}
				taskToArchive := make(chan *reclaim.Task)
				output := make(chan error)

				for i := 0; i < noOfWorkers; i++ {
					wg.Add(1)
					go func() {
						defer wg.Done()
						for task := range taskToArchive {
							task.Status = "ARCHIVED"
							_, err := client.UpdateTask(task)
							output <- err
						}
					}()
				}

				go func() {
					defer close(taskToArchive)
					for _, task := range tasks {
						if maxCount > 0 {
							isOldTask := task.Finished.Before(time.Now().Add(time.Hour * 24 * time.Duration(autoArchiveAge) * -1))
							if isOldTask || input.AskForConfirmation(fmt.Sprintf("Would you like to archive '%s': %v?", task.Title, task.Finished)) {
								logrus.Infof("Archiving task %d", task.Id)
								taskToArchive <- task
								if !isOldTask {
									maxCount -= 1
								}

							}
						}
					}
				}()

				go func() {
					defer close(output)
					wg.Wait()
					logrus.Info("finished archiving tasks")
				}()

				foundErr := false
				for err := range output {
					if err != nil {
						foundErr = true
						logrus.Error(err)
					}
				}

				if foundErr {
					return errors.New("error when archiving tasks - see logs")
				}

				return nil
			},
			Flags: []cli.Flag{
				&cli.UintFlag{
					Name:        "max-count",
					DefaultText: "10",
					Required:    false,
				},
				&cli.UintFlag{
					Name:     "auto-archive-age",
					Required: false,
					Usage:    "no of days after which tasks should be auto archived",
				},
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

func removeGitlabTaskIfClosed(client *reclaim.Client, task *reclaim.Task) error {
	gitlabUrl := os.Getenv("GITLAB_URL")
	if gitlabUrl == "" {
		logrus.Warn("GITLAB_URL not set, skipping gitlab tasks")
		return nil
	}

	if strings.Contains(task.Title, gitlabUrl) && !strings.Contains(task.Title, "#") {
		pid := strings.Split(strings.Split(task.Title, "/-/")[0], "com/")[1]

		token := os.Getenv("GITLAB_TOKEN")
		if token == "" {
			logrus.Warn("GITLAB_TOKEN not set, skipping gitlab tasks")
			return nil
		}

		git, err := gitlab.NewClient(token, gitlab.WithBaseURL(gitlabUrl))
		if err != nil {
			log.Fatalf("Failed to create client: %v", err)
			return err
		}

		mr, _, err := git.MergeRequests.GetMergeRequest(pid, getLastSegmentAsInt(task.Title), nil)
		if err != nil {
			return err
		}
		if mr.State == "merged" {
			logrus.Infof("removing task: %s", task.Title)
			return client.DeleteTask(task.Id)
		}
		if mr.State != "opened" {
			return fmt.Errorf("MR state: %s", mr.State)
		}
	}
	return nil
}

func getLastSegmentAsInt(s string) int {
	lastSegment := getLastSegment(s)
	lastSegmentInt, err := strconv.Atoi(strings.TrimRightFunc(lastSegment, func(r rune) bool {
		return r < '0' || r > '9'
	}))
	if err != nil {
		log.Fatalf("Failed to convert last segment to int: %v", err)
	}
	return lastSegmentInt
}

func getLastSegment(s string) string {
	lastSlashIndex := strings.LastIndex(s, "/")
	if lastSlashIndex == -1 {
		return s
	}
	return s[lastSlashIndex+1:]
}
