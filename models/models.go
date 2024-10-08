package models

type SignupPlan struct {
	Name            string `json:"name" yaml:"name"`
	Region          string `json:"region,omitempty" yaml:"region,omitempty"`
	NumberOfTargets int    `json:"number_of_targets"`
}

type SignupRequest struct {
	Email string     `json:"email"`
	Plan  SignupPlan `json:"plan"`
}

type SignupResponse struct {
	APIKey          string     `json:"api-key"`
	StripeURL       string     `json:"stripe-url"`
	AccountID       string     `json:"account-id"`
	Plan            SignupPlan `json:"plan"`
	NumberOfTargets int        `yaml:"number_of_targets"`
}

type Config struct {
	APIKey          string     `yaml:"api_key"`
	AccountID       string     `yaml:"account_id"`
	Plan            SignupPlan `yaml:"plan"`
	NumberOfTargets int        `yaml:"number_of_targets"`
}
