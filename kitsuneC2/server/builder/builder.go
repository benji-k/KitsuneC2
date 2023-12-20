package builder

import (
	"KitsuneC2/lib/utils"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
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

var implantSource string = "../implant"
var libSource string = "../lib"
var goMod string = "../go.mod"

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
	var err error
	goMod, _ = filepath.Abs(goMod)
	libSource, _ = filepath.Abs(libSource)
	implantSource, err = filepath.Abs(implantSource)
	if err != nil {
		return err
	}

	_, err = os.Stat(goMod)
	if err != nil {
		log.Printf("[ERROR] Cannot find go.mod at location %s.", goMod)
		return err
	}
	_, err = os.Stat(libSource)
	if err != nil {
		log.Printf("[ERROR] Cannot find lib/ at location %s.", libSource)
		return err
	}
	_, err = os.Stat(implantSource)
	if err != nil {
		log.Printf("[ERROR] Cannot find implant/ at location %s.", implantSource)
		return err
	}

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
		log.Printf("[ERROR] Could not access temporary source code at %s.", SourceDir)
		return err
	}
	defer os.Chdir(cwd) //make sure we cd back to the original place we came from

	cmd := exec.Command("go", "build", "-o", config.OutputFile, "-ldflags=-s -w", ".")
	cmd.Env = append(os.Environ(), "GOOS="+config.ImplantOs, "GOARCH="+config.ImplantArch)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

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
	os.Mkdir(filepath.Join(tmpFolder, "lib"), 0700)
	err = os.Mkdir(filepath.Join(tmpFolder, "implant"), 0700)
	if err != nil {
		return "", err
	}
	err = utils.CopyFolder(implantSource, filepath.Join(tmpFolder, "implant"))
	if err != nil {
		log.Printf("[ERROR] builder: Could not copy %s to %s", implantSource, filepath.Join(tmpFolder, "implant"))
		return "", err
	}
	err = utils.CopyFolder(libSource, filepath.Join(tmpFolder, "lib"))
	if err != nil {
		log.Printf("[ERROR] builder: Could not copy %s to %s", libSource, filepath.Join(tmpFolder, "lib"))
		return "", err
	}
	err = os.Link(goMod, filepath.Join(tmpFolder, "go.mod"))
	if err != nil {
		log.Printf("[ERROR] builder: Could to copy %s to %s", goMod, filepath.Join(tmpFolder, "go.mod"))
		return "", err
	}
	return tmpFolder, nil
}
