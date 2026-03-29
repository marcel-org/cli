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

type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	Level        int    `json:"level"`
	Shape        string `json:"shape"`
	Color        string `json:"color"`
	HatID        *int   `json:"hatId"`
	GlassesID    *int   `json:"glassesId"`
	HandID       *int   `json:"handId"`
	BackgroundID *int   `json:"backgroundId"`
	XP           int    `json:"xp"`
	XPMax        int    `json:"xpMax"`
}

type SpaceMember struct {
	ID         int       `json:"id"`
	SpaceID    int       `json:"spaceId"`
	UserID     int       `json:"userId"`
	Permission string    `json:"permission"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	User       User      `json:"user"`
}

type Space struct {
	ID        int           `json:"id"`
	Name      string        `json:"name"`
	OwnerID   int           `json:"ownerId"`
	CreatedAt time.Time     `json:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt"`
	Owner     User          `json:"owner"`
	Members   []SpaceMember `json:"members"`
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
	ID               int        `json:"id"`
	Name             string     `json:"name"`
	AuthorID         int        `json:"authorId"`
	XPReward         int        `json:"xpReward"`
	GoldReward       int        `json:"goldReward"`
	CycleType        string     `json:"cycleType"`
	CycleConfig      any        `json:"cycleConfig"`
	Completed        []string   `json:"completed"`
	CurrentStreak    int        `json:"currentStreak"`
	MaxStreak        int        `json:"maxStreak"`
	StartDate        time.Time  `json:"startDate"`
	EndDate          *time.Time `json:"endDate"`
	CreatedAt        time.Time  `json:"createdAt"`
	UpdatedAt        time.Time  `json:"updatedAt"`
	CycleDescription string     `json:"cycleDescription"`
	IsDueToday       bool       `json:"isDueToday"`
}

type Event struct {
	ID               int        `json:"id"`
	Title            string     `json:"title"`
	Date             time.Time  `json:"date"`
	EndDate          *time.Time `json:"endDate"`
	Time             *string    `json:"time"`
	EndTime          *string    `json:"endTime"`
	Location         *string    `json:"location"`
	Description      *string    `json:"description"`
	AuthorID         int        `json:"authorId"`
	GoogleCalendarID *string    `json:"googleCalendarId"`
	CreatedAt        time.Time  `json:"createdAt"`
	UpdatedAt        time.Time  `json:"updatedAt"`
}

type AppData struct {
	Journeys       []Journey
	Habits         []Habit
	Events         []Event
	Spaces         []Space
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
		Habits:         []Habit{},
		Events:         []Event{},
		Spaces:         []Space{},
		CurrentJourney: 0,
	}
}
