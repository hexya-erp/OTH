package Hextranslate

import (
	"strings"
)

func RegisterLessHelper(entry GenericEntry) {
	split := strings.Split(entry.Path, "static/")
	path := "/static" + strings.Join(split, "")
	key := "LessHelpers"
	genArgs.ControllerAssets[key] = append(genArgs.ControllerAssets[key], path)
}

func RegisterAsset(entry GenericEntry, Id string) {
	entryType := GetEntryTypeTranslator()[entry.Type]
	split := strings.Split(entry.Path, "static/")
	path := "/static" + strings.Join(split, "")
	key := Id + entryType
	genArgs.ControllerAssets[key] = append(genArgs.ControllerAssets[key], path)
}

func ReadEntry(entry GenericEntry, templateId string) {
	switch {
	case len(templateId) == 0:
		return
	case templateId == "Common", templateId == "Backend", templateId == "Frontend":
		RegisterAsset(entry, templateId)
	case templateId == "LessHelpers":
		RegisterLessHelper(entry)
	}
}

func ScriptToGeneric(entry OdooScriptEntry) GenericEntry {
	return GenericEntry{
		Type: entry.Type,
		Path: entry.Src,
	}
}

func LinkToGeneric(entry OdooLinkEntry) GenericEntry {
	return GenericEntry{
		Type: entry.Type,
		Path: entry.Href,
	}
}

func ReadTemplates(templates []OdooTemplate) {
	templateIdTranslator := GetTemplateIdTranslator()
	for _, template := range templates {
		templateRawId := template.InheritId
		if len(templateRawId) == 0 {
			templateRawId = template.Id
		}
		templateId := templateIdTranslator[templateRawId]
		for _, entry := range template.Script {
			ReadEntry(ScriptToGeneric(entry), templateId)
		}
		for _, entry := range template.Link {
			ReadEntry(LinkToGeneric(entry), templateId)
		}
		for _, xpath := range template.Xpath {
			for _, entry := range xpath.Script {
				ReadEntry(ScriptToGeneric(entry), templateId)
			}
			for _, entry := range xpath.Link {
				ReadEntry(LinkToGeneric(entry), templateId)
			}
		}
	}
}
