package util

import "log"

// LogInfo logs a message with the prefix '[INFO ]'.
func LogInfo(message string) {
	log.Printf("[INFO ] %s", message)
}

// LogWarning logs a message with the prefix '[WARN ]' and an optional error.
func LogWarning(message string, error error) {
	if error != nil {
		log.Printf("[WARN ] %s: %s", message, error)
	} else {
		log.Printf("[WARN ] %s", message)
	}
}

// LogError logs a message with the prefix '[ERROR]' and an optional error.
func LogError(message string, error error) {
	if error != nil {
		log.Printf("[ERROR] %s: %s", message, error)
	} else {
		log.Printf("[ERROR] %s", message)
	}
}

// LogFatal logs a message with the prefix '[FATAL]' and an optional error. The
// application will terminate after this call.
func LogFatal(message string, error error) {
	if error != nil {
		log.Fatalf("[FATAL] %s: %s", message, error)
	} else {
		log.Fatalf("[FATAL] %s", message)
	}
}
