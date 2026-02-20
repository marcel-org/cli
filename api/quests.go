package api

import (
	"encoding/json"
	"fmt"
	"io"

	"marcel-cli/models"
)

type QuestsResponse struct {
	Quests []models.Quest `json:"quests"`
}

type QuestResponse struct {
	Quest models.Quest `json:"quest"`
}

type CreateQuestRequest struct {
	Title      string `json:"title"`
	Note       string `json:"note"`
	Difficulty string `json:"difficulty"`
	JourneyID  *int   `json:"journeyId,omitempty"`
}

type UpdateQuestRequest struct {
	Title      *string `json:"title,omitempty"`
	Note       *string `json:"note,omitempty"`
	Done       *bool   `json:"done,omitempty"`
	Difficulty *string `json:"difficulty,omitempty"`
}

func (c *Client) GetQuests() ([]models.Quest, error) {
	resp, err := c.doRequest("GET", "/quest", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get quests: status %d, body: %s", resp.StatusCode, string(body))
	}

	var result QuestsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Quests, nil
}

func (c *Client) CreateQuest(title, note, difficulty string, journeyID *int) (*models.Quest, error) {
	req := CreateQuestRequest{
		Title:      title,
		Note:       note,
		Difficulty: difficulty,
		JourneyID:  journeyID,
	}

	resp, err := c.doRequest("POST", "/quest", req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to create quest: status %d, body: %s", resp.StatusCode, string(body))
	}

	var result QuestResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result.Quest, nil
}

func (c *Client) UpdateQuest(questID int, updates UpdateQuestRequest) (*models.Quest, error) {
	path := fmt.Sprintf("/quest/%d", questID)
	resp, err := c.doRequest("PUT", path, updates)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to update quest: status %d, body: %s", resp.StatusCode, string(body))
	}

	var result QuestResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result.Quest, nil
}

func (c *Client) ToggleQuest(questID int, done bool) (*models.Quest, error) {
	return c.UpdateQuest(questID, UpdateQuestRequest{
		Done: &done,
	})
}

func (c *Client) DeleteQuest(questID int) error {
	path := fmt.Sprintf("/quest/%d", questID)
	resp, err := c.doRequest("DELETE", path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete quest: status %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}
