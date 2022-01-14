package types

import (
	"github.com/gotk3/gotk3/gtk"
	"main/util"
)

// SessionTab contains the GTK components for the 'Session' GTKNotebook page.
type SessionTab struct {
	StatusLabel        *gtk.Label
	ServerLabel        *gtk.Label
	CountryLabel       *gtk.Label
	CityLabel          *gtk.Label
	ServerIPLabel      *gtk.Label
	TechnologyLabel    *gtk.Label
	ProtocolLabel      *gtk.Label
	BytesReceivedLabel *gtk.Label
	BytesSentLabel     *gtk.Label
	UptimeLabel        *gtk.Label
}

// BuildSessionTab constructs the GTKNotebook page for the 'Session' tab from
// the provided builder.
func BuildSessionTab(builder *gtk.Builder) *SessionTab {
	return &SessionTab{
		StatusLabel: util.BuilderGetLabel(builder,
			"session_status_label"),
		ServerLabel: util.BuilderGetLabel(builder,
			"session_server_label"),
		CountryLabel: util.BuilderGetLabel(builder,
			"session_country_label"),
		CityLabel: util.BuilderGetLabel(builder,
			"session_city_label"),
		ServerIPLabel: util.BuilderGetLabel(builder,
			"session_server_ip_label"),
		TechnologyLabel: util.BuilderGetLabel(builder,
			"session_technology_label"),
		ProtocolLabel: util.BuilderGetLabel(builder,
			"session_protocol_label"),
		BytesReceivedLabel: util.BuilderGetLabel(builder,
			"session_bytes_received_label"),
		BytesSentLabel: util.BuilderGetLabel(builder,
			"session_bytes_sent_label"),
		UptimeLabel: util.BuilderGetLabel(builder,
			"session_uptime_label"),
	}
}
