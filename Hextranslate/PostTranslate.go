package Hextranslate

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rayman520/GoRayUtils"
)

func RunTranslatedModule() {
	if genArgs.RecordSet.PostRun() {
		cmdGen := exec.Command("dropdb", "PostTranslate")
		HexFilesDir, _ := filepath.Abs("../..")
		cmdGen.Dir = HexFilesDir
		rayUtils.ExecPrintCmd(cmdGen)
		cmdGen = exec.Command("createdb", "PostTranslate")
		HexFilesDir, _ = filepath.Abs("../..")
		cmdGen.Dir = HexFilesDir
		rayUtils.ExecPrintCmd(cmdGen)
		cmdGen = exec.Command("hexya", "updatedb", "-o", "--db-name", "PostTranslate")
		HexFilesDir, _ = filepath.Abs("../..")
		cmdGen.Dir = HexFilesDir
		rayUtils.ExecPrintCmd(cmdGen)

		cmdGen = exec.Command("hexya", "server", "-o", "--db-name", "PostTranslate")
		HexFilesDir, _ = filepath.Abs("../..")
		cmdGen.Dir = HexFilesDir
		rayUtils.ExecPrintCmd(cmdGen)
	}
}

func GenerateTranslatedModule() {
	if genArgs.RecordSet.PostGenerate() {
		content, err := rayUtils.ReadFileContent("../HexyaFilesReference/hexya.toml")
		if err != nil {
			rayUtils.ExitFatal("Could not open file. %v", err)
		}
		gopath := os.Getenv("GOPATH") + "/src/"
		content = strings.Replace(content, "$$MODULEPATH$$", fmt.Sprintf(",\n\t\"%s\"", strings.Replace(genArgs.outputPath, gopath, "", -1)), -1)
		rayUtils.WriteFileContent("../../hexya.toml", content)
		cmdGen := exec.Command("hexya", "generate", "-o")
		HexFilesDir, _ := filepath.Abs("../..")
		cmdGen.Dir = HexFilesDir
		rayUtils.ExecPrintCmd(cmdGen)
		// PO
		if genArgs.FilesPo != nil && !strings.HasSuffix(genArgs.FilesPo[0], ".pot") {
			genArgs.FilesPo = append([]string{""}, genArgs.FilesPo...)
		}
		RunTranslatedModule()
	}
}

func PostTranslate() {
	GenerateTranslatedModule()
}
