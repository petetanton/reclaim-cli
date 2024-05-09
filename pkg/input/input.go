package input

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	sterm "github.com/AlecAivazis/survey/v2/terminal"
	"github.com/sirupsen/logrus"
	terminal "golang.org/x/term"
)

var tries = 0

type Config struct {
	required  bool
	secure    bool
	multiline bool
}

func isTerminal() bool {
	return terminal.IsTerminal(int(os.Stdout.Fd()))
}

func getAskOptions(options *survey.AskOptions) (err error) {
	// use stdout if not piping
	if isTerminal() {
		options.Stdio = sterm.Stdio{
			In:  os.Stdin,
			Out: os.Stdout,
			Err: os.Stderr,
		}
		return
	}

	options.Stdio = sterm.Stdio{
		In:  os.Stdin,
		Out: os.Stderr,
		Err: os.Stderr,
	}
	return
}

func AskStringWithOptions(question, def string, options Config) (string, error) {
	tries++

	message := fmt.Sprintf("%s: ", question)

	response := ""

	var prompt survey.Prompt
	switch {
	case options.secure:
		prompt = &survey.Password{
			Message: message,
		}
	case options.multiline:
		prompt = &survey.Multiline{
			Message: message,
			Default: def,
		}
	default:
		prompt = &survey.Input{
			Message: message,
			Default: def,
		}
	}

	if err := survey.AskOne(prompt, &response, nil, getAskOptions); err != nil {
		return "", err
	}

	if response == "" {
		if options.required {
			if tries > 1 {
				logrus.Fatal("Field is required and cannot be empty, exiting")
			}
			fmt.Println("Field is required and cannot be empty")
			return AskStringWithOptions(question, def, options)
		}

		response = def
	}

	tries = 0
	return strings.TrimSpace(response), nil
}

func AskMultiSelect(question string, options []string) []string {
	response, err := AskMultiSelectWithError(question, options)
	if err != nil {
		logrus.Fatal(err)
	}

	return response
}

func AskMultiSelectWithError(question string, options []string) ([]string, error) {
	response := []string{}

	prompt := &survey.MultiSelect{
		Message:  fmt.Sprintf("%s:", question),
		Options:  options,
		PageSize: 10,
	}

	err := survey.AskOne(prompt, &response, nil, getAskOptions)

	return response, err
}

func AskSelectMapKeys(question string, options map[string]string) string {
	var ss []string
	for s, _ := range options {
		ss = append(ss, s)
	}

	return AskSelect(question, ss)
}

func AskSelect(question string, options []string) string {
	time.Sleep(time.Millisecond * 250)
	response, err := AskSelectWithError(question, options)
	if err != nil {
		logrus.Fatal(err)
	}

	return response
}

func AskSelectWithFilter(question string, options []string, filter string) string {
	var newOptions []string

	for _, option := range options {
		if strings.Contains(option, filter) {
			newOptions = append(newOptions, option)
		}
	}
	return AskSelect(question, newOptions)
}

func AskSelectWithError(question string, options []string) (string, error) {
	response := ""
	prompt := &survey.Select{
		Message:  fmt.Sprintf("%s:", question),
		Options:  options,
		PageSize: 10,
	}
	err := survey.AskOne(prompt, &response, nil, getAskOptions)

	return response, err
}

func AskForConfirmation(question string) bool {
	response, err := AskForConfirmationWithError(question)
	if err != nil {
		logrus.Fatal(err)
	}

	return response
}

func AskForConfirmationWithError(question string) (bool, error) {
	response := false
	prompt := &survey.Confirm{
		Message: question,
	}
	err := survey.AskOne(prompt, &response, nil, getAskOptions)

	return response, err
}
