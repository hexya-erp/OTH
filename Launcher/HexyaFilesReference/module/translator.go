package translator

import (
	"path/filepath"

	"github.com/hexya-erp/OTH/Hextranslate"
	"github.com/hexya-erp/hexya/src/models"
	"github.com/hexya-erp/hexya/src/models/fields"
	"github.com/hexya-erp/hexya/src/models/types"
	"github.com/hexya-erp/pool/h"
	"github.com/hexya-erp/pool/m"
	"github.com/rayman520/GoRayUtils"
)

func GetInputPath(env models.Environment) interface{} {
	content, _ := rayUtils.ReadFileContent("./inputPath")
	return content
}

func GetOutputPath(env models.Environment) interface{} {
	content, _ := rayUtils.ReadFileContent("./outputPath")
	return content
}

func GetGoModulePath(env models.Environment) interface{} {
	content, _ := rayUtils.ReadFileContent("./goModPath")
	return content
}

func GetModelName(env models.Environment) interface{} {
	content, _ := rayUtils.ReadFileContent("./inputPath")
	name := filepath.Base(content)
	return name
}

var fields_Translator = map[string]models.FieldDefinition{
	"InputPath": fields.Char{
		Required: true,
		ReadOnly: true,
		Default:  GetInputPath,
	},
	"OutputPath": fields.Char{
		Required: true,
		Default:  GetOutputPath,
	},
	"ModuleName": fields.Char{
		Default:  GetModelName,
		Required: true,
	},
	"GoModulePath": fields.Char{
		Default:  GetGoModulePath,
		Required: true,
	},
	"LogLvl": fields.Selection{
		String:  "Logger Level",
		Default: models.DefaultValue("1"),
		Selection: types.Selection{
			"0": "1 - Debug",
			"1": "2 - Info",
			"2": "3 - Warning",
			"3": "4 - Error",
			"4": "5 - Fatal",
			"5": "6 - No Log",
		},
		Help: "Every message with a log level inferior than specified won't be displayed"},
	"TranslatePython": fields.Boolean{
		String: "Python to Go",
		Help:   "Attempt to translate raw python to Go code",
		Stored: true,
	},
	"PostGenerate": fields.Boolean{
		String: "try generate",
		Stored: true,
	},
	"PostRun": fields.Boolean{
		String: "try run",
		Stored: true,
	},
}

// StartHextranslate starts HexTranslate
func translator_StartHextranslate(rs m.TranslatorSet) {
	Hextranslate.Init(rs)
}

func init() {
	models.NewModel("Translator")
	h.Translator().AddFields(fields_Translator)
	h.Translator().NewMethod("StartHextranslate", translator_StartHextranslate)
}
