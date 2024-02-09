package types

const (
	EnableSetMetadata   = "enable_metadata"
	EnableForceTransfer = "enable_force_transfer"
	EnableBurnFrom      = "enable_burn_from"
	// Allows Authorities of the module to mint any token they wish (including base tokens)
	EnableAuthoritiesSudoMint = "enable_admin_sudo_mint"
)

func IsCapabilityEnabled(enabledCapabilities []string, capability string) bool {
	if len(enabledCapabilities) == 0 {
		return true
	}

	for _, v := range enabledCapabilities {
		if v == capability {
			return true
		}
	}

	return false
}
