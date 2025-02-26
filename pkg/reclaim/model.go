package reclaim

import "time"

type Task struct {
	Id                  int     `json:"id"`
	Title               string  `json:"title"`
	Notes               string  `json:"notes"`
	EventCategory       string  `json:"eventCategory"`
	EventSubType        string  `json:"eventSubType"`
	Status              string  `json:"status"`
	TimeChunksRequired  int     `json:"timeChunksRequired"`
	TimeChunksSpent     int     `json:"timeChunksSpent"`
	TimeChunksRemaining int     `json:"timeChunksRemaining"`
	MinChunkSize        int     `json:"minChunkSize"`
	MaxChunkSize        int     `json:"maxChunkSize"`
	AlwaysPrivate       bool    `json:"alwaysPrivate"`
	Deleted             bool    `json:"deleted"`
	Index               float64 `json:"index"`
	//Due                 time.Time `json:"due"`  // this field is a strange format e.g. 0000-12-31T23:58:45-00:01:15
	Created      time.Time `json:"created"`
	Updated      time.Time `json:"updated"`
	Finished     time.Time `json:"finished"`
	Adjusted     bool      `json:"adjusted"`
	AtRisk       bool      `json:"atRisk"`
	TimeSchemeId string    `json:"timeSchemeId"`
	Priority     string    `json:"priority"`
	OnDeck       bool      `json:"onDeck"`
	Deferred     bool      `json:"deferred"`
	SortKey      float64   `json:"sortKey"`
	TaskSource   struct {
		Type string `json:"type"`
	} `json:"taskSource"`
	ReadOnlyFields          []interface{} `json:"readOnlyFields"`
	RecurringAssignmentType string        `json:"recurringAssignmentType"`
	Type                    string        `json:"type"`
}

