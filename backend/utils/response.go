package utils

func SuccessResponse(data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"code":    0,
		"message": "success",
		"data":    data,
	}
}

func ErrorResponse(code int, message string) map[string]interface{} {
	return map[string]interface{}{
		"code":    code,
		"message": message,
		"data":    nil,
	}
}
