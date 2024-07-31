package main

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

type Option struct {
	Name    string
	Value   string
	Choices []string
	Default string
}

var options = []Option{
	{"Create cluster admin user", "No", []string{"No", "Yes"}, "No"},
	{"Deploy cluster using AWS STS", "Yes", []string{"Yes", "No"}, "Yes"},
	{"OpenShift version", "4.12.46", []string{"4.12.46", "4.14.8", "4.14.7", "4.13.5", "4.12"}, "4.12.46"},
	{"External ID (optional)", "", []string{"BLANK", "Custom typing"}, ""},
	{"Operator roles prefix", "clustername-version1", []string{"clustername-version1"}, "clustername-version1"},
	{"Tags (optional)", "", nil, ""},
	{"Multiple availability zones", "No", []string{"No", "Yes"}, "No"},
	{"AWS region", "us-east-2", []string{"us-east-2", "eu-west-3", "sa-east-1", "us-west-1", "us-west-2"}, "us-east-2"},
	{"PrivateLink cluster", "No", []string{"No", "Yes"}, "No"},
	{"Machine CIDR", "10.0.0.0/16", []string{"10.0.0.0/16", "Custom typing"}, "10.0.0.0/16"},
	{"Service CIDR", "172.30.0.0/16", []string{"172.30.0.0/16", "Custom typing"}, "172.30.0.0/16"},
	{"Pod CIDR", "10.128.0.0/14", []string{"10.128.0.0/14", "Custom typing"}, "10.128.0.0/14"},
	{"Install into an existing VPC", "No", []string{"No", "Yes"}, "No"},
	{"Select availability zones", "No", []string{"No", "Yes"}, "No"},
	{"Enable Customer Managed key", "No", []string{"No", "Yes"}, "No"},
	{"Enable autoscaling", "No", []string{"No", "Yes"}, "No"},
	{"Compute nodes", "2", []string{"2", "Custom typing"}, "2"},
	{"Worker machine pool labels (optional)", "", []string{"BLANK", "Custom typing"}, ""},
	{"Host prefix", "23", []string{"23", "Custom typing"}, "23"},
	{"Machine pool root disk size (GiB or TiB)", "300GiB", []string{"300GiB", "Custom typing"}, "300GiB"},
	{"Enable FIPS support", "No", []string{"No", "Yes"}, "No"},
	{"Encrypt etcd data", "No", []string{"No", "Yes"}, "No"},
	{"Disable Workload monitoring", "No", []string{"No", "Yes"}, "No"},
}

func main() {
	// Ask for cluster name
	clusterName := ""
	prompt := &survey.Input{Message: "Enter the cluster name:"}
	survey.AskOne(prompt, &clusterName)
	fmt.Printf("Cluster name: %s\n", clusterName)

	// Ask if user wants to use default settings
	useDefaultSettings := true
	promptConfirm := &survey.Confirm{
		Message: "Create cluster using default settings?",
		Default: true,
	}
	survey.AskOne(promptConfirm, &useDefaultSettings)

	if !useDefaultSettings {
		for {
			// Show updated list of options
			optionChoices := make([]string, len(options))
			for i, option := range options {
				optionChoices[i] = fmt.Sprintf("%s: %s", option.Name, option.Value)
			}
			optionChoices = append(optionChoices, "---\nFINISH CONFIGURATION")

			chosenOption := ""
			prompt := &survey.Select{
				Message:  "Select an option to configure:",
				Options:  optionChoices,
				PageSize: 30,
			}
			survey.AskOne(prompt, &chosenOption)

			if chosenOption == "---\nFINISH CONFIGURATION" {
				break
			} else {
				optionName := strings.Split(chosenOption, ": ")[0]
				configureOption(optionName)
			}
		}
	}

	confirmExecution(clusterName)
}

func configureOption(optionName string) {
	for i := range options {
		if options[i].Name == optionName {
			if options[i].Choices != nil {
				// Display choices if available
				chosenSetting := ""
				prompt := &survey.Select{
					Message: fmt.Sprintf("\033[1m%s\033[0m (default = '%s'):", options[i].Name, options[i].Default),
					Options: options[i].Choices,
					Default: options[i].Default,
				}
				survey.AskOne(prompt, &chosenSetting)
				options[i].Value = chosenSetting
			} else {
				// Simple input for options with no predefined choices
				prompt := &survey.Input{Message: fmt.Sprintf("\033[1m%s\033[0m:", options[i].Name)}
				survey.AskOne(prompt, &options[i].Value)
			}
			// Re-render the options list silently
			printUpdatedOptions()
			break
		}
	}
}

func printUpdatedOptions() {
	fmt.Print("\033[H\033[2J") // Clear the console
	fmt.Println("Updated Configuration:")
	for _, option := range options {
		fmt.Printf("%s: %s\n", option.Name, option.Value)
	}
	fmt.Println()
}

func confirmExecution(clusterName string) {
	fmt.Println("\nFINAL CONFIGURATION:")
	for _, option := range options {
		// Bold only the answer if it is the default value
		if option.Value == option.Default {
			fmt.Printf("%s: \033[1m%s\033[0m\n", option.Name, option.Value)
		} else {
			fmt.Printf("%s: %s\n", option.Name, option.Value)
		}
	}
	fmt.Println("")

	executeCreation := true
	prompt := &survey.Confirm{
		Message: "Do you want to execute with these settings?",
		Default: true,
	}
	survey.AskOne(prompt, &executeCreation)

	if executeCreation {
		fmt.Printf("I: Executing with '%s'\n", clusterName)
		fmt.Println("I: To execute again in the future, you can run:")
		fmt.Printf("go run classic.go --clusterName %s", clusterName)
		for _, option := range options {
			fmt.Printf(" --%s %s", strings.ReplaceAll(strings.ToLower(option.Name), " ", "-"), option.Value)
		}
		fmt.Println("\nExecuting logic...")
	} else {
		// Show updated list of options
		for {
			optionChoices := make([]string, len(options))
			for i, option := range options {
				optionChoices[i] = fmt.Sprintf("%s: %s", option.Name, option.Value)
			}
			optionChoices = append(optionChoices, "---\nFINISH CONFIGURATION")

			chosenOption := ""
			prompt := &survey.Select{
				Message:  "Select an option to configure:",
				Options:  optionChoices,
				PageSize: 20,
			}
			survey.AskOne(prompt, &chosenOption)

			if chosenOption == "---\nFINISH CONFIGURATION" {
				break
			} else {
				optionName := strings.Split(chosenOption, ": ")[0]
				configureOption(optionName)
			}
		}
		confirmExecution(clusterName)
	}
}
