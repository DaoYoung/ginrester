package gorester

type RestError struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Title   string `json:"title"`
	Details string `json:"details"`
}

func NewRestError(status int, code string, title string, details string) *RestError {
	return &RestError{
		Status:  status,
		Code:    code,
		Title:   title,
		Details: details,
	}
}
func NOContentError(err error) *RestError {
	return NewRestError(400, "not_content", "Not Content", err.Error())
}
func NOChangeError(err error) *RestError {
	return NewRestError(400, "not_change", "Not Change", err.Error())
}
func JsonTypeError(err error) *RestError {
	return NewRestError(400, "json_type", "Json type error", err.Error())
}
func ForbidError(err error) *RestError {
	return NewRestError(400, "forbid", "forbid", err.Error())
}

func NotFoundDaoError(err error) *RestError {
	return NewRestError(400, "not_found", "Not Found", err.Error())
}
func NotExistDaoError(err error) *RestError {
	return NewRestError(400, "not_exist", "Not Exist", err.Error())
}
func QueryDaoError(err error) *RestError {
	return NewRestError(400, "db_query", "DB query error", err.Error())
}
