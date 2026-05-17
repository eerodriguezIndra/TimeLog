package autostart

const appName = "TimeLog"

// Enable activa el autoarranque al iniciar sesión apuntando al ejecutable actual.
func Enable() error { return enable() }

// Disable elimina el autoarranque para el usuario actual.
func Disable() error { return disable() }

// IsEnabled indica si el autoarranque está actualmente configurado.
func IsEnabled() bool { return isEnabled() }
