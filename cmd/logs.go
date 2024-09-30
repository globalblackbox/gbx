// cmd/logs.go
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"globalblackbox.io/globalblackbox-cli/models"
)

// Define the logs command
var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Retrieve and download logs from Global Blackbox",
	Long:  `Interact with Global Blackbox logs by listing available log files or downloading specific logs.`,
}

// Define the list subcommand
var logsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available log files",
	Long:  `List available log files based on region, target domain, and date.`,
	Run: func(cmd *cobra.Command, args []string) {
		runLogsList(cmd, args)
	},
}

// Define the download subcommand
var logsDownloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download a specific log file",
	Long:  `Download a specific log file by providing the file name, region, target domain, and date.`,
	Run: func(cmd *cobra.Command, args []string) {
		runLogsDownload(cmd, args)
	},
}

func init() {
	// Add list and download as subcommands of logs
	logsCmd.AddCommand(logsListCmd)
	logsCmd.AddCommand(logsDownloadCmd)

	// Define flags for the list command
	logsListCmd.Flags().StringP("region", "r", "", "Region code (e.g., london.europe) (required)")
	logsListCmd.Flags().StringP("target_domain", "t", "", "Target domain (e.g., booking.com) (required)")
	logsListCmd.Flags().StringP("date", "d", "", "Date in YYYY-MM-DD format (required)")
	logsListCmd.Flags().IntP("limit", "l", 10, "Number of log files to retrieve (max 50)")
	logsListCmd.MarkFlagRequired("region")
	logsListCmd.MarkFlagRequired("target_domain")
	logsListCmd.MarkFlagRequired("date")

	// Define flags for the download command
	logsDownloadCmd.Flags().StringP("fileName", "f", "", "Name of the log file to download (required)")
	logsDownloadCmd.Flags().StringP("region", "r", "", "Region code (e.g., london.europe) (required)")
	logsDownloadCmd.Flags().StringP("target_domain", "t", "", "Target domain (e.g., booking.com) (required)")
	logsDownloadCmd.Flags().StringP("date", "d", "", "Date in YYYY-MM-DD format (required)")
	logsDownloadCmd.MarkFlagRequired("fileName")
	logsDownloadCmd.MarkFlagRequired("region")
	logsDownloadCmd.MarkFlagRequired("target_domain")
	logsDownloadCmd.MarkFlagRequired("date")
}

