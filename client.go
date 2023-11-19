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

// Fetch executes a GROQ query and returns the results. (Placeholder implementation)
func (c *SanityClient) Fetch(groqQuery string, queryParams ...map[string]string) (string, error) {
	if groqQuery == "" {
		return "", fmt.Errorf("please provide a query")
	}

	// Construct the URL for the Sanity API
	baseURL := fmt.Sprintf("https://%s.api.sanity.io/%s/data/query/%s", c.ProjectID, c.APIVersion, c.Dataset)

	// Prepare the full URL with the encoded query and additional parameters
	params := url.Values{}
	params.Set("query", groqQuery)

	if len(queryParams) > 0 {
		for key, value := range queryParams[0] {
			params.Set(key, value)
		}
	}

	fullURL := baseURL + "?" + params.Encode()

	// Create a new HTTP request
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	if c.Token != "" {
		// Add necessary headers
		req.Header.Add("Authorization", "Bearer "+c.Token)
	}

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error executing request: %w", err)
	}
	defer resp.Body.Close()

	// Read and return the response
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	var apiResponse ApiResponse
	err = json.Unmarshal(responseBody, &apiResponse)
	if err != nil {
		return "", fmt.Errorf("error parsing JSON: %w", err)
	}

	return string(apiResponse.Result), nil
}
