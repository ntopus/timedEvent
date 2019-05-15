package routes

import (
	"encoding/json"
)

type JsendMessage struct {
	status  int16
	message string
	data    interface{}
}

func (j *JsendMessage) Data() interface{} {
	return j.data
}

func (j *JsendMessage) SetData(data interface{}) {
	j.data = data
}

func (j *JsendMessage) Message() string {
	return j.message
}

func (j *JsendMessage) SetMessage(message string) {
	j.message = message
}

func (j *JsendMessage) Status() int16 {
	return j.status
}

func (j *JsendMessage) SetStatus(status int16) {
	j.status = status
}

func (j *JsendMessage) MarshalJSON() ([]byte, error) {
	if j.status >= 200 && j.status <= 299 {
		return json.Marshal(struct {
			Status int16       `json:"Status"`
			Data   interface{} `json:"Data"`
		}{
			Status: j.status,
			Data:   j.data,
		})
	} else {
		return json.Marshal(struct {
			Status  int16  `json:"Status"`
			Message string `json:"Message"`
		}{
			Status:  j.status,
			Message: j.message,
		})
	}
}

func (j *JsendMessage) UnmarshalJSON(data []byte) error {
	aux := struct {
		Status  *int16      `json:"Status"`
		Message *string     `json:"Message"`
		Data    interface{} `json:"Data"`
	}{
		Status:  &j.status,
		Message: &j.message,
		Data:    j.data,
	}
	return json.Unmarshal(data, &aux)
}
