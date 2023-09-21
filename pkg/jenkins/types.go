package jenkins

type PendingAction struct {
	ID          string `json:"id"`
	Message     string `json:"message"`
	ProceedURL  string `json:"proceedUrl"`
	AbortURL    string `json:"abortUrl"`
	ProceedText string `json:"proceedText"`
}