// runLogsList handles the 'logs list' command
func runLogsList(cmd *cobra.Command, args []string) {
	// Retrieve flags
	region, _ := cmd.Flags().GetString("region")
	targetDomain, _ := cmd.Flags().GetString("target_domain")
	date, _ := cmd.Flags().GetString("date")
	limit, _ := cmd.Flags().GetInt("limit")

	// Validate date format
	if err := validateDate(date); err != nil {
		exitWithError(err)
	}

	// Validate limit
	if limit > 50 {
		fmt.Println("Limit cannot exceed 50. Setting limit to 50.")
		limit = 50
	}

	// Retrieve API key from config
	apiKey, err := getAPIKey()
	if err != nil {
		exitWithError(err)
	}

	// Construct the request URL
	url := fmt.Sprintf("https://api.globalblackbox.io/logs?region=%s&target_domain=%s&date=%s&limit=%d",
		region, targetDomain, date, limit)

	// Create HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		exitWithError(fmt.Errorf("failed to create HTTP request: %v", err))
	}

	// Add Authorization header
	req.Header.Add("x-api-key", apiKey)

	// Execute HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		exitWithError(fmt.Errorf("failed to execute HTTP request: %v", err))
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		var errorResp map[string]interface{}
		bodyBytes, _ := io.ReadAll(resp.Body)
		json.Unmarshal(bodyBytes, &errorResp)
		exitWithError(fmt.Errorf("API request failed with status %s: %v", resp.Status, errorResp))
	}

	// Parse response
	var logsResponse struct {
		LogFiles []string `json:"logs"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&logsResponse); err != nil {
		exitWithError(fmt.Errorf("failed to parse API response: %v", err))
	}

	// Display log files
	if len(logsResponse.LogFiles) == 0 {
		fmt.Println("No log files found for the given parameters.")
		return
	}

	listStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#A9A9A9")) // Medium Grey
	a_str := fmt.Sprintf("\nAvailable log files for %s, target domain %s, and date %s:\n", region, targetDomain, date)
	fmt.Println(listStyle.Render(a_str))
	for i, file := range logsResponse.LogFiles {
		fmt.Printf("%d. %s\n", i+1, file)
	}
	fmt.Println()
}

// runLogsDownload handles the 'logs download' command
func runLogsDownload(cmd *cobra.Command, args []string) {
	// Retrieve flags
	fileName, _ := cmd.Flags().GetString("fileName")
	region, _ := cmd.Flags().GetString("region")
	targetDomain, _ := cmd.Flags().GetString("target_domain")
	date, _ := cmd.Flags().GetString("date")

	// Validate date format
	if err := validateDate(date); err != nil {
		exitWithError(err)
	}

	// Retrieve API key from config
	apiKey, err := getAPIKey()
	if err != nil {
		exitWithError(err)
	}

	// Construct the request URL
	url := fmt.Sprintf("https://api.globalblackbox.io/logs/%s?region=%s&target_domain=%s&date=%s",
		fileName, region, targetDomain, date)

	// Create HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		exitWithError(fmt.Errorf("failed to create HTTP request: %v", err))
	}

	// Add Authorization header
	req.Header.Add("x-api-key", apiKey)

	// Execute HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		exitWithError(fmt.Errorf("failed to execute HTTP request: %v", err))
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		var errorResp map[string]interface{}
		bodyBytes, _ := io.ReadAll(resp.Body)
		json.Unmarshal(bodyBytes, &errorResp)
		exitWithError(fmt.Errorf("API request failed with status %s: %v", resp.Status, errorResp))
	}

	// Create the logs directory if it doesn't exist
	logsDir := "logs"
	if _, err := os.Stat(logsDir); os.IsNotExist(err) {
		if err := os.Mkdir(logsDir, 0755); err != nil {
			exitWithError(fmt.Errorf("failed to create logs directory: %v", err))
		}
	}

	// Create the file
	filePath := filepath.Join(logsDir, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		exitWithError(fmt.Errorf("failed to create file: %v", err))
	}
	defer file.Close()

	// Write the response body to the file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		exitWithError(fmt.Errorf("failed to write to file: %v", err))
	}

	// Inform the user
	downloadStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#D3D3D3")) // Light Grey
	fmt.Printf("\n%s: %s has been downloaded to the '%s' directory.\n\n", downloadStyle.Render("Success"), fileName, logsDir)
}

// getAPIKey retrieves the API key from the configuration file
func getAPIKey() (string, error) {
	// Get the user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("unable to determine home directory: %v", err)
	}

	// Define the config file path
	configFile := filepath.Join(homeDir, ".gbx", "config.yaml")

	// Check if the config file exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return "", fmt.Errorf("config file not found at %s", configFile)
	}

	// Read the config file
	data, err := os.ReadFile(configFile)
	if err != nil {
		return "", fmt.Errorf("failed to read config file: %v", err)
	}

	// Parse the YAML
	var config models.Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return "", fmt.Errorf("failed to parse config file: %v", err)
	}

	// Check if APIKey is present
	if strings.TrimSpace(config.APIKey) == "" {
		return "", fmt.Errorf("API key not found in config file")
	}

	return config.APIKey, nil
}

// validateDate checks if the provided date is in YYYY-MM-DD format
func validateDate(dateStr string) error {
	_, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return fmt.Errorf("invalid date format. Please use YYYY-MM-DD")
	}
	return nil
}
