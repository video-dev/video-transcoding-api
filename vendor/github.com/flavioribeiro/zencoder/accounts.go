package zencoder

type AccountDetails struct {
	AccountState    string `json:"account_state,omitempty"`
	Plan            string `json:"plan,omitempty"`
	MinutesUsed     int32  `json:"minutes_used,omitempty"`
	MinutesIncluded int32  `json:"minutes_included,omitempty"`
	BillingState    string `json:"billing_state,omitempty"`
	IntegrationMode bool   `json:"integration_mode,omitempty"`
}

type CreateAccountRequest struct {
	Email                string  `json:"email,omitempty"`
	TermsOfService       string  `json:"terms_of_service,omitempty"`
	Password             *string `json:"password,omitempty,omitempty"`
	PasswordConfirmation *string `json:"password_confirmation,omitempty,omitempty"`
}

type CreateAccountResponse struct {
	ApiKey   string `json:"api_key,omitempty"`
	Password string `json:"password,omitempty"`
}

// Create an account
func (z *Zencoder) CreateAccount(email, password string) (*CreateAccountResponse, error) {
	request := &CreateAccountRequest{
		Email:          email,
		TermsOfService: "1",
	}

	if len(password) > 0 {
		request.Password = &password
		request.PasswordConfirmation = &password
	}

	var result CreateAccountResponse

	if err := z.post("account", request, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Get Account Details
func (z *Zencoder) GetAccount() (*AccountDetails, error) {
	var details AccountDetails

	if err := z.getBody("account", &details); err != nil {
		return nil, err
	}

	return &details, nil
}

// Set Integration Mode
func (z *Zencoder) SetIntegrationMode() error {
	return z.putNoContent("account/integration")
}

// Set Live Mode
func (z *Zencoder) SetLiveMode() error {
	return z.putNoContent("account/live")
}
