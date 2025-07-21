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

//go:embed scripts/system_autostart.plist
var systemAutostartPlist string

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

func getSystemLaunchDaemonPath() string {
	return "/Library/LaunchDaemons/pl.com.t3.macostranslate.plist"
}

func getUserLaunchAgentPath() string {
	return filepath.Join(os.Getenv("HOME"), "Library", "LaunchAgents", "pl.com.t3.macostranslate.plist")
}

func getServicePath() string {
	return filepath.Join(os.Getenv("HOME"), "Library", "Services", "macostranslate.workflow")
}

func isServiceInstalled() bool {
	// Check both user and system level installations
	userAutostartExists := fileExists(getUserLaunchAgentPath())
	systemAutostartExists := fileExists(getSystemLaunchDaemonPath())
	serviceExists := fileExists(getServicePath())

	return (userAutostartExists || systemAutostartExists) && serviceExists
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func installServiceWithChoice(preferSystem bool) error {
	binaryPath := getBinaryPath()

	// Install autostart service based on choice
	if preferSystem {
		if err := installSystemAutostartService(binaryPath); err != nil {
			// Fallback to user installation if system fails
			if err := installUserAutostartService(binaryPath); err != nil {
				return fmt.Errorf("failed to install autostart service: %v", err)
			}
		}
	} else {
		if err := installUserAutostartService(binaryPath); err != nil {
			return fmt.Errorf("failed to install autostart service: %v", err)
		}
	}

	// Install keyboard shortcut service (always user-level)
	if err := installKeyboardService(binaryPath); err != nil {
		return fmt.Errorf("failed to install keyboard service: %v", err)
	}

	return nil
}

func uninstallService() error {
	// Try to uninstall both user and system level services
	var errors []string

	if err := uninstallAutostartService(); err != nil {
		errors = append(errors, fmt.Sprintf("autostart service: %v", err))
	}

	if err := uninstallSystemAutostartService(); err != nil {
		errors = append(errors, fmt.Sprintf("system autostart service: %v", err))
	}

	if err := uninstallKeyboardService(); err != nil {
		errors = append(errors, fmt.Sprintf("keyboard service: %v", err))
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to uninstall some services: %s", strings.Join(errors, "; "))
	}

	return nil
}

func installUserAutostartService(binaryPath string) error {
	// Create LaunchAgents directory if it doesn't exist
	launchAgentsDir := filepath.Join(os.Getenv("HOME"), "Library", "LaunchAgents")
	if err := os.MkdirAll(launchAgentsDir, 0755); err != nil {
		return fmt.Errorf("failed to create LaunchAgents directory: %v", err)
	}

	// Replace placeholder with actual binary path
	content := strings.ReplaceAll(autostartPlist, "{{BINARY_PATH}}", binaryPath)

	// Write the plist file
	plistPath := getUserLaunchAgentPath()
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

func installSystemAutostartService(binaryPath string) error {
	// Create system-level plist content with user context
	content := strings.ReplaceAll(systemAutostartPlist, "{{BINARY_PATH}}", binaryPath)
	content = strings.ReplaceAll(content, "{{USER_NAME}}", os.Getenv("USER"))

	// Use sudo to write the system plist
	plistPath := getSystemLaunchDaemonPath()

	// Create temporary file
	tmpFile, err := os.CreateTemp("", "macostranslate-*.plist")
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(content); err != nil {
		tmpFile.Close()
		return fmt.Errorf("failed to write temporary file: %v", err)
	}
	tmpFile.Close()

	// Use AppleScript with admin privileges to copy the file
	script := fmt.Sprintf(`
do shell script "cp '%s' '%s' && chown root:wheel '%s' && chmod 644 '%s'" with administrator privileges
`, tmpFile.Name(), plistPath, plistPath, plistPath)

	cmd := exec.Command("osascript", "-e", script)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install system plist: %v", err)
	}

	// Try to load the service
	loadScript := fmt.Sprintf(`
do shell script "launchctl load '%s'" with administrator privileges
`, plistPath)

	cmd = exec.Command("osascript", "-e", loadScript)
	if err := cmd.Run(); err != nil {
		fmt.Printf("Warning: Could not load system autostart service: %v\n", err)
	}

	return nil
}

func uninstallAutostartService() error {
	// Try to uninstall user-level service first
	userPlistPath := getUserLaunchAgentPath()
	if fileExists(userPlistPath) {
		if err := uninstallUserAutostartService(); err != nil {
			return err
		}
	}
	return nil
}

func uninstallUserAutostartService() error {
	plistPath := getUserLaunchAgentPath()

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

func uninstallSystemAutostartService() error {
	plistPath := getSystemLaunchDaemonPath()

	if fileExists(plistPath) {
		// Try to unload the service first using sudo
		unloadScript := fmt.Sprintf(`
do shell script "launchctl unload '%s'" with administrator privileges
`, plistPath)

		cmd := exec.Command("osascript", "-e", unloadScript)
		if err := cmd.Run(); err != nil {
			fmt.Printf("Warning: Could not unload system autostart service: %v\n", err)
		}

		// Remove the plist file using sudo
		removeScript := fmt.Sprintf(`
do shell script "rm -f '%s'" with administrator privileges
`, plistPath)

		cmd = exec.Command("osascript", "-e", removeScript)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to remove system autostart plist: %v", err)
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
