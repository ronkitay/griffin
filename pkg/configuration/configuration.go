package configuration

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"ronkitay.com/griffin/pkg/terminal"
)

type IdeConfiguration struct {
	DefaultIDE            string `json:"default"`
	DefaultIDEAlternative string `json:"defaultAlternative"`
	GoLang                string `json:"go"`
	GoLangAlternative     string `json:"goAlternative"`
	Java                  string `json:"java"`
	JavaAlternative       string `json:"javaAlternative"`
	Kotlin                string `json:"kotlin"`
	KotlinAlternative     string `json:"kotlinAlternative"`
	Python                string `json:"python"`
	PythonAlternative     string `json:"pythonAlternative"`
	NodeJS                string `json:"node"`
	NodeJSAlternative     string `json:"nodeAlternative"`
}

type UserConfiguration struct {
	RepoRoots        []string         `json:"repoRoots"`
	IdeConfiguration IdeConfiguration `json:"ideConfiguration"`
}

type Configuration struct {
	RepoListLocation    string
	ProjectListLocation string
	UserConfiguration   UserConfiguration
}

type ConfigurationManager struct {
	config     UserConfiguration
	configFile string
}

func NewConfigurationManager() (*ConfigurationManager, error) {
	configDir := os.Getenv("HOME") + "/.config/griffin"
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("error creating config directory: %v", err)
	}

	configFile := configDir + "/config.json"
	var config UserConfiguration

	// Load existing configuration if it exists
	if exists, _ := fileExists(configFile); exists {
		data, err := os.ReadFile(configFile)
		if err != nil {
			return nil, fmt.Errorf("error reading config file: %v", err)
		}

		if err := json.Unmarshal(data, &config); err != nil {
			return nil, fmt.Errorf("error parsing config file: %v", err)
		}
	}

	return &ConfigurationManager{
		config:     config,
		configFile: configFile,
	}, nil
}

func (cm *ConfigurationManager) AddRepoRoot(path string) error {
	// Expand and verify the path
	expandedPath, err := expandPath(path)
	if err != nil {
		return fmt.Errorf("error expanding path variables: %v", err)
	}

	// Check if the expanded directory exists
	if _, err := os.Stat(expandedPath); os.IsNotExist(err) {
		return fmt.Errorf("directory does not exist: %s (expanded from %s)", expandedPath, path)
	}

	// Add repo root if it's not already in the list (store original path with variables)
	for _, root := range cm.config.RepoRoots {
		if root == path {
			return nil 
		}
	}
	
	cm.config.RepoRoots = append(cm.config.RepoRoots, path)
	return nil
}

func (cm *ConfigurationManager) GetRepoRoots() ([]string, error) {
	var expandedPaths []string
	for _, root := range cm.config.RepoRoots {
		expandedPath, err := expandPath(root)
		if err != nil {
			return nil, fmt.Errorf("error expanding path %s: %v", root, err)
		}
		expandedPaths = append(expandedPaths, expandedPath)
	}
	return expandedPaths, nil
}

func expandPath(path string) (string, error) {
	result := path
	// Find all ${VAR} patterns
	for start := strings.Index(result, "${"); start != -1; start = strings.Index(result, "${") {
		end := strings.Index(result[start:], "}")
		if end == -1 {
			return "", fmt.Errorf("unclosed variable reference in path: %s", path)
		}
		end = start + end + 1

		varName := result[start+2 : end-1]
		varValue := os.Getenv(varName)
		if varValue == "" {
			return "", fmt.Errorf("environment variable not set: %s", varName)
		}

		result = result[:start] + varValue + result[end:]
	}
	return result, nil
}

func (cm *ConfigurationManager) SetDefaultIDE(ide string) {
	cm.config.IdeConfiguration.DefaultIDE = ide
}

func (cm *ConfigurationManager) SetGoIDE(ide string) {
	cm.config.IdeConfiguration.GoLang = ide
}

func (cm *ConfigurationManager) SetJavaIDE(ide string) {
	cm.config.IdeConfiguration.Java = ide
}

func (cm *ConfigurationManager) SetKotlinIDE(ide string) {
	cm.config.IdeConfiguration.Kotlin = ide
}

