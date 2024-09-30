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
	"globalblackbox.io/globalblackbox-cli/models"
)

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

// runSignup orchestrates the sign-up process
func runSignup() {
	// Welcome Message
	welcomeStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#D3D3D3")) // Light Grey
	fmt.Println(welcomeStyle.Render("Welcome to GBX Sign-Up CLI"))

	// Display Pricing Information
	displayPricingInfo()

	// Display Plans Documentation Link
	docStyle := lipgloss.NewStyle().
		Italic(true).
		Foreground(lipgloss.Color("#A9A9A9")) // Medium Grey
	fmt.Println(docStyle.Render("For more information on subscription plans, visit: https://docs.globalblackbox.io/plans/\n"))

	// Prompt for Email
	email, err := promptEmail()
	if err != nil {
		exitWithError(err)
	}

	var planName string
	for {
		// Prompt for Subscription Plan
		planName, err = promptPlan()
		if err != nil {
			exitWithError(err)
		}

		// Confirm the Selected Plan
		confirmed, err := confirmPlan(planName)
		if err != nil {
			exitWithError(err)
		}

		if confirmed {
			break
		} else {
			// User chose to re-select plan
			fmt.Println("\n" + lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#696969")).
				Render("Let's re-select your subscription plan.\n"))
		}
	}

	// If 'single-region', prompt for region
	var region string
	if planName == "single-region" {
		region, err = promptRegion()
		if err != nil {
			exitWithError(err)
		}
	}

	// Prompt for Number of Targets
	numberOfTargets, err := promptNumberOfTargets()
	if err != nil {
		exitWithError(err)
	}

	// Construct Signup Request
	signupReq := models.SignupRequest{
		Email: email,
		Plan: models.SignupPlan{
			Name:   planName,
			Region: region,
		},
		NumberOfTargets: numberOfTargets,
	}

	// Send Signup Request
	response, err := sendSignupRequest(signupReq)
	if err != nil {
		exitWithError(err)
	}

	// Display Response
	displayResponse(response)
}

// displayPricingInfo displays information about how pricing works
func displayPricingInfo() {
	pricingStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#A9A9A9")) // Medium Grey
	fmt.Println(pricingStyle.Render("Pricing Structure:\n"))
	fmt.Println(`- Available plans: single-region, all-continents, or worldwide.
- Number of targets: The number of endpoints you wish to monitor.

The total cost is determined by the combination of your selected plan and the number of targets. Detailed pricing information is available during the subscription process via the Stripe payment link.`)
	fmt.Println()
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

// confirmPlan displays the plan details and prompts the user to confirm or re-select
func confirmPlan(planName string) (bool, error) {
	description, exists := models.PlanDetails[planName]
	if !exists {
		return false, fmt.Errorf("no details found for the selected plan: %s", planName)
	}

	// Display Plan Details
	planStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#A9A9A9")) // Medium Grey
	fmt.Println("\n" + planStyle.Render("Selected Plan Details:\n"))
	fmt.Println(description)

	// Prompt for Confirmation
	confirmPrompt := promptui.Select{
		Label: "Do you want to proceed with this plan?",
		Items: []string{"Confirm", "Re-select"},
	}

	_, result, err := confirmPrompt.Run()
	if err != nil {
		return false, fmt.Errorf("confirmation prompt failed: %v", err)
	}

	if result == "Confirm" {
		return true, nil
	} else {
		return false, nil
	}
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

// promptNumberOfTargets prompts the user to enter the desired number of probe targets
func promptNumberOfTargets() (int, error) {
	validate := func(input string) error {
		input = strings.TrimSpace(input)
		if input == "" {
			return fmt.Errorf("number of targets cannot be empty")
		}
		var number int
		_, err := fmt.Sscanf(input, "%d", &number)
		if err != nil || number <= 0 {
			return fmt.Errorf("please enter a valid positive integer for the number of targets")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Enter the number of probe targets you wish to monitor",
		Validate: validate,
	}

	input, err := prompt.Run()
	if err != nil {
		return 0, err
	}

	var numberOfTargets int
	_, err = fmt.Sscanf(strings.TrimSpace(input), "%d", &numberOfTargets)
	if err != nil {
		return 0, fmt.Errorf("failed to parse number of targets: %v", err)
	}

	return numberOfTargets, nil
}

// sendSignupRequest sends the signup request to the API and returns the response
func sendSignupRequest(req models.SignupRequest) (*models.SignupResponse, error) {
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

	var signupResp models.SignupResponse
	if err := json.NewDecoder(resp.Body).Decode(&signupResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &signupResp, nil
}

// displayResponse displays the API response in a user-friendly format
func displayResponse(resp *models.SignupResponse) {
	// Update styling to use grey shades
	style := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#A9A9A9")) // Medium Grey

	fmt.Println("\n" + style.Render("Sign-Up Successful!\n"))

	fmt.Printf("%s: %s\n", style.Render("Account ID"), resp.AccountID)
	fmt.Printf("%s: %s\n", style.Render("API Key"), resp.APIKey)
	fmt.Printf("%s: %s\n", style.Render("Stripe URL"), resp.StripeURL)
	fmt.Printf("%s: %s\n", style.Render("Plan Name"), resp.Plan.Name)

	if resp.Plan.Name == "single-region" {
		fmt.Printf("%s: %s\n", style.Render("Region"), resp.Plan.Region)
	}

	// Save API key to config.yaml
	config := &models.Config{
		APIKey:    resp.APIKey,
		AccountID: resp.AccountID,
		Plan: models.SignupPlan{
			Name:   resp.Plan.Name,
			Region: resp.Plan.Region,
		},
		NumberOfTargets: resp.NumberOfTargets, // Assuming API will eventually return this
	}

	if err := SaveConfig(config); err != nil {
		errorStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#696969")) // Dark Grey
		fmt.Fprintf(os.Stderr, "%s: %v\n", errorStyle.Render("Warning"), err)
	} else {
		fmt.Println(style.Render("\nAPI key has been saved to ~/.gbx/config.yaml"))
	}

	// Next Steps
	nextStepsStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#D3D3D3")) // Light Grey
	fmt.Println("\n" + nextStepsStyle.Render("Next Steps:"))
	fmt.Println("1. Complete Subscription Payment by visiting the Stripe URL provided.")
	fmt.Println("2. Secure your API Key for authenticating your Prometheus scrape jobs.")
	fmt.Println("3. Configure Prometheus with your account details. Refer to the Prometheus Configuration documentation for guidance.\n")

	// Support Information
	supportStyle := lipgloss.NewStyle().
		Italic(true).
		Foreground(lipgloss.Color("#A9A9A9")) // Medium Grey
	fmt.Println(supportStyle.Render("For support, contact support@globalblackbox.io"))
}

// exitWithError prints the error and exits the application
func exitWithError(err error) {
	errorStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#696969")) // Dark Grey
	fmt.Fprintf(os.Stderr, "%s: %v\n", errorStyle.Render("Error"), err)
	os.Exit(1)
}
