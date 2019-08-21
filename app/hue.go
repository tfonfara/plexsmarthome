package app

type HueGroup struct {
	Name    string   `json:"name"`
	Lights  []string `json:"lights"`
	Sensors []string `json:"sensors"`
	Type    string   `json:"type"`
	State   struct {
		AllOn  bool `json:"all_on"`
		AllOff bool `json:"all_off"`
	} `json:"state"`
	Recycle bool   `json:"recycle"`
	Class   string `json:"class"`
}

type HueRequestData struct {
	Scene string `json:"scene"`
}
