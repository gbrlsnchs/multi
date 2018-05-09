package multi

import (
	"encoding/json"
	"log"
)

// Start parses the configuration file and starts the task list.
func Start(f []byte, noColor bool) error {
	var tl TaskList

	if err := json.Unmarshal(f, &tl); err != nil {
		return err
	}

	log.Printf("Initializing task list %s\n", tl.Name)
	log.Printf("Description: %s\n", tl.Desc)

	return tl.Start(noColor)
}
