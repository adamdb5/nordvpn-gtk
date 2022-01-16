package types

import (
	"github.com/adamdb5/opennord/pb"
	"github.com/gotk3/gotk3/gtk"
	"log"
	"main/util"
)

type WhitelistTab struct {
	SubnetListBox      *gtk.ListBox
	SubnetEntry        *gtk.Entry
	SubnetAddButton    *gtk.Button
	SubnetRemoveButton *gtk.Button
	UDPListBox         *gtk.ListBox
	UDPEntry           *gtk.Entry
	UDPAddButton       *gtk.Button
	UDPRemoveButton    *gtk.Button
	TCPListBox         *gtk.ListBox
	TCPEntry           *gtk.Entry
	TCPAddButton       *gtk.Button
	TCPRemoveButton    *gtk.Button
	ApplyButton        *gtk.Button
	SaveButton         *gtk.Button
}

func BuildWhitelistTab(builder *gtk.Builder) *WhitelistTab {
	return &WhitelistTab{
		SubnetListBox: util.BuilderGetListBox(builder,
			"whitelist_subnet_list_box"),
		SubnetEntry: util.BuilderGetEntry(builder,
			"whitelist_subnet_entry"),
		SubnetAddButton: util.BuilderGetButton(builder,
			"whitelist_subnet_add_button"),
		SubnetRemoveButton: util.BuilderGetButton(builder,
			"whitelist_subnet_remove_button"),

		UDPListBox: util.BuilderGetListBox(builder,
			"whitelist_udp_list_box"),
		UDPEntry: util.BuilderGetEntry(builder,
			"whitelist_udp_entry"),
		UDPAddButton: util.BuilderGetButton(builder,
			"whitelist_udp_add_button"),
		UDPRemoveButton: util.BuilderGetButton(builder,
			"whitelist_udp_remove_button"),

		TCPListBox: util.BuilderGetListBox(builder,
			"whitelist_tcp_list_box"),
		TCPEntry: util.BuilderGetEntry(builder,
			"whitelist_tcp_entry"),
		TCPAddButton: util.BuilderGetButton(builder,
			"whitelist_tcp_add_button"),
		TCPRemoveButton: util.BuilderGetButton(builder,
			"whitelist_tcp_remove_button"),

		ApplyButton: util.BuilderGetButton(builder,
			"whitelist_apply_button"),
		SaveButton: util.BuilderGetButton(builder, "whitelist_save_button"),
	}
}

func SubnetAddButtonClicked(app *Application) error {
	subnet, _ := app.Window.WhiteListTab.SubnetEntry.GetText()
	label, _ := gtk.LabelNew(subnet)
	row, _ := gtk.ListBoxRowNew()
	row.SetHAlign(gtk.ALIGN_START)
	row.Add(label)
	app.Window.WhiteListTab.SubnetListBox.Add(row)
	row.ShowAll()
	app.Window.WhiteListTab.SubnetEntry.SetText("")
	app.Config.WhiteList.Subnets = append(app.Config.WhiteList.Subnets, "")
	return nil
}

func SubnetRemoveButtonClicked(app *Application) error {
	row := app.Window.WhiteListTab.SubnetListBox.GetSelectedRow()
	if row != nil {
		app.Window.WhiteListTab.SubnetListBox.Remove(row)
		var newSubnets []string
		for _, subnet := range app.Config.WhiteList.Subnets {
			widget, _ := row.GetChild()
			label, _ := widget.(*gtk.Label)
			text, _ := label.GetText()
			if subnet != text {
				newSubnets = append(newSubnets, subnet)
			}
		}
		app.Config.WhiteList.Subnets = newSubnets
		log.Printf("%d", len(app.Config.WhiteList.Subnets))
	}
	return nil
}

func UDPAddButtonClicked(app *Application) error {
	subnet, _ := app.Window.WhiteListTab.UDPEntry.GetText()
	label, _ := gtk.LabelNew(subnet)
	row, _ := gtk.ListBoxRowNew()
	row.SetHAlign(gtk.ALIGN_START)
	row.Add(label)
	app.Window.WhiteListTab.UDPListBox.Add(row)
	row.ShowAll()
	app.Window.WhiteListTab.UDPEntry.SetText("")
	return nil
}

func UDPRemoveButtonClicked(app *Application) error {
	row := app.Window.WhiteListTab.UDPListBox.GetSelectedRow()
	if row != nil {
		app.Window.WhiteListTab.UDPListBox.Remove(row)
	}
	return nil
}

func TCPAddButtonClicked(app *Application) error {
	subnet, _ := app.Window.WhiteListTab.TCPEntry.GetText()
	label, _ := gtk.LabelNew(subnet)
	row, _ := gtk.ListBoxRowNew()
	row.SetHAlign(gtk.ALIGN_START)
	row.Add(label)
	app.Window.WhiteListTab.TCPListBox.Add(row)
	row.ShowAll()
	app.Window.WhiteListTab.TCPEntry.SetText("")
	return nil
}

func TCPRemoveButtonClicked(app *Application) error {
	row := app.Window.WhiteListTab.TCPListBox.GetSelectedRow()
	if row != nil {
		app.Window.WhiteListTab.TCPListBox.Remove(row)
	}
	return nil
}

func WhitelistApplyButtonClicked(app *Application) error {

	err := app.Client.SetWhitelist(&pb.SetWhitelistRequest{
		Whitelist: &pb.Whitelist{
			Ports: &pb.Ports{
				Udp: nil,
				Tcp: nil,
			},
			Subnets: nil,
		},
	})

	if err != nil {
		util.LogError("Unable to apply whitelist", err)
		infoBar := app.Window.InfoBar
		infoBar.Button.SetLabel("Dismiss")
		infoBar.Button.Connect("clicked", infoBar.HideMessage)
		infoBar.DisplayMessage("Unable to apply whitelist",
			gtk.MESSAGE_ERROR)
		return err
	}

	return nil
}

func WhitelistSaveButtonClicked(app *Application) error {
	return nil
}
