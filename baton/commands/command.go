package commands

type Command struct {
	Module string      `json:"module"`
	Call   string      `json:"call"`
	ID     string      `json:"id"`
	Body   interface{} `json:"body"`
}