func (cm *ConfigurationManager) SetPythonIDE(ide string) {
	cm.config.IdeConfiguration.Python = ide
}

func (cm *ConfigurationManager) SetNodeIDE(ide string) {
	cm.config.IdeConfiguration.NodeJS = ide
}

func (cm *ConfigurationManager) SetDefaultIDEAlternative(ide string) {
	cm.config.IdeConfiguration.DefaultIDEAlternative = ide
}

func (cm *ConfigurationManager) SetGoIDEAlternative(ide string) {
	cm.config.IdeConfiguration.GoLangAlternative = ide
}

func (cm *ConfigurationManager) SetJavaIDEAlternative(ide string) {
	cm.config.IdeConfiguration.JavaAlternative = ide
}

func (cm *ConfigurationManager) SetKotlinIDEAlternative(ide string) {
	cm.config.IdeConfiguration.KotlinAlternative = ide
}

func (cm *ConfigurationManager) SetPythonIDEAlternative(ide string) {
	cm.config.IdeConfiguration.PythonAlternative = ide
}

func (cm *ConfigurationManager) SetNodeIDEAlternative(ide string) {
	cm.config.IdeConfiguration.NodeJSAlternative = ide
}

func (cm *ConfigurationManager) Save() error {
	data, err := json.MarshalIndent(cm.config, "", "    ")
	if err != nil {
		return fmt.Errorf("error encoding config: %v", err)
	}

	if err := os.WriteFile(cm.configFile, data, 0644); err != nil {
		return fmt.Errorf("error writing config file: %v", err)
	}

	return nil
}

func (cm *ConfigurationManager) GetConfiguration() UserConfiguration {
	return cm.config
}

func fileExists(filename string) (bool, error) {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		// File does not exist
		return false, nil
	} else if err != nil {
		// An error occurred (other than file not existing)
		return false, err
	}

	// File exists
	return true, nil
}

func LoadConfiguration() Configuration {
	configurationDirectory := os.Getenv("HOME") + "/.config/griffin"
	repoListLocation := configurationDirectory + "/repo.list"
	projectListLocation := configurationDirectory + "/project.list"

	var userConfiguration UserConfiguration

	if exists, _ := fileExists(configurationDirectory + "/config.json"); exists == true {
		jsonConfigFile, fileOpenError := os.ReadFile(configurationDirectory + "/config.json")
		if fileOpenError != nil {
			os.Exit(1)
		}

		jsonReadError := json.Unmarshal(jsonConfigFile, &userConfiguration)
		if jsonReadError != nil {
			panic("Error reading JSON Configuration:" + jsonReadError.Error())
		}
	} else {
		userConfiguration = UserConfiguration{}
	}

	return Configuration{RepoListLocation: repoListLocation, ProjectListLocation: projectListLocation, UserConfiguration: userConfiguration}
}

func RegisterFlags() {
	flag.String("add-repo-root", "", "Add a repository root directory")
	flag.String("default-ide", "", "Set default IDE")
	flag.String("go-ide", "", "Set Go IDE")
	flag.String("java-ide", "", "Set Java IDE")
	flag.String("kotlin-ide", "", "Set Kotlin IDE")
	flag.String("python-ide", "", "Set Python IDE")
	flag.String("node-ide", "", "Set NodeJS IDE")
	flag.String("default-ide-alt", "", "Set default IDE alternative")
	flag.String("go-ide-alt", "", "Set Go IDE alternative")
	flag.String("java-ide-alt", "", "Set Java IDE alternative")
	flag.String("kotlin-ide-alt", "", "Set Kotlin IDE alternative")
	flag.String("python-ide-alt", "", "Set Python IDE alternative")
	flag.String("node-ide-alt", "", "Set NodeJS IDE alternative")
}

