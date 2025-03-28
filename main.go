// This program backs up Grafana dashboards by fetching them via the Grafana API
// and saving them as JSON files in a directory structure that mirrors the
// Grafana folder organization.

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// Get Dashboard uid
type SearchResult struct {
	UID string `json:"uid"`
}

// Get Dashboard folder name
type DashboardInfo struct {
	Dashboard struct {
		Title string `json:"title"`
	} `json:"dashboard"`
	Meta struct {
		FolderTitle string `json:"folderTitle"`
	} `json:"meta"`
}

// Get Enviroment
func getRequiredEnv(key string) (string, error) {
	val := os.Getenv(key)
	if val == "" {
		return "", fmt.Errorf("required environment variable %s is not set", key)
	}
	return val, nil
}

// Get Authorization header and parse JSON
func fetchJSON(url, token string, target interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("request failed: %s\n%s", res.Status, body)
	}
	return json.NewDecoder(res.Body).Decode(target)
}

// Dashboard list uid
func fetchDashboardList(baseURL, token string) ([]SearchResult, error) {
	var dashboards []SearchResult
	url := fmt.Sprintf("%s/api/search?type=dash-db", baseURL)
	err := fetchJSON(url, token, &dashboards)
	return dashboards, err
}

// Get use uid Dashboard info 
func fetchDashboardDetail(baseURL, token, uid string) (*DashboardInfo, error) {
	var detail DashboardInfo
	url := fmt.Sprintf("%s/api/dashboards/uid/%s", baseURL, uid)
	err := fetchJSON(url, token, &detail)
	return &detail, err
}

// Replace characters that can’t be used in file names (like / or spaces) with -
func sanitizeFileName(title string) string {
	title = strings.ReplaceAll(title, "/", "-")
	title = strings.ReplaceAll(title, " ", "-")
	return title
}

// Get Dashboard JSON and save file
func saveDashboardRawJSON(baseURL, token, uid, folder, filename string) error {
	url := fmt.Sprintf("%s/api/dashboards/uid/%s", baseURL, uid)
	var raw map[string]interface{}
	if err := fetchJSON(url, token, &raw); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}
	outfile, err := os.Create(fmt.Sprintf("%s/%s", folder, filename))
	if err != nil {
		return err
	}
	defer outfile.Close()
	enc := json.NewEncoder(outfile)
	enc.SetIndent("", "  ")
	return enc.Encode(raw)
}

// make a one Dashboard process
func saveDashboard(baseURL, token, outputDir, uid string) {
	detail, err := fetchDashboardDetail(baseURL, token, uid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed Dashboard: %v\n", err)
		return
	}

	title := sanitizeFileName(detail.Dashboard.Title)
	filename := fmt.Sprintf("%s_%s.json", title, uid)

	folder := fmt.Sprintf("%s/%s", outputDir, detail.Meta.FolderTitle)
	os.MkdirAll(folder, 0755)

	if err := saveDashboardRawJSON(baseURL, token, uid, folder, filename); err != nil {
		fmt.Fprintf(os.Stderr, "failed saving dashboard: %v\n", err)
	}
}

func main() {
	baseURL, err := getRequiredEnv("MONITORING_URL")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	apiToken, err := getRequiredEnv("GRAFANA_API_KEY")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	outputDir := "./dashboard"
	os.RemoveAll(outputDir)

	dashboards, err := fetchDashboardList(baseURL, apiToken)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed search API: %v\n", err)
		os.Exit(1)
	}

	for _, d := range dashboards {
		saveDashboard(baseURL, apiToken, outputDir, d.UID)
	}
}
