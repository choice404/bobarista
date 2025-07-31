package bobarista

import (
	"fmt"
	"os"
	"path/filepath"
)

// Logger defines the interface for logging operations in Bobarista.
// Implementations can provide custom logging behavior for different environments.
type Logger interface {
	// LogError logs an error message.
	LogError(err error)

	// LogInfo logs an informational message.
	LogInfo(message string)

	// LogDebug logs a debug message.
	LogDebug(message string)

	// LogWarning logs a warning message.
	LogWarning(message string)
}

// LogFilename specifies the name of the log file to use.
// If empty, a default filename with the process ID will be generated.
var LogFilename string

// openLogFile opens or creates the log file for writing.
// It creates the necessary directory structure in the user's config directory.
// Returns an open file handle or an error if the operation fails.
func openLogFile() (*os.File, error) {
	// Generate default filename if not specified
	if LogFilename == "" {
		LogFilename = fmt.Sprintf("cupsleeve_%d.log", os.Getpid())
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	logDir := filepath.Join(homeDir, ".config/bobarista")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	logFilePath := filepath.Join(homeDir, ".config/bobarista/", LogFilename)
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}
	return file, nil
}

// LogError logs an error message to the log file.
// If the log file cannot be opened or written to, the error is printed to stderr.
func LogError(err error) {
	file, errOpen := openLogFile()
	if errOpen != nil {
		fmt.Fprintf(os.Stderr, "Failed to open log file: %v\n", errOpen)
		return
	}
	defer file.Close()

	_, errWrite := file.WriteString(fmt.Sprintf("ERROR: %v\n", err))
	if errWrite != nil {
		fmt.Fprintf(os.Stderr, "Failed to write to log file: %v\n", errWrite)
	}
}

// LogInfo logs an informational message to the log file.
// If the log file cannot be opened or written to, the error is printed to stderr.
func LogInfo(message string) {
	file, errOpen := openLogFile()
	if errOpen != nil {
		fmt.Fprintf(os.Stderr, "Failed to open log file: %v\n", errOpen)
		return
	}
	defer file.Close()

	_, errWrite := file.WriteString(fmt.Sprintf("INFO: %s\n", message))
	if errWrite != nil {
		fmt.Fprintf(os.Stderr, "Failed to write to log file: %v\n", errWrite)
	}
}

// LogDebug logs a debug message to the log file.
// If the log file cannot be opened or written to, the error is printed to stderr.
func LogDebug(message string) {
	file, errOpen := openLogFile()
	if errOpen != nil {
		fmt.Fprintf(os.Stderr, "Failed to open log file: %v\n", errOpen)
		return
	}
	defer file.Close()

	_, errWrite := file.WriteString(fmt.Sprintf("DEBUG: %s\n", message))
	if errWrite != nil {
		fmt.Fprintf(os.Stderr, "Failed to write to log file: %v\n", errWrite)
	}
}

// LogWarning logs a warning message to the log file.
// If the log file cannot be opened or written to, the error is printed to stderr.
func LogWarning(message string) {
	file, errOpen := openLogFile()
	if errOpen != nil {
		fmt.Fprintf(os.Stderr, "Failed to open log file: %v\n", errOpen)
		return
	}
	defer file.Close()

	_, errWrite := file.WriteString(fmt.Sprintf("WARNING: %s\n", message))
	if errWrite != nil {
		fmt.Fprintf(os.Stderr, "Failed to write to log file: %v\n", errWrite)
	}
}

// debugLog logs a debug message if debug mode is enabled.
// This is a convenience method for conditional debug logging.
func (f *Bobarista) debugLog(message string) {
	if f.config.Debug {
		LogDebug(message)
	}
}

// infoLog logs an informational message if debug mode is enabled.
// This is a convenience method for conditional info logging.
func (f *Bobarista) infoLog(message string) {
	if f.config.Debug {
		LogInfo(message)
	}
}

// warningLog logs a warning message if debug mode is enabled.
// This is a convenience method for conditional warning logging.
func (f *Bobarista) warningLog(message string) {
	if f.config.Debug {
		LogWarning(message)
	}
}

// errorLog logs an error message if debug mode is enabled.
// This is a convenience method for conditional error logging.
func (f *Bobarista) errorLog(err error) {
	if f.config.Debug {
		LogError(err)
	}
}
