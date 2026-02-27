package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"marcel-cli/api"
	"marcel-cli/config"
	"marcel-cli/models"
)

type Storage struct {
	config    *config.Config
	apiClient *api.Client
}

type CacheData struct {
	Timestamp time.Time       `json:"timestamp"`
	Journeys  []models.Journey `json:"journeys"`
	Quests    []models.Quest   `json:"quests"`
	Habits    []models.Habit   `json:"habits"`
	Events    []models.Event   `json:"events"`
}

func New() (*Storage, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	client := api.NewClient(cfg)

	return &Storage{
		config:    cfg,
		apiClient: client,
	}, nil
}

func (s *Storage) Load() (*models.AppData, error) {
	return s.LoadAll()
}

func (s *Storage) Save(data *models.AppData) error {
	return nil
}

func (s *Storage) GetConfig() *config.Config {
	return s.config
}

func (s *Storage) GetAPIClient() *api.Client {
	return s.apiClient
}

func (s *Storage) getCachePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	cacheDir := filepath.Join(homeDir, ".marcel")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return "", err
	}

	return filepath.Join(cacheDir, "cache.json"), nil
}

func (s *Storage) LoadFromCache() (*models.AppData, error) {
	cachePath, err := s.getCachePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(cachePath)
	if err != nil {
		return nil, err
	}

	var cache CacheData
	if err := json.Unmarshal(data, &cache); err != nil {
		return nil, err
	}

	appData := models.NewAppData()

	questsByJourney := make(map[int][]models.Quest)
	var unassignedQuests []models.Quest

	for _, quest := range cache.Quests {
		if quest.JourneyID != nil {
			questsByJourney[*quest.JourneyID] = append(questsByJourney[*quest.JourneyID], quest)
		} else {
			unassignedQuests = append(unassignedQuests, quest)
		}
	}

	if len(unassignedQuests) > 0 {
		appData.Journeys = append(appData.Journeys, models.Journey{
			ID:     0,
			Name:   "My Quests",
			Quests: unassignedQuests,
		})
	}

	for _, journey := range cache.Journeys {
		journey.Quests = questsByJourney[journey.ID]
		if len(journey.Quests) > 0 || journey.ID != 0 {
			appData.Journeys = append(appData.Journeys, journey)
		}
	}

	appData.Habits = cache.Habits
	appData.Events = cache.Events
	appData.CurrentSection = "quests"

	return &appData, nil
}

func (s *Storage) SaveToCache(journeys []models.Journey, quests []models.Quest, habits []models.Habit, events []models.Event) error {
	cachePath, err := s.getCachePath()
	if err != nil {
		return err
	}

	cache := CacheData{
		Timestamp: time.Now(),
		Journeys:  journeys,
		Quests:    quests,
		Habits:    habits,
		Events:    events,
	}

	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(cachePath, data, 0644)
}

func (s *Storage) LoadWithCache() (*models.AppData, error) {
	cachedData, err := s.LoadFromCache()
	if err == nil {
		go func() {
			s.LoadAll()
		}()
		return cachedData, nil
	}

	return s.LoadQuestsOnly()
}

func (s *Storage) LoadQuestsOnly() (*models.AppData, error) {
	data := models.NewAppData()

	quests, err := s.apiClient.GetQuests()
	if err != nil {
		return &data, err
	}

	var unassignedQuests []models.Quest
	for _, quest := range quests {
		if quest.JourneyID == nil {
			unassignedQuests = append(unassignedQuests, quest)
		}
	}

	if len(unassignedQuests) > 0 {
		data.Journeys = append(data.Journeys, models.Journey{
			ID:     0,
			Name:   "My Quests",
			Quests: unassignedQuests,
		})
	}

	data.CurrentSection = "quests"

	return &data, nil
}

func (s *Storage) LoadAll() (*models.AppData, error) {
	data := models.NewAppData()

	quests, err := s.apiClient.GetQuests()
	if err != nil {
		return &data, err
	}

	journeys, err := s.apiClient.GetJourneys()
	if err != nil {
		return &data, err
	}

	habits, err := s.apiClient.GetHabits()
	if err != nil {
		return &data, err
	}

	events, err := s.apiClient.GetEvents()
	if err != nil {
		return &data, err
	}

	questsByJourney := make(map[int][]models.Quest)
	var unassignedQuests []models.Quest

	for _, quest := range quests {
		if quest.JourneyID != nil {
			questsByJourney[*quest.JourneyID] = append(questsByJourney[*quest.JourneyID], quest)
		} else {
			unassignedQuests = append(unassignedQuests, quest)
		}
	}

	if len(unassignedQuests) > 0 {
		data.Journeys = append(data.Journeys, models.Journey{
			ID:     0,
			Name:   "My Quests",
			Quests: unassignedQuests,
		})
	}

	for _, journey := range journeys {
		journey.Quests = questsByJourney[journey.ID]
		if len(journey.Quests) > 0 || journey.ID != 0 {
			data.Journeys = append(data.Journeys, journey)
		}
	}

	data.Habits = habits
	data.Events = events
	data.CurrentSection = "quests"

	s.SaveToCache(journeys, quests, habits, events)

	return &data, nil
}
