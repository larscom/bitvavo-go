package util

func IfOrElse[T any](condition bool, If func() T, Else T) T {
	if condition {
		return If()
	}
	return Else
}

func GetOrEmpty[T any](key string, data map[string]any) T {
	var empty T
	value, exist := data[key]
	return IfOrElse(exist && value != nil, func() T { return value.(T) }, empty)
}
