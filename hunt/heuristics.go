package hunt

import "strings"

// TagSuspicious adds tags and explain text based on simple heuristics.
func TagSuspicious(d *Detection) {
	if d.Severity == "" {
		d.Severity = "info"
	}
	msg := d.Message
	if d.EventID == "4625" {
		d.Tags = append(d.Tags, "failed_logon")
		d.Explain = "Failed logon; check for brute force if many from same IP."
		d.Severity = "medium"
	}
	if d.EventID == "4624" {
		d.Tags = append(d.Tags, "logon")
	}
	if d.EventID == "4728" || d.EventID == "4732" || d.EventID == "4756" {
		d.Tags = append(d.Tags, "group_add")
		d.Explain = "User added to privileged group; verify expected."
		d.Severity = "high"
	}
	if d.EventID == "7045" {
		d.Tags = append(d.Tags, "service_installed")
		d.Explain = "New service installed; verify legitimate."
		d.Severity = "medium"
	}
	if d.EventID == "1102" {
		d.Tags = append(d.Tags, "log_cleared")
		d.Explain = "Audit log was cleared; possible tampering."
		d.Severity = "high"
	}
	if len(d.Tags) == 0 && (strings.Contains(msg, "Failed password") || strings.Contains(msg, "Invalid user")) {
		d.Tags = append(d.Tags, "ssh_failed")
		d.Explain = "SSH failed/invalid login; check for brute force."
		d.Severity = "medium"
	}
	if strings.Contains(msg, "Accepted password") || strings.Contains(msg, "Accepted publickey") {
		d.Tags = append(d.Tags, "ssh_success")
	}
}
