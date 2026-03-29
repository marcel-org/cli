package api

import (
	"encoding/json"
	"fmt"
	"io"

	"marcel-cli/models"
)

type SpacesResponse struct {
	Spaces []models.Space `json:"spaces"`
}

type SpaceResponse struct {
	Space models.Space `json:"space"`
}

type CreateSpaceRequest struct {
	Name string `json:"name"`
}

type UpdateSpaceRequest struct {
	Name *string `json:"name,omitempty"`
}

func (c *Client) GetSpaces() ([]models.Space, error) {
	resp, err := c.doRequest("GET", "/space", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get spaces: status %d, body: %s", resp.StatusCode, string(body))
	}

	var result SpacesResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Spaces, nil
}

func (c *Client) CreateSpace(name string) (*models.Space, error) {
	resp, err := c.doRequest("POST", "/space", CreateSpaceRequest{Name: name})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to create space: status %d, body: %s", resp.StatusCode, string(body))
	}

	var result SpaceResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result.Space, nil
}

func (c *Client) UpdateSpace(spaceID int, updates UpdateSpaceRequest) (*models.Space, error) {
	path := fmt.Sprintf("/space/%d", spaceID)
	resp, err := c.doRequest("PUT", path, updates)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to update space: status %d, body: %s", resp.StatusCode, string(body))
	}

	var result SpaceResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result.Space, nil
}
