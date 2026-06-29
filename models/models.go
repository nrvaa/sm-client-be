package models

type TimelineEvent struct {
	ID          string `json:"id"`
	Status      string `json:"status"`
	Date        string `json:"date"`
	Description string `json:"description"`
	Image       string `json:"image,omitempty"`
	Completed   bool   `json:"completed"`
}

type ProgressPhoto struct {
	ID    string `json:"id"`
	URL   string `json:"url"`
	Title string `json:"title"`
	Date  string `json:"date"`
	Stage string `json:"stage"`
	Type  string `json:"type,omitempty"`
}

type ReplacedPart struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	From string `json:"from"`
	To   string `json:"to"`
}

type Vehicle struct {
	ID                   string          `json:"id"`
	ClientCode           string          `json:"clientCode"`
	ClientName           string          `json:"clientName"`
	Brand                string          `json:"brand"`
	Model                string          `json:"model"`
	Year                 int             `json:"year"`
	LicensePlate         string          `json:"licensePlate"`
	VIN                  string          `json:"vin"`
	Status               string          `json:"status"`
	CompletionPercentage int             `json:"completionPercentage"`
	EstimatedCompletion  string          `json:"estimatedCompletion"`
	BannerImage          string          `json:"bannerImage"`
	RestorationType      string          `json:"restorationType"`
	Timeline             []TimelineEvent `json:"timeline"`
	Gallery              []ProgressPhoto `json:"gallery"`
	ReplacedParts        []ReplacedPart  `json:"replacedParts,omitempty"`
}

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Slug     string `json:"slug"`
}
