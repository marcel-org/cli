package api

import (
	"encoding/json"
	"fmt"
	"io"

	"marcel-cli/models"
)

type HabitsResponse struct {
	Habits []models.Habit `json:"habits"`
}

type HabitResponse struct {
	Habit models.Habit `json:"habit"`
}

type CreateHabitRequest struct {
	Name       string `json:"name"`
	CycleType  string `json:"cycleType"`
	CycleConfig any   `json:"cycleConfig,omitempty"`
}

type UpdateHabitRequest struct {
	Name          *string `json:"name,omitempty"`
	CompleteToday *bool   `json:"completeToday,omitempty"`
	CycleType     *string `json:"cycleType,omitempty"`
	CycleConfig   any     `json:"cycleConfig,omitempty"`
}

func (c *Client) GetHabits() ([]models.Habit, error) {
	resp, err := c.doRequest("GET", "/habit", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get habits: status %d, body: %s", resp.StatusCode, string(body))
	}

	var result HabitsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Habits, nil
}

func (c *Client) CreateHabit(name, cycleType string, cycleConfig any) (*models.Habit, error) {
	req := CreateHabitRequest{
		Name:       name,
		CycleType:  cycleType,
		CycleConfig: cycleConfig,
	}

	resp, err := c.doRequest("POST", "/habit", req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to create habit: status %d, body: %s", resp.StatusCode, string(body))
	}

	var result HabitResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result.Habit, nil
}

func (c *Client) UpdateHabit(habitID int, updates UpdateHabitRequest) (*models.Habit, error) {
	path := fmt.Sprintf("/habit/%d", habitID)
	resp, err := c.doRequest("PUT", path, updates)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to update habit: status %d, body: %s", resp.StatusCode, string(body))
	}

	var result HabitResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result.Habit, nil
}

func (c *Client) ToggleHabit(habitID int, completeToday bool) (*models.Habit, error) {
	return c.UpdateHabit(habitID, UpdateHabitRequest{
		CompleteToday: &completeToday,
	})
}

func (c *Client) DeleteHabit(habitID int) error {
	path := fmt.Sprintf("/habit/%d", habitID)
	resp, err := c.doRequest("DELETE", path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete habit: status %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}
