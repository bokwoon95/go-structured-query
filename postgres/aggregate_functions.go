package sq

// Count represents the COUNT(*) aggregate function.
func Count() NumberField {
	format := "COUNT(*)"
	return NumberField{
		format: &format,
	}
}

// CountOver represents the COUNT(*) OVER window function.
func CountOver(window Window) NumberField {
	format := "COUNT(*) OVER ?"
	return NumberField{
		format: &format,
		values: []interface{}{window},
	}
}

// Sum represents the SUM() aggregate function.
func Sum(field interface{}) NumberField {
	format := "SUM(?)"
	return NumberField{
		format: &format,
		values: []interface{}{field},
	}
}

// SumOver represents the SUM() OVER window function.
func SumOver(field interface{}, window Window) NumberField {
	format := "SUM(?) OVER ?"
	return NumberField{
		format: &format,
		values: []interface{}{field, window},
	}
}

// Avg represents the AVG() aggregate function.
func Avg(field interface{}) NumberField {
	format := "AVG(?)"
	return NumberField{
		format: &format,
		values: []interface{}{field},
	}
}

// AvgOver represents the AVG() OVER window function.
func AvgOver(field interface{}, window Window) NumberField {
	format := "AVG(?) OVER ?"
	return NumberField{
		format: &format,
		values: []interface{}{field, window},
	}
}

// Min represents the MIN() aggregate function.
func Min(field interface{}) NumberField {
	format := "MIN(?)"
	return NumberField{
		format: &format,
		values: []interface{}{field},
	}
}

// MinOver represents the MIN() OVER window function.
func MinOver(field interface{}, window Window) NumberField {
	format := "MIN(?) OVER ?"
	return NumberField{
		format: &format,
		values: []interface{}{field, window},
	}
}

// Max represents the MAX() aggregate function.
func Max(field interface{}) NumberField {
	format := "MAX(?)"
	return NumberField{
		format: &format,
		values: []interface{}{field},
	}
}

// MaxOver represents the MAX() OVER window function.
func MaxOver(field interface{}, window Window) NumberField {
	format := "MAX(?) OVER ?"
	return NumberField{
		format: &format,
		values: []interface{}{field, window},
	}
}
