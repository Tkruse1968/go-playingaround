package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
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
	// ... (rest of the code remains the same)

	// Assuming the response is a JSON object with a "data" field containing an array of projects
	type ProjectResponse struct {
		Data []Project `json:"data"`
	}

	var projectResponse ProjectResponse
	err = json.Unmarshal(body, &projectResponse)
	if err != nil {
		return nil, err
	}

	return projectResponse.Data, nil
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
