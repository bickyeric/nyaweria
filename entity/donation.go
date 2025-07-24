package entity

type Donation struct {
	ID          string
	From        string `json:"from"`
	To          string `json:"to"`
	Amount      string `json:"amount"`
	Message     string `json:"message"`
	AudioPath   string `json:"audio_path"`
	RecipientID string
}

type DonationSummary struct{}
