package httperror

type WeaverResponse struct {
	Errors []ErrorDetails `json:"errors"`
}

type ErrorDetails struct {
	Code            string `json:"code"`
	Message         string `json:"message"`
	MessageTitle    string `json:"message_title"`
	MessageSeverity string `json:"message_severity"`
}
