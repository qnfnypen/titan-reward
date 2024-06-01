package myerror

// GetMsg 获取错误信息
func GetMsg(code ErrCode, language ...string) string {
	lan := "en"
	if len(language) > 0 && language[0] != "" {
		lan = language[0]
	}

	emaps, ok := lanErrMaps[lan]
	if !ok {
		return "don't support this language"
	}
	if value, ok := emaps[code]; ok {
		return value
	}

	return "get error message error"
}
