package postgres

import "strings"

// 'Some text” not valid, so I duplicate ' :')
func PrepareString(str string) string {
	return strings.ReplaceAll(str, "'", "''")
}
