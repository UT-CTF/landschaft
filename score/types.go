package score

// Listener is a listening endpoint (candidate scored service).
type Listener struct {
	Port    uint16
	Proto   string
	Process string
	Bind    string
	Explain string
}
