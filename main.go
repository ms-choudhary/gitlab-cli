package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

var (
	GitlabBaseURL = os.Getenv("GITLAB_BASE_URL")
	AccessToken   = os.Getenv("PERSONAL_ACCESS_TOKEN")
)

type Issue struct {
	ID          int `json:"id"`
	ProjectID   int `json:"project_id"`
	Author      User
	Assignee    User
	Title       string
	Description string
}

type User struct {
	Name string
}

func issuesByProjectPath(p int) string {
	return fmt.Sprintf("%s/api/v4/projects/%d/issues", GitlabBaseURL, p)
}

func ListIssues(projectID int) ([]Issue, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", issuesByProjectPath(projectID), nil)
	if err != nil {
		return []Issue{}, err
	}

	req.Header.Add("PRIVATE-TOKEN", AccessToken)

	resp, err := client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		return []Issue{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return []Issue{}, fmt.Errorf("failed to list issues with status: %s", resp.Status)
	}

	var result []Issue
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return []Issue{}, fmt.Errorf("failed to decode: %v", err)
	}

	return result, nil
}

func main() {
	issues, err := ListIssues(76)
	if err != nil {
		log.Fatal(err)
	}

	for _, iss := range issues {
		fmt.Printf("%3d %3d %9.9s %9.9s %.55s\n", iss.ProjectID, iss.ID, iss.Author, iss.Assignee, iss.Title)
	}
}
