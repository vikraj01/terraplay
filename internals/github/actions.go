package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const GitHubAPIURL = "https://api.github.com/repos/%s/%s/actions/workflows/%s/dispatches"

type DispatchRequest struct {
	Ref    string            `json:"ref"`
	Inputs map[string]string `json:"inputs,omitempty"`
}

func TriggerGithubAction(owner, repo, workflowID, ref string, inputs map[string]string) error {
	token := os.Getenv("GITHUB_TOKEN")

	url := fmt.Sprintf(GitHubAPIURL, owner, repo, workflowID)

	payload := DispatchRequest{
		Ref:    ref,
		Inputs: inputs,
	}

	jsonData, err := json.Marshal(payload)

	if err != nil {
		return fmt.Errorf("failed to marshal request: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("GitHub Action failed with status code: %d", resp.StatusCode)
	}

	fmt.Println("GitHub Action triggered successfully!")
	return nil
}
