package types

import (
	"github.com/gotk3/gotk3/gtk"
	"main/util"
)

// AboutTab contains the GTK components for the 'About' GTKNotebook page.
type AboutTab struct {
	NameLabel        *gtk.Label
	VersionLabel     *gtk.Label
	DescriptionLabel *gtk.Label
	WebsiteLabel     *gtk.Label
	CopyrightLabel   *gtk.Label
	LicenseLabel     *gtk.Label
}

// BuildAboutTab constructs the GTKNotebook page for the 'About' tab from the
// provided builder.
func BuildAboutTab(builder *gtk.Builder) *AboutTab {
	nameLabel := util.BuilderGetLabel(builder, "about_name_label")
	versionLabel := util.BuilderGetLabel(builder, "about_version_label")
	descLabel := util.BuilderGetLabel(builder, "about_description_label")
	websiteLabel := util.BuilderGetLabel(builder, "about_website_label")
	copyLabel := util.BuilderGetLabel(builder, "about_copyright_label")
	licenseLabel := util.BuilderGetLabel(builder, "about_license_label")

	nameLabel.SetMarkup("<span size=\"20pt\"><b>" + AppName + "</b></span>")
	versionLabel.SetMarkup(AppVersion)
	descLabel.SetMarkup(AppDescription)
	websiteLabel.SetMarkup("<a href=\"" + AppWebsite + "\">" + AppName +
		"</a>")
	copyLabel.SetMarkup("\302\251 " + AppCopyright)
	licenseLabel.SetMarkup("This software is distributed under the" +
		" " + AppLicense + ".")

	return &AboutTab{
		NameLabel:        nameLabel,
		VersionLabel:     versionLabel,
		DescriptionLabel: descLabel,
		WebsiteLabel:     websiteLabel,
		CopyrightLabel:   copyLabel,
		LicenseLabel:     licenseLabel,
	}
}
