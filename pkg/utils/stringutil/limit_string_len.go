package stringutil

func LimitStringLen(text string, limit int) string {
	if len(text) > limit {
		return text[:limit]
	}
	return text
}
