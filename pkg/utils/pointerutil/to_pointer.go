package pointerutil

// ToString returns a pointer of string
func ToString(s string) *string {
	return &s
}

// ToInt returns a pointer of int
func ToInt(i int) *int {
	return &i
}
