package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/manifoldco/promptui"
)

// Define ANSI color codes for custom colors
const (
	ColorGreen = "\033[32m"
	ColorRed   = "\033[31m"
	ColorReset = "\033[0m"
)

// Simulate shell commands (for demonstration purposes only)
func runCommand(command string) (string, error) {
	time.Sleep(1 * time.Second) // Simulate delay

	switch command {
	case "create-vpc":
		return "vpc-123456", nil
	case "create-subnet":
		return "subnet-123456", nil
	case "create-nat-gateway":
		return "", fmt.Errorf("You do not have the proper permissions") // Simulate failure
	default:
		return "success", nil
	}
}

// Loading icon function
func loadingIcon(message string, success bool) error {
	icons := []string{"◐", "◓", "◑", "◒"}
	index := 0
	ticker := time.NewTicker(250 * time.Millisecond)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				ticker.Stop()
				return
			case <-ticker.C:
				fmt.Printf("\r%s %s", icons[index], message)
				index = (index + 1) % len(icons)
			}
		}
	}()

	// Simulate process completion or failure
	time.Sleep(1 * time.Second) // Simulate delay
	done <- true                // Stop the loading icon

	if success {
		fmt.Printf("\r%s✓ %s%s\n", ColorGreen, message, ColorReset) // Entire line in green for success
		return nil
	} else {
		// Ensure the loading icon is cleared
		fmt.Print("\r") // Move the cursor to the beginning of the line

		// Print the error message in red
		fmt.Printf("%s✘ %s%s\n", ColorRed, message, ColorReset) // Error message with red cross
		return fmt.Errorf("%s failed", message)
	}
}

func displayHelpText() {
	fmt.Println("Help text:")
	fmt.Println("Provide a unique name for your cluster.")
	fmt.Println("This name will be used to create a specific web address for your cluster.")
}

// Setup cluster simulation
func setupCluster() {
	prompt := promptui.Select{
		Label: "Select an action",
		Items: []string{
			"Create Classic Cluster",
			"Create HCP Cluster",
			"Create User Roles",
			"How to setup ROSA (prerequisites)",
			"Exit",
		},
	}

	_, result, err := prompt.Run()
	if err != nil {
		fmt.Println("Prompt failed:", err)
		return
	}

	switch result {
	case "Create Classic Cluster":
		createClassicCluster()
	case "Create HCP Cluster":
		setupClusterInfrastructure()
	case "Create User Roles":
		manageRoles()
	case "How to setup ROSA (prerequisites)":
		showROSAInfo()
	case "Exit":
		fmt.Println("Exiting setup.")
	default:
		fmt.Println("Invalid choice.")
	}
}

func createClassicCluster() {
	fmt.Println("Running classic.go...")

	// Define the command to run classic.go
	cmd := exec.Command("go", "run", "classic.go")

	// Set the standard input, output, and error to the current process's IO
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command and wait for it to finish
	err := cmd.Run()
	if err != nil {
		fmt.Printf("%sError running classic.go: %s%s\n", ColorRed, err, ColorReset)
		return
	}

	// Optionally print success message or handle further actions
	fmt.Printf("%sclassic.go executed successfully.%s\n", ColorGreen, ColorReset)
	fmt.Println() // Line break
}

func setupClusterInfrastructure() {
	prompt := promptui.Prompt{
		Label: "Enter the cluster name",
	}

	clusterName, err := prompt.Run()
	if err != nil {
		fmt.Println("Prompt failed:", err)
		return
	}

	if clusterName == "" {
		fmt.Println("Cluster name cannot be empty.")
		return
	}

	// Simulate various setup steps
	if err := loadingIcon("Creating VPC...", true); err == nil {
		vpcId, _ := runCommand("create-vpc")
		fmt.Printf("%sVPC created with ID: %s%s\n", ColorGreen, vpcId, ColorReset)
		fmt.Println() // Line break
	}

	if err := loadingIcon("Creating public subnet...", true); err == nil {
		publicSubnetId, _ := runCommand("create-subnet")
		fmt.Printf("%sPublic subnet created with ID: %s%s\n", ColorGreen, publicSubnetId, ColorReset)
		fmt.Println() // Line break
	}

	if err := loadingIcon("Creating NAT Gateway...", false); err != nil {
		// Ensure the loading icon is cleared
		fmt.Print("\r") // Move the cursor to the beginning of the line

		// Print the error message in red
		fmt.Printf("%sError: %s%s\n", ColorRed, err, ColorReset)
		fmt.Printf("%sReason: You do not have the proper permissions%s\n", ColorRed, ColorReset)
		fmt.Printf("%sSolution:%s\n", ColorRed, ColorReset)
		fmt.Printf("%s- Contact your admin to validate your permissions%s\n", ColorRed, ColorReset)
		fmt.Printf("%s- Try exiting your IDE and try again%s\n", ColorRed, ColorReset)
		fmt.Printf("%s- Honestly, you might just be SOL%s\n", ColorRed, ColorReset)
	}

	// Line break before the completion message
	fmt.Println()
	fmt.Printf("%s✓ Setup complete.%s\n", ColorGreen, ColorReset) // Success message with green checkmark
	fmt.Printf("To recreate this cluster in the future, run the following command:\n")
	fmt.Printf("setupCluster --clusterName %s\n", clusterName)
	fmt.Println() // Line break

}

func manageRoles() {
	prompt := promptui.Select{
		Label: "Select role actions",
		Items: []string{
			"Create OCM role",
			"Create user role",
			"Exit",
		},
	}

	_, result, err := prompt.Run()
	if err != nil {
		fmt.Println("Prompt failed:", err)
		return
	}

	switch result {
	case "Create OCM role":
		createOCMRole()
	case "Create user role":
		createUserRole()
	case "Exit":
		fmt.Println("Exiting role management.")
	default:
		fmt.Println("Invalid choice.")
	}
}

func createOCMRole() {
	if err := loadingIcon("Creating OCM role...", true); err == nil {
		fmt.Printf("%sOCM role created successfully.%s\n", ColorGreen, ColorReset)
	}
}

func createUserRole() {
	if err := loadingIcon("Creating user role...", true); err == nil {
		fmt.Printf("%sUser role created successfully.%s\n", ColorGreen, ColorReset)
	}
}

func showROSAInfo() {
	fmt.Println("To setup ROSA (Red Hat OpenShift Service on AWS), ensure you meet the following prerequisites:")
	fmt.Println("1. Have a valid AWS account.")
	fmt.Println("2. Have an OpenShift cluster or subscription.")
	fmt.Println("3. Ensure you have the necessary IAM roles and permissions.")
	fmt.Println("4. Follow the official documentation for detailed setup instructions.")
}

func main() {
	setupCluster()
}
