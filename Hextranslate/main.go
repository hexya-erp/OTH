package Hextranslate

import (
	"os"

	"github.com/hexya-erp/pool/m"
)

func Init(input m.TranslatorSet) {
	genArgs = ReadArgs(input)
	InitSignal()
	PrintVersion()
	CheckRequestedHelp()
	ReadManifest()
	ReadInit()
	os.RemoveAll(genArgs.outputPath)
	StartTranslator()
	WriteInitFiles()
	PostTranslate()
}
