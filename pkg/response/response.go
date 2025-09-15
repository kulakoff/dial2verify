package response

type PhoneResponse struct {
	Found  bool   `json:"found"`
	Phone  string `json:"phone"`
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

type PhoneResponseErr struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

func SuccessPhoneCheck(phone string, found bool) *PhoneResponse {
	return &PhoneResponse{
		Found:  found,
		Phone:  phone,
		Status: "success",
	}
}

func Error(msg string) *PhoneResponseErr {
	return &PhoneResponseErr{
		Status:  "error",
		Message: msg,
	}
}
