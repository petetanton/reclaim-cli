package reclaim

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	apiUrl = "https://api.app.reclaim.ai"
)

type Client struct {
	h      http.Client
	apiKey string
}

func New() *Client {
	return &Client{h: http.Client{
		Timeout: time.Second * 10,
	},
		apiKey: os.Getenv("RECLAIM_API_KEY")}
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	return c.h.Do(req)
}

func (c *Client) CreateTask(title string, minChunkSize int, maxChunkSize int, timeChunksRequired int) (*Task, error) {
	requestBody := fmt.Sprintf(`{
		"title": "%s",
		"status": "NEW",
		"minChunkSize": %d,
		"maxChunkSize": %d,
		"timeChunksRequired": %d,
		"eventCategory": "WORK"
}`, title, minChunkSize, maxChunkSize, timeChunksRequired)
	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/tasks", apiUrl), strings.NewReader(requestBody))

	if err != nil {
		return nil, err
	}

	response, err := c.do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d when creating a task:\n%s", response.StatusCode, requestBody)
	}

	var task *Task
	reqBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(reqBytes, &task)

	if err != nil {
		return nil, err
	}

	return task, nil
}

func (c *Client) SnoozeTask(taskId int, snoozeUntil time.Time) error {
	requestBody := fmt.Sprintf(`{
		"snoozeUntil": "%s"
}`, snoozeUntil.Format(time.RFC3339Nano))
	request, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/api/tasks/%d", apiUrl, taskId), strings.NewReader(requestBody))

	if err != nil {
		return err
	}

	response, err := c.do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code %d when snoozing a task:\n%s", response.StatusCode, requestBody)
	}

	return nil
}

func (c *Client) GetTasks() ([]*Task, error) {
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/tasks?status=NEW,SCHEDULED,IN_PROGRESS,COMPLETE", apiUrl), nil)
	if err != nil {
		return nil, err
	}

	response, err := c.do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d when getting tasks", response.StatusCode)
	}
	var tasks []*Task
	reqBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(reqBytes, &tasks)

	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (c *Client) DeleteTask(taskId int) error {
	request, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/tasks/%d", apiUrl, taskId), nil)
	if err != nil {
		return err
	}

	response, err := c.do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code %d when deleting a task", response.StatusCode)
	}

	return nil
}

func (c *Client) UpdateTask(task *Task) (*Task, error) {
	requestBody, err := json.Marshal(task)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/api/tasks/%d", apiUrl, task.Id), strings.NewReader(string(requestBody)))
	if err != nil {
		return nil, err
	}

	response, err := c.do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var updatedTask *Task
	reqBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(reqBytes, &updatedTask)

	if err != nil {
		return nil, err
	}

	return updatedTask, nil
}
