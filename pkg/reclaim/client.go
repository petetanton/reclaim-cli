package reclaim

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	apiUrl = "https://api.app.reclaim.ai"
)

type TaskPriority string

const (
	P1 TaskPriority = "P1"
	P2 TaskPriority = "P2"
	P3 TaskPriority = "P3"
	P4 TaskPriority = "P4"
)

type Client struct {
	h http.Client
	// baseUrl is a field rather than the apiUrl constant so tests can point the
	// client at a stub server.
	baseUrl string
	apiKey  string
}

func New() *Client {
	return &Client{h: http.Client{
		Timeout: time.Second * 10,
	},
		baseUrl: apiUrl,
		apiKey:  os.Getenv("RECLAIM_API_KEY")}
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	return c.h.Do(req)
}

func (c *Client) CreateTask(task *CreateTaskRequest) (*Task, error) {
	if task.Status == "" {
		task.Status = "NEW"
	}
	if task.EventCategory == "" {
		task.EventCategory = "WORK"
	}

	requestBody, err := json.Marshal(task)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/tasks", c.baseUrl), strings.NewReader(string(requestBody)))
	if err != nil {
		return nil, err
	}

	response, err := c.do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	responseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d when creating a task:\nrequest: %s\nresponse: %s", response.StatusCode, requestBody, responseBytes)
	}

	var createdTask *Task
	if err := json.Unmarshal(responseBytes, &createdTask); err != nil {
		return nil, err
	}

	return createdTask, nil
}

func (c *Client) SnoozeTask(taskId int, snoozeUntil time.Time) (*Task, error) {
	requestBody, err := json.Marshal(&snoozeTaskRequest{SnoozeUntil: snoozeUntil})
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/api/tasks/%d", c.baseUrl, taskId), strings.NewReader(string(requestBody)))
	if err != nil {
		return nil, err
	}

	response, err := c.do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	responseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d when snoozing a task:\nrequest: %s\nresponse: %s", response.StatusCode, requestBody, responseBytes)
	}

	var snoozedTask *Task
	if err := json.Unmarshal(responseBytes, &snoozedTask); err != nil {
		return nil, err
	}

	return snoozedTask, nil
}

func (c *Client) GetTasks(statuses []string) ([]*Task, error) {
	if len(statuses) == 0 {
		statuses = []string{"NEW", "SCHEDULED", "IN_PROGRESS", "COMPLETE"}
	}
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/tasks?status=%s", c.baseUrl, strings.Join(statuses, ",")), nil)
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
	request, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/tasks/%d", c.baseUrl, taskId), nil)
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

	request, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/api/tasks/%d", c.baseUrl, task.Id), strings.NewReader(string(requestBody)))
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

func (c *Client) GetNextMeetingTime(linkId string) (*MeetingTime, error) {
	now := time.Now()
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/scheduling-link/%s/meeting/availability/V2?date=%s&zoneId=Europe/London&conferenceType=ZOOM", c.baseUrl, linkId, now.Format("2006-01-02")), nil)
	if err != nil {
		return nil, err
	}

	response, err := c.do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d when getting next meeting time", response.StatusCode)
	}

	var mtr *MeetingTimeResponse
	reqBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(reqBytes, &mtr)
	if err != nil {
		return nil, err
	}

	return mtr.AvailableTimes.ThirtyMinuteSlots[0], nil
}

func (c *Client) GetScheduleLinks() ([]*ScheduleLink, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/scheduling-link", c.baseUrl), nil)
	if err != nil {
		return nil, err
	}

	response, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d when getting scheduling links", response.StatusCode)
	}
	var scheduleLinks []*ScheduleLink
	reqBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(reqBytes, &scheduleLinks)
	if err != nil {
		return nil, err
	}

	return scheduleLinks, nil
}

func (c *Client) CreateMeeting(inviteeName string, inviteeEmail string, title string, meetingTime *MeetingTime, linkId string) (*MeetingResponse, error) {
	request := &MeetingRequest{
		InviteeName: inviteeName,
		Message:     title,
		MeetingLocation: struct {
			ConferenceType string `json:"conferenceType"`
		}{ConferenceType: "ZOOM"},
		AttendeeTimeZone: "Europe/London",
		Start:            meetingTime.StartTime,
		End:              meetingTime.EndTime,
		InviteeEmail:     inviteeEmail,
		InviteeZoneId:    "Europe/London",
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	logrus.Info(string(requestBody))

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/scheduling-link/%s/meeting", c.baseUrl, linkId), strings.NewReader(string(requestBody)))
	if err != nil {
		return nil, err
	}

	response, err := c.do(req)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	responseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d when creating a meeting: %s", response.StatusCode, string(responseBytes))
	}

	var meetingResponse *MeetingResponse
	err = json.Unmarshal(responseBytes, &meetingResponse)
	if err != nil {
		return nil, err
	}

	return meetingResponse, nil
}
