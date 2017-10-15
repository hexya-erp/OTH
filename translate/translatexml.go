package translate

import (
	"fmt"

	"strings"

	"github.com/beevik/etree"
	"github.com/hexya-erp/hexya/hexya/tools/strutils"
)

//translate an odoo xml into an hexya xml
func TransXML(sourcedoc *etree.Document, pack string) *etree.Document {

	doc := etree.NewDocument()
	hexya := doc.CreateElement("hexya")
	data := hexya.CreateElement("data")

	var recs []*etree.Element
	var path string

	if sourcedoc.SelectElement("odoo").SelectElement("data") != nil {
		path = "odoo/data"
	} else {
		path = "odoo"
	}

	recs = sourcedoc.FindElements(path + "/record")

	for _, rec := range recs {

		recType := rec.SelectAttrValue("model", "")

		switch recType {

		case "ir.ui.view":

			view := data.CreateElement("view")

			view.CreateAttr("id", AbsolutizeExternalID(rec.SelectAttr("id").Value, pack))

			fields := rec.FindElements("field")

			for _, fi := range fields {

				if fi.SelectAttr("name").Value == "model" {

					view.CreateAttr("model", CamelCase(fi.Text()))

				} else if fi.SelectAttr("name").Value == "arch" {

					for _, child := range fi.ChildElements() {

						view.AddChild(child)
					}
				}

			}

		case "ir.actions.act_window":

			viewmodeset := false
			action := data.CreateElement("action")

			action.CreateAttr("id", AbsolutizeExternalID(rec.SelectAttr("id").Value, pack))
			action.CreateAttr("type", "ir.actions.act_window")

			fields := rec.FindElements("field")

			for _, fi := range fields {

				switch fi.SelectAttr("name").Value {

				case "res_model":
					action.CreateAttr("model", CamelCase(fi.Text()))

				case "name":
					action.CreateAttr("name", CamelCase(fi.Text()))

				case "view_mode":
					action.CreateAttr("view_mode", fi.Text())
					viewmodeset = true

				case "search_view_id":
					action.CreateAttr("search_view_id", AbsolutizeExternalID(fi.SelectAttr("ref").Value, pack))

				case "help":
					help := action.CreateElement("help")
					for _, markup := range fi.ChildElements() {
						help.AddChild(markup)
					}

				case "view_id":
					if fi.SelectAttr("ref") != nil {
						action.CreateAttr("view_id", AbsolutizeExternalID(fi.SelectAttr("ref").Value, pack))
					}

				}

			}

			if viewmodeset == false {
				action.CreateAttr("view_mode", "tree,form")
			}

		case "ir.actions.act_window.view":

			var id string
			var mode string
			var action string

			fields := rec.FindElements("field")

			for _, fi := range fields {

				switch fi.SelectAttr("name").Value {

				case "view_mode":
					mode = fi.Text()
				case "view_id":
					id = fi.SelectAttr("ref").Value
				case "act_window_id":
					action = fi.SelectAttr("ref").Value
				}

			}

			for _, d := range data.ChildElements() {

				if d.SelectAttr("id").Value == action {

					v := d.CreateElement("view")
					v.CreateAttr("id", AbsolutizeExternalID(id, pack))
					v.CreateAttr("type", mode)
				}
			}

		}

	}

	actwindows := sourcedoc.FindElements(path + "/act_window")

	for _, act := range actwindows {

		act_window := data.CreateElement("act_window")

		act_window.CreateAttr("id", AbsolutizeExternalID(act.SelectAttrValue("id", ""), pack))
		act_window.CreateAttr("name", act.SelectAttrValue("name", ""))
		act_window.CreateAttr("model", CamelCase(act.SelectAttrValue("res_model", "")))
		act_window.CreateAttr("src_model", act.SelectAttrValue("src_model", ""))
		act_window.CreateAttr("view_mode", act.SelectAttrValue("view_mode", "tree,form"))
		act_window.CreateAttr("target", act.SelectAttrValue("target", ""))

	}

	menuItems := sourcedoc.FindElements(path + "/menuitem")

	for _, act := range menuItems {

		mItem := data.CreateElement("menuitem")

		mItem.CreateAttr("id", AbsolutizeExternalID(act.SelectAttrValue("id", ""), pack))
		name := act.SelectAttrValue("name", "")
		if name != "" {
			mItem.CreateAttr("name", name)
		}
		sequence := act.SelectAttrValue("sequence", "")
		if sequence != "" {
			mItem.CreateAttr("sequence", sequence)
		}
		parent := AbsolutizeExternalID(act.SelectAttrValue("parent", ""), pack)
		if parent != "" {
			mItem.CreateAttr("parent", parent)
		}
		action := AbsolutizeExternalID(act.SelectAttrValue("action", ""), pack)
		if action != "" {
			mItem.CreateAttr("action", action)
		}

		groups := strings.Split(act.SelectAttrValue("groups", ""), ",")
		var absGroups []string
		for _, group := range groups {
			ag := AbsolutizeExternalID(group, pack)
			if ag != "" {
				absGroups = append(absGroups, ag)
			}
		}
		if len(absGroups) > 0 {
			mItem.CreateAttr("groups", strings.Join(groups, ","))
		}
	}

	doc.Indent(4)

	return doc
}

func AbsolutizeExternalID(id, packName string) string {
	if id == "" {
		return ""
	}
	tokens := strings.Split(id, ".")
	if len(tokens) == 1 {
		packPrefix := strutils.SnakeCaseString(packName)
		return fmt.Sprintf("%s_%s", packPrefix, id)
	}
	if len(tokens) == 2 {
		return strings.Replace(id, ".", "_", 1)
	}
	panic(fmt.Errorf("unable to Absolutize external ID: %s", id))
}
