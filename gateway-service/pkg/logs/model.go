package logs

type LogEntry struct {
	Date string `json:"date"`
	Log  string `json:"log"`
}

type LogsResponse struct {
	Logs []LogEntry `json:"logs"`
}
