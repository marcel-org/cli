package storage

import (
	"marcel-cli/api"
	"marcel-cli/config"
	"marcel-cli/models"
)

type Storage struct {
	config    *config.Config
	apiClient *api.Client
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
	data := models.NewAppData()

	journeys, err := s.apiClient.GetJourneys()
	if err != nil {
		return &data, err
	}

	quests, err := s.apiClient.GetQuests()
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

	return &data, nil
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
