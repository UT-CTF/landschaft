package score

// PortExplain returns a short explanation for a given port (e.g. "SSH", "HTTP").
func PortExplain(port uint16) string {
	switch port {
	case 22:
		return "SSH"
	case 80:
		return "HTTP"
	case 443:
		return "HTTPS"
	case 3389:
		return "RDP"
	case 389:
		return "LDAP"
	case 636:
		return "LDAPS"
	case 88:
		return "Kerberos"
	case 53:
		return "DNS"
	case 25:
		return "SMTP"
	case 110:
		return "POP3"
	case 143:
		return "IMAP"
	case 445:
		return "SMB"
	case 3306:
		return "MySQL"
	case 5432:
		return "PostgreSQL"
	case 5985:
		return "WinRM HTTP"
	case 5986:
		return "WinRM HTTPS"
	case 135:
		return "RPC"
	case 139:
		return "NetBIOS"
	default:
		return ""
	}
}
