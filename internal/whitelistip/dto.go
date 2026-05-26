package whitelistip

// CreateWhitelistIPRequest is the JSON body for creating a whitelist entry.
// IPAddress is stored as PostgreSQL `cidr`: use CIDR notation (include /32 for a single IPv4 host, /128 for IPv6).
type CreateWhitelistIPRequest struct {
	IPAddress string `json:"ip_address" binding:"required,max=64" example:"203.0.113.10/32"`
	// Description optional human-readable note (e.g. office name, ticket id).
	Description string `json:"description" binding:"max=255" example:"HQ VPN egress - INC-1234"`
	// IsActive when false, the entry is persisted but must not grant access (enforcement is middleware responsibility).
	IsActive *bool `json:"is_active" binding:"required" example:"true"`
}