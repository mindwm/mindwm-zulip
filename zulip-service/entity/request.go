package entity

type CloudEventData struct {
	IODocument struct {
		PS1     string `json:"ps1"`
		Input   string `json:"input"`
		Output  string `json:"output"`
		FullCMD string `json:"full_cmd"`
	} `json:"iodocument"`
}
