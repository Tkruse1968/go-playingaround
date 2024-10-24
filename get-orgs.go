package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/time/rate"
)

const (
	snykAPIURL = "https://api.snyk.io/rest/orgs/%s/projects?version=2024-10-24"
	// Set the desired rate limit (e.g., 10 requests per minute)
	rateLimit         = 10
	rateLimitDuration = time.Minute
)

type Project struct {
	Name string `json:"name"`
}

func main() {
	if len(os.Args) != 3 {
		log.Fatal("Usage: go run main.go <API_TOKEN> <ORG_ID>")
	}

	apiToken := os.Args[1]
	orgID := os.Args[2]

	limiter := rate.NewLimiter(rate.Every(rateLimitDuration), rateLimit)

	projects, err := getProjects(limiter, apiToken, orgID)
	if err != nil {
		log.Fatal(err)
	}

	err = exportProjectsToCSV(projects)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Projects exported to projects.csv successfully!")
}

func getProjects(limiter *rate.Limiter, apiToken, orgID string) ([]Project, error) {
	url := fmt.Sprintf(snykAPIURL, orgID)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("token %s", apiToken))
	req.Header.Set("Content-Type", "application/vnd.api+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Assuming the response is a JSON array of projects
	var projects []Project
	err = json.Unmarshal(body, &projects)
	if err != nil {
		return nil, err
	}

	return projects, nil
}

func exportProjectsToCSV(projects []Project) error {
	file, err := os.Create("projects.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header row
	err = writer.Write([]string{"Name"})
	if err != nil {
		return err
	}

	// Write project data
	for _, project := range projects {
		err = writer.Write([]string{project.Name})
		if err != nil {
			return err
		}
	}

	return nil
}
