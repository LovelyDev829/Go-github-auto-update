package nik

import (
	"io/ioutil"
)

// SetVariableToFile creates a file (if it doesn't exist) and writes the variable value to it.
func SetVariableToFile(filename string, value string) error {
	// Convert the value to a string	
	// valueStr := strconv.Itoa(value)
	valueStr := value
	// Write the value to the file
	err := ioutil.WriteFile(filename, []byte(valueStr), 0644)
	if err != nil {
		return err
	}
	return nil
}

// GetVariableFromFile reads the variable value from the file.
func GetVariableFromFile(filename string) (string) {
	// Read the content of the file
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return ""
	}
	
	// Parse the content as an integer
	// value, err := strconv.Atoi(string(content))
	value := string(content)
	if err != nil {
		return ""
	}
	return value
}