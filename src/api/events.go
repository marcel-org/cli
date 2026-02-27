package api

import (
	"encoding/json"
	"fmt"
	"io"

	"marcel-cli/models"
)

type EventsResponse struct {
	Events []models.Event `json:"events"`
}

type EventResponse struct {
	Event models.Event `json:"event"`
}

type CreateEventRequest struct {
	Title       string  `json:"title"`
	Date        string  `json:"date"`
	EndDate     *string `json:"endDate,omitempty"`
	Time        *string `json:"time,omitempty"`
	EndTime     *string `json:"endTime,omitempty"`
	Location    *string `json:"location,omitempty"`
	Description *string `json:"description,omitempty"`
}

type UpdateEventRequest struct {
	Title       *string `json:"title,omitempty"`
	Date        *string `json:"date,omitempty"`
	EndDate     *string `json:"endDate,omitempty"`
	Time        *string `json:"time,omitempty"`
	EndTime     *string `json:"endTime,omitempty"`
	Location    *string `json:"location,omitempty"`
	Description *string `json:"description,omitempty"`
}

func (c *Client) GetEvents() ([]models.Event, error) {
	resp, err := c.doRequest("GET", "/event", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get events: status %d, body: %s", resp.StatusCode, string(body))
	}

	var result EventsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Events, nil
}

func (c *Client) CreateEvent(req CreateEventRequest) (*models.Event, error) {
	resp, err := c.doRequest("POST", "/event", req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to create event: status %d, body: %s", resp.StatusCode, string(body))
	}

	var result EventResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result.Event, nil
}

func (c *Client) UpdateEvent(eventID int, updates UpdateEventRequest) (*models.Event, error) {
	path := fmt.Sprintf("/event/%d", eventID)
	resp, err := c.doRequest("PUT", path, updates)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to update event: status %d, body: %s", resp.StatusCode, string(body))
	}

	var result EventResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result.Event, nil
}

func (c *Client) DeleteEvent(eventID int) error {
	path := fmt.Sprintf("/event/%d", eventID)
	resp, err := c.doRequest("DELETE", path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete event: status %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}
