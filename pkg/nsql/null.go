package nsql

import "database/sql"

func NullString(str string) sql.NullString {
	return sql.NullString{
		String: str,
		Valid:  str != "",
	}
}
