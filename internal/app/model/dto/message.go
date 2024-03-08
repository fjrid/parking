package dto

import "github.com/labstack/echo/v4"

var (
	Langs = []string{"id", "en"}

	MessageSucessfully  = Message{Messages: []string{"Sucessfully", "Berhasil"}}
	MessageBadRequest   = Message{Messages: []string{"Bad Request", "Bad Request"}}
	MessageUnauthorized = Message{Messages: []string{"Akses Ditolak", "Unauthorized"}}
)

type (
	Message struct {
		Ctx      echo.Context
		Messages []string
	}
)

func (msg Message) Translate(args ...string) (rs map[string]string) {
	res := make(map[string]string)
	for k, v := range msg.Messages {
		if len(Langs) > k {
			res[Langs[k]] = v
		}
	}

	return res
}
