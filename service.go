package main

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

//go:embed scripts/autostart.plist
var autostartPlist string

//go:embed scripts/service_info.plist
var serviceInfoPlist string

//go:embed scripts/service_workflow.plist
var serviceWorkflowPlist string

// Service management functions
func getBinaryPath() string {
	if executable, err := os.Executable(); err == nil {
		return executable
	}
	return "/usr/local/bin/macostranslate" // fallback
}

func getAutostartPlistPath() string {
	return filepath.Join(os.Getenv("HOME"), "Library", "LaunchAgents", "pl.com.t3.macostranslate.plist")
}

func getServicePath() string {
	return filepath.Join(os.Getenv("HOME"), "Library", "Services", "macostranslate.workflow")
}

func isServiceInstalled() bool {
	autostartExists := fileExists(getAutostartPlistPath())
	serviceExists := fileExists(getServicePath())
	return autostartExists && serviceExists
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func installService() error {
	binaryPath := getBinaryPath()

	// Install autostart service
	if err := installAutostartService(binaryPath); err != nil {
		return fmt.Errorf("failed to install autostart service: %v", err)
	}

	// Install keyboard shortcut service
	if err := installKeyboardService(binaryPath); err != nil {
		return fmt.Errorf("failed to install keyboard service: %v", err)
	}

	return nil
}

func uninstallService() error {
	if err := uninstallAutostartService(); err != nil {
		return fmt.Errorf("failed to uninstall autostart service: %v", err)
	}

	if err := uninstallKeyboardService(); err != nil {
		return fmt.Errorf("failed to uninstall keyboard service: %v", err)
	}

	return nil
}

func installAutostartService(binaryPath string) error {
	// Create LaunchAgents directory if it doesn't exist
	launchAgentsDir := filepath.Join(os.Getenv("HOME"), "Library", "LaunchAgents")
	if err := os.MkdirAll(launchAgentsDir, 0755); err != nil {
		return fmt.Errorf("failed to create LaunchAgents directory: %v", err)
	}

	// Replace placeholder with actual binary path
	content := strings.ReplaceAll(autostartPlist, "{{BINARY_PATH}}", binaryPath)

	// Write the plist file
	plistPath := getAutostartPlistPath()
	if err := os.WriteFile(plistPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write autostart plist: %v", err)
	}

	// Try to load the service
	if err := runCommand("launchctl", "load", plistPath); err != nil {
		// Don't fail if load fails - it might already be loaded
		fmt.Printf("Warning: Could not load autostart service: %v\n", err)
	}

	return nil
}

func uninstallAutostartService() error {
	plistPath := getAutostartPlistPath()

	if fileExists(plistPath) {
		// Try to unload the service first
		if err := runCommand("launchctl", "unload", plistPath); err != nil {
			// Don't fail if unload fails - service might not be loaded
			fmt.Printf("Warning: Could not unload autostart service: %v\n", err)
		}

		// Remove the plist file
		if err := os.Remove(plistPath); err != nil {
			return fmt.Errorf("failed to remove autostart plist: %v", err)
		}
	}

	return nil
}

func installKeyboardService(binaryPath string) error {
	// Create Services directory if it doesn't exist
	servicesDir := filepath.Join(os.Getenv("HOME"), "Library", "Services")
	if err := os.MkdirAll(servicesDir, 0755); err != nil {
		return fmt.Errorf("failed to create Services directory: %v", err)
	}

	// Create the workflow directory structure
	workflowPath := getServicePath()
	contentsPath := filepath.Join(workflowPath, "Contents")
	if err := os.MkdirAll(contentsPath, 0755); err != nil {
		return fmt.Errorf("failed to create workflow Contents directory: %v", err)
	}

	// Write Info.plist
	infoPath := filepath.Join(contentsPath, "Info.plist")
	if err := os.WriteFile(infoPath, []byte(serviceInfoPlist), 0644); err != nil {
		return fmt.Errorf("failed to write service Info.plist: %v", err)
	}

	// Write document.wflow with binary path
	workflowContent := strings.ReplaceAll(serviceWorkflowPlist, "{{BINARY_PATH}}", binaryPath)
	workflowFilePath := filepath.Join(contentsPath, "document.wflow")
	if err := os.WriteFile(workflowFilePath, []byte(workflowContent), 0644); err != nil {
		return fmt.Errorf("failed to write service workflow: %v", err)
	}

	return nil
}

func uninstallKeyboardService() error {
	servicePath := getServicePath()

	if fileExists(servicePath) {
		if err := os.RemoveAll(servicePath); err != nil {
			return fmt.Errorf("failed to remove keyboard service: %v", err)
		}
	}

	return nil
}

func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
