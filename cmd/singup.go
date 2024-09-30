// cmd/signup.go
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// SignupRequest represents the JSON payload for the sign-up API
type SignupRequest struct {
	Email string     `json:"email"`
	Plan  SignupPlan `json:"plan"`
}

// SignupResponse represents the response from the sign-up API
type SignupResponse struct {
	APIKey    string     `json:"api-key"`
	StripeURL string     `json:"stripe-url"`
	AccountID string     `json:"account-id"`
	Plan      SignupPlan `json:"plan"`
}

// Define the signup command
var signupCmd = &cobra.Command{
	Use:   "signup",
	Short: "Sign up for a Global Blackbox account",
	Long:  `Interactively sign up for a Global Blackbox account by providing your email and selecting a subscription plan.`,
	Run: func(cmd *cobra.Command, args []string) {
		runSignup()
	},
}

func init() {
	// You can define flags and configuration settings here if needed in the future
}

func runSignup() {
	// Welcome Message
	welcomeStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#D3D3D3")) // Light Grey
	fmt.Println(welcomeStyle.Render("Welcome to GBX Sign-Up CLI"))

	// Display Plans Documentation Link
	docStyle := lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("#A9A9A9")) // Medium Grey
	fmt.Println(docStyle.Render("For more information on subscription plans, visit: https://docs.globalblackbox.io/plans/\n"))

	// Prompt for Email
	email, err := promptEmail()
	if err != nil {
		exitWithError(err)
	}

	// Prompt for Subscription Plan
	planName, err := promptPlan()
	if err != nil {
		exitWithError(err)
	}

	// If 'single-region', prompt for region
	var region string
	if planName == "single-region" {
		region, err = promptRegion()
		if err != nil {
			exitWithError(err)
		}
	}

	// Construct Signup Request
	signupReq := SignupRequest{
		Email: email,
		Plan: SignupPlan{
			Name:   planName,
			Region: region,
		},
	}

	// Send Signup Request
	response, err := sendSignupRequest(signupReq)
	if err != nil {
		exitWithError(err)
	}

	// Display Response
	displayResponse(response)
}

// promptEmail prompts the user to enter their email address
func promptEmail() (string, error) {
	validate := func(input string) error {
		input = strings.TrimSpace(input)
		if input == "" {
			return fmt.Errorf("email cannot be empty")
		}
		if !strings.Contains(input, "@") || !strings.Contains(input, ".") {
			return fmt.Errorf("invalid email address")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Enter your email address",
		Validate: validate,
	}

	email, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return email, nil
}

// promptPlan prompts the user to select a subscription plan
func promptPlan() (string, error) {
	plans := []string{"single-region", "all-continents", "worldwide"}

	prompt := promptui.Select{
		Label: "Select a subscription plan",
		Items: plans,
		Size:  len(plans),
	}

	_, result, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return result, nil
}

// promptRegion prompts the user to enter a region if "single-region" plan is selected
func promptRegion() (string, error) {
	validate := func(input string) error {
		if strings.TrimSpace(input) == "" {
			return fmt.Errorf("region cannot be empty")
		}
		// Optionally, you can add more validation for region format
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Enter your desired region (e.g., sao-paulo.americas)",
		Validate: validate,
	}

	region, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return region, nil
}

// sendSignupRequest sends the signup request to the API and returns the response
func sendSignupRequest(req SignupRequest) (*SignupResponse, error) {
	url := "https://api.globalblackbox.io/sign-up"

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %v", err)
	}

	// Inform the user that the request is being submitted
	fmt.Println("\nSubmitting your sign-up request...")

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Read response body for more details
		var errorResp map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			return nil, fmt.Errorf("sign-up failed with status: %s", resp.Status)
		}
		return nil, fmt.Errorf("sign-up failed: %v", errorResp)
	}

	var signupResp SignupResponse
	if err := json.NewDecoder(resp.Body).Decode(&signupResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &signupResp, nil
}

// displayResponse displays the API response in a user-friendly format
func displayResponse(resp *SignupResponse) {
	// Update styling to use grey shades
	style := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#A9A9A9")) // Medium Grey

	fmt.Println("\n" + style.Render("Sign-Up Successful!\n"))

	fmt.Printf("%s: %s\n", style.Render("Account ID"), resp.AccountID)
	fmt.Printf("%s: %s\n", style.Render("API Key"), resp.APIKey)
	fmt.Printf("%s: %s\n", style.Render("Stripe URL"), resp.StripeURL)
	fmt.Printf("%s: %s\n", style.Render("Plan Name"), resp.Plan.Name)

	if resp.Plan.Name == "single-region" {
		fmt.Printf("%s: %s\n", style.Render("Region"), resp.Plan.Region)
	}

	// Save API key to config.yaml
	config := &Config{
		APIKey:    resp.APIKey,
		AccountID: resp.AccountID,
		Plan: SignupPlan{
			Name:   resp.Plan.Name,
			Region: resp.Plan.Region,
		},
	}

	if err := SaveConfig(config); err != nil {
		errorStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#696969")) // Dark Grey
		fmt.Fprintf(os.Stderr, "%s: %v\n", errorStyle.Render("Warning"), err)
	} else {
		fmt.Println(style.Render("\nAPI key has been saved to ~/.gbx/config.yaml"))
	}

	// Next Steps
	nextStepsStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#D3D3D3")) // Light Grey
	fmt.Println("\n" + nextStepsStyle.Render("Next Steps:"))
	fmt.Println("1. Complete Subscription Payment by visiting the Stripe URL provided.")
	fmt.Println("2. Secure your API Key for authenticating your Prometheus scrape jobs.")
	fmt.Println("3. Configure Prometheus with your account details. Refer to the Prometheus Configuration documentation for guidance.\n")

	// Support Information
	supportStyle := lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("#A9A9A9")) // Medium Grey
	fmt.Println(supportStyle.Render("For support, contact support@globalblackbox.io"))
}

// exitWithError prints the error and exits the application
func exitWithError(err error) {
	errorStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#696969")) // Dark Grey
	fmt.Fprintf(os.Stderr, "%s: %v\n", errorStyle.Render("Error"), err)
	os.Exit(1)
}
