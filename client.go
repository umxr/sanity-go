package sanity

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// SanityClient represents the client for Sanity API.
type SanityClient struct {
	ProjectID  string
	Dataset    string
	Token      string
	APIVersion string
}

type ApiResponse struct {
	Result json.RawMessage `json:"result"`
}

// CreateClient creates a new Sanity API client.
func CreateClient(projectID, dataset, token, apiVersion string) (*SanityClient, error) {
	if projectID == "" {
		return nil, fmt.Errorf("projectID is required")
	}
	if dataset == "" {
		return nil, fmt.Errorf("dataset is required")
	}

	// Set defaults for optional parameters
	if apiVersion == "" {
		apiVersion = "v" + time.Now().Format("2006-01-02") // Go uses this specific date as a format layout
	}

	return &SanityClient{
		ProjectID:  projectID,
		Dataset:    dataset,
		Token:      token,
		APIVersion: apiVersion,
	}, nil
}

func (c *SanityClient) Fetch(groqQuery string, queryParams ...map[string]string) (string, error) {
	// Check if the GROQ query is provided
	if groqQuery == "" {
		return "", fmt.Errorf("please provide a query")
	}

	// Construct the base URL for the Sanity API
	baseURL := fmt.Sprintf("https://%s.api.sanity.io/%s/data/query/%s", c.ProjectID, c.APIVersion, c.Dataset)

	// Construct the full URL
	fullURL := fmt.Sprintf("%s?query=%s", baseURL, url.QueryEscape(groqQuery))

	// Append additional parameters if provided
	if len(queryParams) > 0 {
		for key, value := range queryParams[0] {
			// Prepend '$' to the key and append the parameter to the URL
			fullURL += fmt.Sprintf("&$%s=\"%s\"", key, url.QueryEscape(value))
		}
	}

	// Create a new HTTP GET request with the constructed URL
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	// Add Authorization header if a token is provided
	if c.Token != "" {
		req.Header.Add("Authorization", "Bearer "+c.Token)
	}

	// Execute the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error executing request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	// Parse the JSON response
	var apiResponse ApiResponse
	err = json.Unmarshal(responseBody, &apiResponse)
	if err != nil {
		return "", fmt.Errorf("error parsing JSON: %w", err)
	}

	// Return the result part of the response
	return string(apiResponse.Result), nil
}

// Clone creates a new instance of SanityClient with the same configuration.
func (c *SanityClient) Clone() *SanityClient {
	return &SanityClient{
		ProjectID:  c.ProjectID,
		Dataset:    c.Dataset,
		Token:      c.Token,
		APIVersion: c.APIVersion,
	}
}

// Returns the current configuration of the client.
func (c *SanityClient) Config() SanityClient {
	return SanityClient{
		ProjectID:  c.ProjectID,
		Dataset:    c.Dataset,
		Token:      c.Token,
		APIVersion: c.APIVersion,
	}
}
