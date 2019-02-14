package info

// Build information. Populated at build-time.
var (
	Version   string
	Revision  string
	Branch    string
	BuildDate string
	GoVersion string
	SettingsComponentEntryPoint string
	CheckComponentEntryPoint    string
)


type PluginInfo struct {
	// detailed plugin description
	Description string `json:"description,omitempty"`

	// unique ID of installed plugin
	// basically it is slugged URI of plugin repository name e. g. supergiant-request-limits-check
	//
	ID string `json:"id,omitempty"`

	// date/Time the plugin was installed
	// Filled by post-install job
	InstalledAt string `json:"installedAt,omitempty"`

	// name is the name of the plugin.
	Name string `json:"name,omitempty"`

	// service labels
	// Filled by post-install job
	ServiceLabels map[string]string `json:"serviceLabels,omitempty"`

	// name of k8s service which is front of plugin deployment
	// Filled by post-install job
	ServiceName string `json:"serviceName,omitempty"`

	// entry points for web components
	SettingsComponentEntryPoint string `json:"settingsComponentEntryPoint,omitempty"`
	CheckComponentEntryPoint    string `json:"checkComponentEntryPoint,omitempty"`

	// plugin status
	Status string `json:"status,omitempty"`

	// plugin version, major version shall be equal to analyze-core version
	Version string `json:"version,omitempty"`

	Revision string `json:"revision,omitempty"`
	Branch    string `json:"branch,omitempty"`
	BuildDate string `json:"buildDate,omitempty"`
	GoVersion string `json:"goVersion,omitempty"`
}

func Info() PluginInfo {
	return PluginInfo{
		Description:   "Sunsetting plugin shows how it is possible to pack nodes better, and which nodes can be shut down",
		ID:            "analyze-plugin-sunsetting",
		InstalledAt:   "",
		Name:          "Supergiant sunsetting plugin",
		ServiceLabels: nil,
		ServiceName:   "",
		SettingsComponentEntryPoint: SettingsComponentEntryPoint,
		CheckComponentEntryPoint: CheckComponentEntryPoint,
		Status:        "OK",
		Version:       Version,
		Revision:      Revision,
		Branch:        Branch,
		BuildDate:     BuildDate,
		GoVersion:     GoVersion,
	}
}
