package err

import "fmt"

type UIError struct {
	Err        error  `json:"error"`
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

//Interface implementation of string method let allow us to use UIError as type error
func (u UIError) Error() string {
	return fmt.Sprintf("Message: %s, Code: %d\n", u.Message, u.StatusCode)
}
