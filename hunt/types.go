package hunt

// Detection is one normalized detection event (suspicious or notable).
type Detection struct {
	Timestamp string   `json:"timestamp"`
	Host      string   `json:"host"`
	OS        string   `json:"os"`
	Source    string   `json:"source"`
	EventID   string   `json:"event_id,omitempty"`
	Account   string   `json:"account,omitempty"`
	RemoteIP  string   `json:"remote_ip,omitempty"`
	Process   string   `json:"process,omitempty"`
	Message   string   `json:"message"`
	Severity  string   `json:"severity"`
	Tags      []string `json:"tags,omitempty"`
	Explain   string   `json:"explain,omitempty"`
}
