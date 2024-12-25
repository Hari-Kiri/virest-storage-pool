package poolAutostart

type Request struct {
	Uuid      string `json:"uuid"`
	Autostart bool   `json:"autostart"`
}
