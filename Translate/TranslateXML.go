package Translate

import (
	"github.com/beevik/etree"
)

func TransXML(sourcedoc *etree.Document, docname string) {

	doc := etree.NewDocument()
	hexya := doc.CreateElement("hexya")
	data := hexya.CreateElement("data")

	var recs []*etree.Element

	if sourcedoc.SelectElement("odoo").SelectElement("data") != nil {
		recs = sourcedoc.FindElements("odoo/data/record")
	} else {
		recs = sourcedoc.FindElements("odoo/record")
	}

	for _, rec := range recs {

		recType := rec.SelectAttrValue("model", "")

		switch recType {

		case "ir.ui.view":

			view := data.CreateElement("view")

			view.CreateAttr("id", rec.SelectAttr("id").Value)

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

			action := data.CreateElement("action")

			action.CreateAttr("id", rec.SelectAttr("id").Value)
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

				case "search_view_id":
					action.CreateAttr("search_view_id", fi.SelectAttr("ref").Value)

				case "help":
					help := action.CreateElement("help")
					for _, markup := range fi.ChildElements() {
						help.AddChild(markup)
					}

				case "view_id":
					if fi.SelectAttr("ref") != nil {
						action.CreateAttr("view_id", fi.SelectAttr("ref").Value)
					} else {
						action.CreateAttr("view_id", "")
					}

				}

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
					v.CreateAttr("id", id)
					v.CreateAttr("type", mode)
				}
			}

		}

	}

	doc.Indent(4)
	doc.WriteToFile("ResultXML/" + docname)
}
