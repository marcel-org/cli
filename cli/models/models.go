package models

import (
	"time"
)

type Quest struct {
	ID         int       `json:"id"`
	Title      string    `json:"title"`
	Note       string    `json:"note"`
	Done       bool      `json:"done"`
	Difficulty string    `json:"difficulty"`
	AuthorID   int       `json:"authorId"`
	XPReward   int       `json:"xpReward"`
	GoldReward int       `json:"goldReward"`
	Date       *string   `json:"date"`
	Time       *string   `json:"time"`
	JourneyID  *int      `json:"journeyId"`
	SpaceID    *int      `json:"spaceId"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type Journey struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	AuthorID  int       `json:"authorId"`
	SpaceID   *int      `json:"spaceId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Quests    []Quest   `json:"quests,omitempty"`
}

type Habit struct {
	ID              int       `json:"id"`
	Name            string    `json:"name"`
	AuthorID        int       `json:"authorId"`
	XPReward        int       `json:"xpReward"`
	GoldReward      int       `json:"goldReward"`
	CycleType       string    `json:"cycleType"`
	CycleConfig     any       `json:"cycleConfig"`
	Completed       []string  `json:"completed"`
	CurrentStreak   int       `json:"currentStreak"`
	MaxStreak       int       `json:"maxStreak"`
	StartDate       time.Time `json:"startDate"`
	EndDate         *time.Time `json:"endDate"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
	CycleDescription string   `json:"cycleDescription"`
	IsDueToday      bool      `json:"isDueToday"`
}

type Event struct {
	ID               int       `json:"id"`
	Title            string    `json:"title"`
	Date             time.Time `json:"date"`
	EndDate          *time.Time `json:"endDate"`
	Time             *string   `json:"time"`
	EndTime          *string   `json:"endTime"`
	Location         *string   `json:"location"`
	Description      *string   `json:"description"`
	AuthorID         int       `json:"authorId"`
	GoogleCalendarID *string   `json:"googleCalendarId"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

type AppData struct {
	Journeys       []Journey
	Habits         []Habit
	Events         []Event
	CurrentJourney int
	CurrentSection string
}

func NewQuest(title string) Quest {
	return Quest{
		Title:     title,
		Done:      false,
		CreatedAt: time.Now(),
	}
}

func NewJourney(name string) Journey {
	return Journey{
		Name:   name,
		Quests: []Quest{},
	}
}

func NewAppData() AppData {
	return AppData{
		Journeys:       []Journey{},
		CurrentJourney: 0,
	}
}
