package commands

type ID string

type Command struct {
	Module string      `json:"module"`
	Call   string      `json:"call"`
	ID     ID          `json:"id"`
	Body   interface{} `json:"body"`
}
