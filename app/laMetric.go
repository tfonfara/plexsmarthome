package app

type LaMetricRequestData struct {
	Priority string        `json:"priority"`
	IconType string        `json:"icon_type"`
	Model    LaMetricModel `json:"model"`
}

type LaMetricModel struct {
	Frames []LaMetricFrame `json:"frames"`
	Cycles int             `json:"cycles"`
}

type LaMetricFrame struct {
	Text string `json:"text"`
	Icon int    `json:"icon"`
}
