package api

import (
	"encoding/json"
	"fmt"
	"io"

	"marcel-cli/models"
)

type JourneysResponse struct {
	Journeys []models.Journey `json:"journeys"`
}

type JourneyResponse struct {
	Journey models.Journey `json:"journey"`
}

type CreateJourneyRequest struct {
	Name string `json:"name"`
}

type UpdateJourneyRequest struct {
	Name *string `json:"name,omitempty"`
}

func (c *Client) GetJourneys() ([]models.Journey, error) {
	resp, err := c.doRequest("GET", "/journey", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get journeys: status %d, body: %s", resp.StatusCode, string(body))
	}

	var result JourneysResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Journeys, nil
}

func (c *Client) CreateJourney(name string) (*models.Journey, error) {
	req := CreateJourneyRequest{
		Name: name,
	}

	resp, err := c.doRequest("POST", "/journey", req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to create journey: status %d, body: %s", resp.StatusCode, string(body))
	}

	var result JourneyResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result.Journey, nil
}

func (c *Client) UpdateJourney(journeyID int, updates UpdateJourneyRequest) (*models.Journey, error) {
	path := fmt.Sprintf("/journey/%d", journeyID)
	resp, err := c.doRequest("PUT", path, updates)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to update journey: status %d, body: %s", resp.StatusCode, string(body))
	}

	var result JourneyResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result.Journey, nil
}

func (c *Client) DeleteJourney(journeyID int) error {
	path := fmt.Sprintf("/journey/%d", journeyID)
	resp, err := c.doRequest("DELETE", path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete journey: status %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}