type MeetingResponse struct {
	MeetingId string `json:"meetingId"`
	Event     struct {
		Key              string      `json:"key"`
		Public           bool        `json:"public"`
		Locked           interface{} `json:"locked"`
		Duration         float64     `json:"duration"`
		Version          string      `json:"version"`
		StartTime        time.Time   `json:"startTime"`
		Status           string      `json:"status"`
		EventId          string      `json:"eventId"`
		EndTime          time.Time   `json:"endTime"`
		ReclaimManaged   bool        `json:"reclaimManaged"`
		CalendarId       int         `json:"calendarId"`
		Etag             string      `json:"etag"`
		UserId           interface{} `json:"userId"`
		Title            string      `json:"title"`
		OnlineMeetingUrl string      `json:"onlineMeetingUrl"`
		PersonalSync     bool        `json:"personalSync"`
		Attendees        struct {
			PtantonIndeedCom string `json:"ptanton@indeed.com"`
		} `json:"attendees"`
		Color             string      `json:"color"`
		Category          string      `json:"category"`
		RsvpStatus        string      `json:"rsvpStatus"`
		Free              bool        `json:"free"`
		Rrule             interface{} `json:"rrule"`
		ConferenceDetails struct {
			Solution string `json:"solution"`
			Url      string `json:"url"`
			Status   string `json:"status"`
			Source   string `json:"source"`
		} `json:"conferenceDetails"`
		DisplayColorHex  string `json:"displayColorHex"`
		ReclaimEventType string `json:"reclaimEventType"`
		EventDescription struct {
			Raw       string `json:"raw"`
			Processed string `json:"processed"`
		} `json:"eventDescription"`
	} `json:"event"`
	SchedulingLinkId string `json:"schedulingLinkId"`
	Organizer        struct {
		UserId    string `json:"userId"`
		Email     string `json:"email"`
		Name      string `json:"name"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		AvatarUrl string `json:"avatarUrl"`
	} `json:"organizer"`
	Attendee struct {
		UserId         string `json:"userId"`
		Email          string `json:"email"`
		Name           string `json:"name"`
		FirstName      string `json:"firstName"`
		LastName       string `json:"lastName"`
		AvatarUrl      string `json:"avatarUrl"`
		AttendanceType string `json:"attendanceType"`
	} `json:"attendee"`
	AttendeeZoneId struct {
		Id           string `json:"id"`
		DisplayName  string `json:"displayName"`
		Abbreviation string `json:"abbreviation"`
	} `json:"attendeeZoneId"`
	TargetCalendarId int `json:"targetCalendarId"`
	ConferenceData   struct {
		JoinUrl     string    `json:"join_url"`
		Timezone    string    `json:"timezone"`
		CreatedAt   time.Time `json:"created_at"`
		Type        int       `json:"type"`
		Uuid        string    `json:"uuid"`
		HostId      string    `json:"host_id"`
		Duration    int       `json:"duration"`
		StartTime   time.Time `json:"start_time"`
		Password    string    `json:"password"`
		WaitingRoom bool      `json:"waiting_room"`
		Topic       string    `json:"topic"`
		Id          int64     `json:"id"`
		HostEmail   string    `json:"host_email"`
		Status      string    `json:"status"`
	} `json:"conferenceData"`
	UserSlug struct {
		Slug string `json:"slug"`
	} `json:"userSlug"`
	SchedulingLink struct {
		Id              string `json:"id"`
		Title           string `json:"title"`
		Slug            string `json:"slug"`
		PageSlug        string `json:"pageSlug"`
		Description     string `json:"description"`
		Enabled         bool   `json:"enabled"`
		Hidden          bool   `json:"hidden"`
		MainOrganizerId string `json:"mainOrganizerId"`
		HostId          string `json:"hostId"`
		Owner           struct {
			UserId    string `json:"userId"`
			Email     string `json:"email"`
			Name      string `json:"name"`
			FirstName string `json:"firstName"`
			LastName  string `json:"lastName"`
			AvatarUrl string `json:"avatarUrl"`
		} `json:"owner"`
		Organizers []struct {
			Organizer struct {
				UserId    string `json:"userId"`
				Email     string `json:"email"`
				Name      string `json:"name"`
				FirstName string `json:"firstName"`
				LastName  string `json:"lastName"`
				AvatarUrl string `json:"avatarUrl"`
			} `json:"organizer"`
			Role     string `json:"role"`
			Timezone struct {
				Id           string `json:"id"`
				DisplayName  string `json:"displayName"`
				Abbreviation string `json:"abbreviation"`
			} `json:"timezone"`
			TimePolicyType     string `json:"timePolicyType"`
			TimeSchemeId       string `json:"timeSchemeId"`
			ResolvedTimePolicy struct {
				DayHours struct {
					FRIDAY struct {
						Intervals []struct {
							Start    string  `json:"start"`
							End      string  `json:"end"`
							Duration float64 `json:"duration"`
						} `json:"intervals"`
						EndOfDay   string `json:"endOfDay"`
						StartOfDay string `json:"startOfDay"`
					} `json:"FRIDAY"`
					MONDAY struct {
						Intervals []struct {
							Start    string  `json:"start"`
							End      string  `json:"end"`
							Duration float64 `json:"duration"`
						} `json:"intervals"`
						EndOfDay   string `json:"endOfDay"`
						StartOfDay string `json:"startOfDay"`
					} `json:"MONDAY"`
					TUESDAY struct {
						Intervals []struct {
							Start    string  `json:"start"`
							End      string  `json:"end"`
							Duration float64 `json:"duration"`
						} `json:"intervals"`
						EndOfDay   string `json:"endOfDay"`
						StartOfDay string `json:"startOfDay"`
					} `json:"TUESDAY"`
					THURSDAY struct {
						Intervals []struct {
							Start    string  `json:"start"`
							End      string  `json:"end"`
							Duration float64 `json:"duration"`
						} `json:"intervals"`
						EndOfDay   string `json:"endOfDay"`
						StartOfDay string `json:"startOfDay"`
					} `json:"THURSDAY"`
					WEDNESDAY struct {
						Intervals []struct {
							Start    string  `json:"start"`
							End      string  `json:"end"`
							Duration float64 `json:"duration"`
						} `json:"intervals"`
						EndOfDay   string `json:"endOfDay"`
						StartOfDay string `json:"startOfDay"`
					} `json:"WEDNESDAY"`
				} `json:"dayHours"`
				StartOfWeek string `json:"startOfWeek"`
				EndOfWeek   string `json:"endOfWeek"`
			} `json:"resolvedTimePolicy"`
			Status               string   `json:"status"`
			Optional             bool     `json:"optional"`
			AttendanceType       string   `json:"attendanceType"`
			ValidConferenceTypes []string `json:"validConferenceTypes"`
			TargetCalendarId     int      `json:"targetCalendarId"`
		} `json:"organizers"`
		EffectiveTimePolicy struct {
			DayHours struct {
				FRIDAY struct {
					Intervals []struct {
						Start    string  `json:"start"`
						End      string  `json:"end"`
						Duration float64 `json:"duration"`
					} `json:"intervals"`
					EndOfDay   string `json:"endOfDay"`
					StartOfDay string `json:"startOfDay"`
				} `json:"FRIDAY"`
				THURSDAY struct {
					Intervals []struct {
						Start    string  `json:"start"`
						End      string  `json:"end"`
						Duration float64 `json:"duration"`
					} `json:"intervals"`
					EndOfDay   string `json:"endOfDay"`
					StartOfDay string `json:"startOfDay"`
				} `json:"THURSDAY"`
				MONDAY struct {
					Intervals []struct {
						Start    string  `json:"start"`
						End      string  `json:"end"`
						Duration float64 `json:"duration"`
					} `json:"intervals"`
					EndOfDay   string `json:"endOfDay"`
					StartOfDay string `json:"startOfDay"`
				} `json:"MONDAY"`
				TUESDAY struct {
					Intervals []struct {
						Start    string  `json:"start"`
						End      string  `json:"end"`
						Duration float64 `json:"duration"`
					} `json:"intervals"`
					EndOfDay   string `json:"endOfDay"`
					StartOfDay string `json:"startOfDay"`
				} `json:"TUESDAY"`
				WEDNESDAY struct {
					Intervals []struct {
						Start    string  `json:"start"`
						End      string  `json:"end"`
						Duration float64 `json:"duration"`
					} `json:"intervals"`
					EndOfDay   string `json:"endOfDay"`
					StartOfDay string `json:"startOfDay"`
				} `json:"WEDNESDAY"`
			} `json:"dayHours"`
			StartOfWeek string `json:"startOfWeek"`
			EndOfWeek   string `json:"endOfWeek"`
		} `json:"effectiveTimePolicy"`
		Durations       []int  `json:"durations"`
		DefaultDuration int    `json:"defaultDuration"`
		DelayStart      string `json:"delayStart"`
		DelayStartUnits int    `json:"delayStartUnits"`
		DaysIntoFuture  int    `json:"daysIntoFuture"`
		Priority        string `json:"priority"`
		LocationOptions []struct {
			ConferenceType string `json:"conferenceType"`
		} `json:"locationOptions"`
		DefaultLocationIndex int           `json:"defaultLocationIndex"`
		IconType             string        `json:"iconType"`
		OrganizerRefCode     string        `json:"organizerRefCode"`
		OwnerRefCode         string        `json:"ownerRefCode"`
		MeetingTitle         string        `json:"meetingTitle"`
		LinkGroupId          string        `json:"linkGroupId"`
		LinkGroupName        string        `json:"linkGroupName"`
		TargetCalendarId     int           `json:"targetCalendarId"`
		DisableBuffers       bool          `json:"disableBuffers"`
		SharedMeetingTimes   []int         `json:"sharedMeetingTimes"`
		ResolvedBrandingMode string        `json:"resolvedBrandingMode"`
		OptionalOrganizer    bool          `json:"optionalOrganizer"`
		OwnerAttendanceType  string        `json:"ownerAttendanceType"`
		FixedTimePolicy      bool          `json:"fixedTimePolicy"`
		WebhookConfigIds     []interface{} `json:"webhookConfigIds"`
		FourOhFourAlerts     []interface{} `json:"fourOhFourAlerts"`
		Permissions          struct {
			CanView   bool `json:"canView"`
			CanEdit   bool `json:"canEdit"`
			CanEnable bool `json:"canEnable"`
			CanDelete bool `json:"canDelete"`
		} `json:"permissions"`
	} `json:"schedulingLink"`
	Message         string `json:"message"`
	MeetingLocation struct {
		ConferenceType string `json:"conferenceType"`
	} `json:"meetingLocation"`
	Ccs            []interface{} `json:"ccs"`
	SurveyResponse struct {
		Responses []interface{} `json:"responses"`
	} `json:"surveyResponse"`
	Organizers []struct {
		Organizer struct {
			UserId         string `json:"userId"`
			Email          string `json:"email"`
			Name           string `json:"name"`
			FirstName      string `json:"firstName"`
			LastName       string `json:"lastName"`
			AvatarUrl      string `json:"avatarUrl"`
			AttendanceType string `json:"attendanceType"`
		} `json:"organizer"`
		IsHost bool `json:"isHost"`
	} `json:"organizers"`
}

type MeetingRequest struct {
	InviteeName     string   `json:"inviteeName"`
	CcEmails        []string `json:"ccEmails,omitempty"`
	Message         string   `json:"message"`
	MeetingLocation struct {
		ConferenceType string `json:"conferenceType"`
	} `json:"meetingLocation"`
	AttendeeTimeZone string    `json:"attendeeTimeZone"`
	Start            time.Time `json:"start"`
	End              time.Time `json:"end"`
	InviteeEmail     string    `json:"inviteeEmail"`
	CustomData       struct {
	} `json:"customData,omitempty"`
	InviteeZoneId string `json:"inviteeZoneId"`
}

type MeetingTime struct {
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"endTime"`
	IsSuggested bool      `json:"isSuggested"`
}

type MeetingTimeResponse struct {
	AvailableTimes struct {
		ThirtyMinuteSlots []*MeetingTime `json:"30"`
	} `json:"availableTimes"`
}

type ScheduleLink struct {
	Id              string `json:"id"`
	Title           string `json:"title"`
	Slug            string `json:"slug"`
	PageSlug        string `json:"pageSlug"`
	Description     string `json:"description"`
	Enabled         bool   `json:"enabled"`
	Hidden          bool   `json:"hidden"`
	MainOrganizerId string `json:"mainOrganizerId"`
	HostId          string `json:"hostId"`
	Owner           struct {
		UserId    string `json:"userId"`
		Email     string `json:"email"`
		Name      string `json:"name"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		AvatarUrl string `json:"avatarUrl"`
	} `json:"owner"`
	LinkOwnerTeamId int `json:"linkOwnerTeamId"`
	Organizers      []struct {
		Organizer struct {
			UserId    string `json:"userId"`
			Email     string `json:"email"`
			Name      string `json:"name"`
			FirstName string `json:"firstName"`
			LastName  string `json:"lastName"`
			AvatarUrl string `json:"avatarUrl"`
		} `json:"organizer"`
		Role     string `json:"role"`
		Timezone struct {
			Id           string `json:"id"`
			DisplayName  string `json:"displayName"`
			Abbreviation string `json:"abbreviation"`
		} `json:"timezone"`
		TimePolicyType     string `json:"timePolicyType"`
		TimeSchemeId       string `json:"timeSchemeId"`
		ResolvedTimePolicy struct {
			DayHours struct {
				FRIDAY struct {
					Intervals []struct {
						Start    string  `json:"start"`
						End      string  `json:"end"`
						Duration float64 `json:"duration"`
					} `json:"intervals"`
					StartOfDay string `json:"startOfDay"`
					EndOfDay   string `json:"endOfDay"`
				} `json:"FRIDAY"`
				MONDAY struct {
					Intervals []struct {
						Start    string  `json:"start"`
						End      string  `json:"end"`
						Duration float64 `json:"duration"`
					} `json:"intervals"`
					StartOfDay string `json:"startOfDay"`
					EndOfDay   string `json:"endOfDay"`
				} `json:"MONDAY"`
				TUESDAY struct {
					Intervals []struct {
						Start    string  `json:"start"`
						End      string  `json:"end"`
						Duration float64 `json:"duration"`
					} `json:"intervals"`
					StartOfDay string `json:"startOfDay"`
					EndOfDay   string `json:"endOfDay"`
				} `json:"TUESDAY"`
				THURSDAY struct {
					Intervals []struct {
						Start    string  `json:"start"`
						End      string  `json:"end"`
						Duration float64 `json:"duration"`
					} `json:"intervals"`
					StartOfDay string `json:"startOfDay"`
					EndOfDay   string `json:"endOfDay"`
				} `json:"THURSDAY"`
				WEDNESDAY struct {
					Intervals []struct {
						Start    string  `json:"start"`
						End      string  `json:"end"`
						Duration float64 `json:"duration"`
					} `json:"intervals"`
					StartOfDay string `json:"startOfDay"`
					EndOfDay   string `json:"endOfDay"`
				} `json:"WEDNESDAY"`
			} `json:"dayHours"`
			StartOfWeek string `json:"startOfWeek"`
			EndOfWeek   string `json:"endOfWeek"`
		} `json:"resolvedTimePolicy"`
		Status               string   `json:"status"`
		Optional             bool     `json:"optional"`
		AttendanceType       string   `json:"attendanceType"`
		ValidConferenceTypes []string `json:"validConferenceTypes"`
		TargetCalendarId     int      `json:"targetCalendarId"`
	} `json:"organizers"`
	EffectiveTimePolicy struct {
		DayHours struct {
			WEDNESDAY struct {
				Intervals []struct {
					Start    string  `json:"start"`
					End      string  `json:"end"`
					Duration float64 `json:"duration"`
				} `json:"intervals"`
				StartOfDay string `json:"startOfDay"`
				EndOfDay   string `json:"endOfDay"`
			} `json:"WEDNESDAY"`
			TUESDAY struct {
				Intervals []struct {
					Start    string  `json:"start"`
					End      string  `json:"end"`
					Duration float64 `json:"duration"`
				} `json:"intervals"`
				StartOfDay string `json:"startOfDay"`
				EndOfDay   string `json:"endOfDay"`
			} `json:"TUESDAY"`
			MONDAY struct {
				Intervals []struct {
					Start    string  `json:"start"`
					End      string  `json:"end"`
					Duration float64 `json:"duration"`
				} `json:"intervals"`
				StartOfDay string `json:"startOfDay"`
				EndOfDay   string `json:"endOfDay"`
			} `json:"MONDAY"`
			THURSDAY struct {
				Intervals []struct {
					Start    string  `json:"start"`
					End      string  `json:"end"`
					Duration float64 `json:"duration"`
				} `json:"intervals"`
				StartOfDay string `json:"startOfDay"`
				EndOfDay   string `json:"endOfDay"`
			} `json:"THURSDAY"`
			FRIDAY struct {
				Intervals []struct {
					Start    string  `json:"start"`
					End      string  `json:"end"`
					Duration float64 `json:"duration"`
				} `json:"intervals"`
				StartOfDay string `json:"startOfDay"`
				EndOfDay   string `json:"endOfDay"`
			} `json:"FRIDAY"`
		} `json:"dayHours"`
		StartOfWeek string `json:"startOfWeek"`
		EndOfWeek   string `json:"endOfWeek"`
	} `json:"effectiveTimePolicy"`
	Durations       []int  `json:"durations"`
	DefaultDuration int    `json:"defaultDuration"`
	DelayStart      string `json:"delayStart"`
	DelayStartUnits int    `json:"delayStartUnits"`
	DaysIntoFuture  int    `json:"daysIntoFuture"`
	Priority        string `json:"priority"`
	LocationOptions []struct {
		ConferenceType string `json:"conferenceType"`
	} `json:"locationOptions"`
	DefaultLocationIndex int           `json:"defaultLocationIndex"`
	IconType             string        `json:"iconType"`
	OrganizerRefCode     string        `json:"organizerRefCode"`
	OwnerRefCode         string        `json:"ownerRefCode"`
	MeetingTitle         string        `json:"meetingTitle"`
	LinkGroupId          string        `json:"linkGroupId"`
	LinkGroupName        string        `json:"linkGroupName"`
	TargetCalendarId     int           `json:"targetCalendarId"`
	DisableBuffers       bool          `json:"disableBuffers"`
	SharedMeetingTimes   []int         `json:"sharedMeetingTimes"`
	ResolvedBrandingMode string        `json:"resolvedBrandingMode"`
	OptionalOrganizer    bool          `json:"optionalOrganizer"`
	OwnerAttendanceType  string        `json:"ownerAttendanceType"`
	FixedTimePolicy      bool          `json:"fixedTimePolicy"`
	WebhookConfigIds     []interface{} `json:"webhookConfigIds"`
	FourOhFourAlerts     []interface{} `json:"fourOhFourAlerts"`
	Permissions          struct {
		CanView   bool `json:"canView"`
		CanEdit   bool `json:"canEdit"`
		CanEnable bool `json:"canEnable"`
		CanDelete bool `json:"canDelete"`
	} `json:"permissions"`
}
