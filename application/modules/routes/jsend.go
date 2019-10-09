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
		statusMsg := "success"
		return json.Marshal(struct {
			Status string      `json:"status"`
			Data   interface{} `json:"data"`
		}{
			Status: statusMsg,
			Data:   j.data,
		})
	} else {
		statusMsg := "fail"
		return json.Marshal(struct {
			Status  string `json:"status"`
			Message string `json:"message"`
		}{
			Status:  statusMsg,
			Message: j.message,
		})
	}
}

func (j *JsendMessage) UnmarshalJSON(data []byte) error {
	aux := struct {
		Status  *string     `json:"status"`
		Message *string     `json:"message"`
		Data    interface{} `json:"data"`
	}{
		Message: &j.message,
		Data:    j.data,
	}
	return json.Unmarshal(data, &aux)
}
