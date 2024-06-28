package entity

type ZulipBotState struct {
	ID    string `json:"id"`
	State State  `json:"state"`
}

type State string

const (
	StateZulipOn  = "ZULIP_ON"
	StateZulipOff = "ZULIP_OFF"
)
