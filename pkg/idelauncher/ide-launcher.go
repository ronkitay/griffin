package idelauncher

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	config "ronkitay.com/griffin/pkg/configuration"
)

const (
	UNSUPPORTED_LANGUAGE = "Not Supported"

	PYTHON_LANGUAGE  = "python"
	JAVA_LANGUAGE    = "java"
	KOTLIN_LANGUAGE  = "kotlin"
	NODE_JS_LANGUAGE = "node"
	GO_LANGUAGE = "golang"
	RUST_LANGUAGE    = "rust"
)

func OpenInIDE(projectDir string) {
	ideConfiguration := config.LoadConfiguration().UserConfiguration.IdeConfiguration

	if ideConfiguration.DefaultIDE == "" {
		fmt.Println("Missing DefaultIDE configuration")
		os.Exit(12)
	}

	language := detectLanguage(projectDir)
	ide := ideOrDefault(language, ideConfiguration)
	openIDE(ide, projectDir)
}

func OpenInAlternativeIDE(projectDir string) {
	ideConfiguration := config.LoadConfiguration().UserConfiguration.IdeConfiguration

	if ideConfiguration.DefaultIDE == "" {
		fmt.Println("Missing DefaultIDE configuration")
		os.Exit(12)
	}

	language := detectLanguage(projectDir)
	ide := alternativeIdeOrDefault(language, ideConfiguration)
	openIDE(ide, projectDir)
}

func detectLanguage(projectDir string) string {
	if exists(filepath.Join(projectDir, "requirements.txt")) || exists(filepath.Join(projectDir, "Pipfile")) {
		return PYTHON_LANGUAGE
	} else if exists(filepath.Join(projectDir, "build.gradle.kts")) || exists(filepath.Join(projectDir, "settings.gradle.kts")) {
		return KOTLIN_LANGUAGE
	} else if exists(filepath.Join(projectDir, "go.mod")) {
		return GO_LANGUAGE
	} else if exists(filepath.Join(projectDir, "build.gradle")) || exists(filepath.Join(projectDir, "pom.xml")) || exists(filepath.Join(projectDir, "settings.gradle")) {
		return JAVA_LANGUAGE
	} else if exists(filepath.Join(projectDir, "package.json")) {
		return NODE_JS_LANGUAGE
	} else if exists(filepath.Join(projectDir, "Cargo.toml")) {
		return RUST_LANGUAGE
	} else {
		return UNSUPPORTED_LANGUAGE
	}
}

func ideOrDefault(language string, ideConfiguration config.IdeConfiguration) string {
	switch language {
	case PYTHON_LANGUAGE:
		return ifNull(ideConfiguration.Python, ideConfiguration.DefaultIDE)
	case GO_LANGUAGE:
		return ifNull(ideConfiguration.GoLang, ideConfiguration.DefaultIDE)
	case JAVA_LANGUAGE:
		return ifNull(ideConfiguration.Java, ideConfiguration.DefaultIDE)
	case KOTLIN_LANGUAGE:
		return ifNull(ideConfiguration.Kotlin, ideConfiguration.DefaultIDE)
	case NODE_JS_LANGUAGE:
		return ifNull(ideConfiguration.NodeJS, ideConfiguration.DefaultIDE)
	case RUST_LANGUAGE:
		return ifNull(ideConfiguration.Rust, ideConfiguration.DefaultIDE)
	default:
		return ideConfiguration.DefaultIDE
	}
}

func fallbackToDefaultAlternativeIDE(ide string, ideConfiguration config.IdeConfiguration) string {
	return ifNull(ide, ifNull(ideConfiguration.DefaultIDEAlternative, ideConfiguration.DefaultIDE))
}

func alternativeIdeOrDefault(language string, ideConfiguration config.IdeConfiguration) string {
	switch language {
	case PYTHON_LANGUAGE:
		return fallbackToDefaultAlternativeIDE(ideConfiguration.PythonAlternative, ideConfiguration)
	case GO_LANGUAGE:
		return fallbackToDefaultAlternativeIDE(ideConfiguration.GoLangAlternative, ideConfiguration)
	case JAVA_LANGUAGE:
		return fallbackToDefaultAlternativeIDE(ideConfiguration.JavaAlternative, ideConfiguration)
	case KOTLIN_LANGUAGE:
		return fallbackToDefaultAlternativeIDE(ideConfiguration.KotlinAlternative, ideConfiguration)
	case NODE_JS_LANGUAGE:
		return fallbackToDefaultAlternativeIDE(ideConfiguration.NodeJSAlternative, ideConfiguration)
	case RUST_LANGUAGE:
		return fallbackToDefaultAlternativeIDE(ideConfiguration.RustAlternative, ideConfiguration)
	default:
		return fallbackToDefaultAlternativeIDE("", ideConfiguration)
	}
}

func ifNull(nullable string, defaultValue string) string {
	if nullable == "" {
		return defaultValue
	} else {
		return nullable
	}
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func openIDE(ide string, projectDir string) {
	var rootDirectory = projectDir
	if rootDirectory == "." {
		currentDir, err := os.Getwd()
		if err != nil {
			fmt.Println("Cannot resolve working directory:", err)
			os.Exit(1)
		}
		rootDirectory = currentDir
	}
	cmd := exec.Command("open", "-na", ide, "--args", rootDirectory)
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error opening IDE (", ide, "):", err)
	}
}
