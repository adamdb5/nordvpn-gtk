package util

import "github.com/gotk3/gotk3/gtk"

// BuilderGetEntry is a helper function for retrieving a generic GTK widget
// from the builder and casting to a GTK Entry.
func BuilderGetEntry(builder *gtk.Builder, name string) *gtk.Entry {
	obj, _ := builder.GetObject(name)
	return obj.(*gtk.Entry)
}

// BuilderGetWindow is a helper function for retrieving a generic GTK widget
// from the builder and casting to a GTK Window.
func BuilderGetWindow(builder *gtk.Builder, name string) *gtk.Window {
	obj, _ := builder.GetObject(name)
	return obj.(*gtk.Window)
}

// BuilderGetInfoBar is a helper function for retrieving a generic GTK widget
// from the builder and casting to a GTK InfoBar.
func BuilderGetInfoBar(builder *gtk.Builder, name string) *gtk.InfoBar {
	obj, _ := builder.GetObject(name)
	return obj.(*gtk.InfoBar)
}

// BuilderGetLabel is a helper function for retrieving a generic GTK widget
// from the builder and casting to a GTK Label.
func BuilderGetLabel(builder *gtk.Builder, name string) *gtk.Label {
	obj, _ := builder.GetObject(name)
	return obj.(*gtk.Label)
}

// BuilderGetComboBoxText is a helper function for retrieving a generic GTK
// widget from the builder and casting to a GTK ComboBoxText.
func BuilderGetComboBoxText(builder *gtk.Builder, name string) *gtk.ComboBoxText {
	obj, _ := builder.GetObject(name)
	return obj.(*gtk.ComboBoxText)
}

// BuilderGetButton is a helper function for retrieving a generic GTK widget
// from the builder and casting to a GTK Button.
func BuilderGetButton(builder *gtk.Builder, name string) *gtk.Button {
	obj, _ := builder.GetObject(name)
	return obj.(*gtk.Button)
}

// BuilderGetFrame is a helper function for retrieving a generic GTK widget
// from the builder and casting to a GTK Frame.
func BuilderGetFrame(builder *gtk.Builder, name string) *gtk.Frame {
	obj, _ := builder.GetObject(name)
	return obj.(*gtk.Frame)
}

// BuilderGetSwitch is a helper function for retrieving a generic GTK widget
// from the builder and casting to a GTK Switch.
func BuilderGetSwitch(builder *gtk.Builder, name string) *gtk.Switch {
	obj, _ := builder.GetObject(name)
	return obj.(*gtk.Switch)
}

// BuilderGetListBox is a helper function for retrieving a generic GTK widget
// from the builder and casting to a GTK List Box.
func BuilderGetListBox(builder *gtk.Builder, name string) *gtk.ListBox {
	obj, _ := builder.GetObject(name)
	return obj.(*gtk.ListBox)
}
