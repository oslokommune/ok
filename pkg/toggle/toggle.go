package toggle

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	AssumeCDRole bool `json:"assume_cd_role"`
}

func Assume() {
	// Read the JSON file
	data, err := os.ReadFile("_config.auto.tfvars.json")
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	// Unmarshal the JSON data into a map
	config := make(map[string]interface{})
	if err := json.Unmarshal(data, &config); err != nil {
		log.Fatalf("Error unmarshaling JSON: %v", err)
	}

	// Check if 'assume_cd_role' exists and is a boolean, then toggle it
	if assumeCDRole, ok := config["assume_cd_role"].(bool); ok {
		config["assume_cd_role"] = !assumeCDRole
	} else {
		log.Fatal("Error: 'assume_cd_role' not found or is not a boolean")
	}

	// Marshal the map back into JSON
	updatedData, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	// Write the updated JSON back to the file
	if err := os.WriteFile("_config.auto.tfvars.json", updatedData, 0644); err != nil {
		log.Fatalf("Error writing file: %v", err)
	}
}