func HandleConfiguration() error {
	changesMade := false
	flag.Parse()

	// Get values after flags have been parsed
	addRepoRoot := flag.Lookup("add-repo-root").Value.String()
	defaultIDE := flag.Lookup("default-ide").Value.String()
	goIDE := flag.Lookup("go-ide").Value.String()
	javaIDE := flag.Lookup("java-ide").Value.String()
	kotlinIDE := flag.Lookup("kotlin-ide").Value.String()
	pythonIDE := flag.Lookup("python-ide").Value.String()
	nodeIDE := flag.Lookup("node-ide").Value.String()
	defaultIDEAlt := flag.Lookup("default-ide-alt").Value.String()
	goIDEAlt := flag.Lookup("go-ide-alt").Value.String()
	javaIDEAlt := flag.Lookup("java-ide-alt").Value.String()
	kotlinIDEAlt := flag.Lookup("kotlin-ide-alt").Value.String()
	pythonIDEAlt := flag.Lookup("python-ide-alt").Value.String()
	nodeIDEAlt := flag.Lookup("node-ide-alt").Value.String()

	configManager, err := NewConfigurationManager()
	if err != nil {
		return fmt.Errorf("error initializing configuration: %v", err)
	}

	// Update configuration based on provided flags
	if addRepoRoot != "" {
		if err := configManager.AddRepoRoot(addRepoRoot); err != nil {
			return fmt.Errorf("error adding repository root: %v", err)
		}
	}

	// Update IDE configurations
	if defaultIDE != "" {
		configManager.SetDefaultIDE(defaultIDE)
		changesMade = true
	}
	if goIDE != "" {
		configManager.SetGoIDE(goIDE)
		changesMade = true
	}
	if javaIDE != "" {
		configManager.SetJavaIDE(javaIDE)
		changesMade = true
	}
	if kotlinIDE != "" {
		configManager.SetKotlinIDE(kotlinIDE)
		changesMade = true
	}
	if pythonIDE != "" {
		configManager.SetPythonIDE(pythonIDE)
		changesMade = true
	}
	if nodeIDE != "" {
		configManager.SetNodeIDE(nodeIDE)
		changesMade = true
	}

	if defaultIDEAlt != "" {
		configManager.SetDefaultIDEAlternative(defaultIDEAlt)
		changesMade = true
	}
	if goIDEAlt != "" {
		configManager.SetGoIDEAlternative(goIDEAlt)
		changesMade = true
	}
	if javaIDEAlt != "" {
		configManager.SetJavaIDEAlternative(javaIDEAlt)
		changesMade = true
	}
	if kotlinIDEAlt != "" {
		configManager.SetKotlinIDEAlternative(kotlinIDEAlt)
		changesMade = true
	}
	if pythonIDEAlt != "" {
		configManager.SetPythonIDEAlternative(pythonIDEAlt)
		changesMade = true
	}
	if nodeIDEAlt != "" {
		configManager.SetNodeIDEAlternative(nodeIDEAlt)
		changesMade = true
	}

	// Save if any changes were made
	if changesMade {
		if err := configManager.Save(); err != nil {
			return fmt.Errorf("error saving configuration: %v", err)
		}
	}

	// Display current configuration
	config := configManager.GetConfiguration()
	fmt.Printf(terminal.GREEN_COLOR  + terminal.BOLD_COLOR + "\nCurrent Configuration (%s):\n" + terminal.RESET_COLORS, configManager.configFile)
	fmt.Println("  " + terminal.GREEN_COLOR + "Repository Roots" + terminal.RESET_COLORS)
	for _, root := range config.RepoRoots {
		fmt.Printf("    - %s\n", root)
	}
	fmt.Println("  " + terminal.GREEN_COLOR + "IDE Configuration" + terminal.RESET_COLORS)
	fmt.Printf("    Default IDE: %s\n", config.IdeConfiguration.DefaultIDE)
	if config.IdeConfiguration.GoLang != "" { fmt.Printf("    Go IDE: %s\n", config.IdeConfiguration.GoLang) }
	if config.IdeConfiguration.Java != "" { fmt.Printf("    Java IDE: %s\n", config.IdeConfiguration.Java) }
	if config.IdeConfiguration.Kotlin != "" { fmt.Printf("    Kotlin IDE: %s\n", config.IdeConfiguration.Kotlin) }
	if config.IdeConfiguration.Python != "" { fmt.Printf("    Python IDE: %s\n", config.IdeConfiguration.Python) }
	if config.IdeConfiguration.NodeJS != "" { fmt.Printf("    NodeJS IDE: %s\n", config.IdeConfiguration.NodeJS) }

	return nil
}
