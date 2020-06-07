package Hextranslate

import (
	"fmt"
	"os"
	"strconv"

	"github.com/hexya-erp/pool/m"
)

type GenArgs struct {
	RecordSet m.TranslatorSet

	fast             bool
	help             bool
	inputPath        string
	loadbar          bool
	loggerFileOutput bool
	loggerLevel      uint8
	outputPath       string
	goModPath        string
	moduleName       string

	manifest         Manifest
	pythonFiles      []string
	modelDotFieldIds []string

	FilesData   []string
	FilesDemo   []string
	FilesPython []string
	FilesPo     []string

	ControllerAssets map[string][]string

	ComputeMethodNames []string
}

type ArgTraductorEntry struct {
	output string
	input  []string
}

func (g GenArgs) GetArgTraductor() []ArgTraductorEntry {
	return []ArgTraductorEntry{
		{"b", []string{"-b", "--bar", "--loadbar", "--no-bar", "--no-loadbar", "--hide-loadbar"}},
		{"f", []string{"-f", "--fast"}},
		{"h", []string{"-h", "--help"}},
		{"l", []string{"-l", "--ll", "--log", "--logl", "--log-l", "--loglevel", "--log-level"}},
		{"o", []string{"-o", "--out", "--output"}},
	}
}

func (g GenArgs) GetLoggerLevelTraductor() []ArgTraductorEntry {
	return []ArgTraductorEntry{
		{"0", []string{"0", "debug", "all"}},
		{"1", []string{"1", "info"}},
		{"2", []string{"2", "warn"}},
		{"3", []string{"3", "error"}},
		{"4", []string{"4", "fatal"}},
		{"5", []string{"5", "nothing", "none"}},
	}
}

func (g GenArgs) TraductArg(arg string, argTraductor []ArgTraductorEntry) string {
	for _, entry := range argTraductor {
		for _, trad := range entry.input {
			if arg == trad {
				return entry.output
			}
		}
	}
	return "unknown"
}

func (g GenArgs) TraductLoggerLevel(arg string, logleveltraductor []ArgTraductorEntry) string {
	for _, entry := range logleveltraductor {
		for _, trad := range entry.input {
			if arg == trad {
				return entry.output
			}
		}
	}
	fmt.Printf("Fatal: error while parsing mendatory flag argument '%s': log level value unknown\n", arg)
	os.Exit(0)
	return "1"
}

func NewGenArgs() GenArgs {
	return GenArgs{
		fast:             false,
		help:             false,
		inputPath:        "",
		loadbar:          true,
		loggerFileOutput: false,
		loggerLevel:      1,
		outputPath:       "output",
		moduleName:       "myModule",
	}
}

func ReadArgs(rs m.TranslatorSet) GenArgs {
	genArgs = NewGenArgs()
	genArgs.RecordSet = rs
	genArgs.inputPath = rs.InputPath()
	genArgs.outputPath = rs.OutputPath()
	genArgs.goModPath = rs.GoModulePath()
	tmpint, _ := strconv.Atoi(rs.LogLvl())
	genArgs.loggerLevel = uint8(tmpint)
	genArgs.moduleName = rs.ModuleName()
	genArgs.ControllerAssets = make(map[string][]string)
	return genArgs
}
