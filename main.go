package main

import (
	"github.com/petetanton/reclaim-cli/pkg"
	"github.com/petetanton/reclaim-cli/pkg/input"
	"github.com/petetanton/reclaim-cli/pkg/reclaim"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"os"
	"strconv"
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
				//if title == "" {
				//	return errors.New("title is required")
				//}
				mins := input.AskSelect("how many mins for the task", []string{"15", "30", "45", "60", "75", "90", "105", "120"})
				minsInt, err := strconv.Atoi(mins)
				if err != nil {
					return err
				}
				err = client.CreateTask(title, 1, 8, minsInt/15)
				if err != nil {
					return err
				}
				logrus.Printf("task %s created", title)
				return nil
			},
		},
		{
			Name: "version",
			Action: func(c *cli.Context) error {
				logrus.Print(pkg.Version)
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		logrus.Fatal(err)
	}
}
