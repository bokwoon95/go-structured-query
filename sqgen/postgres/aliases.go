package postgres

// sq Field Types
const (
	FieldTypeBoolean = "sq.BooleanField"
	FieldTypeJSON    = "sq.JSONField"
	FieldTypeNumber  = "sq.NumberField"
	FieldTypeString  = "sq.StringField"
	FieldTypeTime    = "sq.TimeField"
	FieldTypeEnum    = "sq.EnumField"
	FieldTypeArray   = "sq.ArrayField"
	FieldTypeBinary  = "sq.BinaryField"
	FieldTypeUUID    = "sq.UUIDField"

	FieldConstructorBoolean = "sq.NewBooleanField"
	FieldConstructorJSON    = "sq.NewJSONField"
	FieldConstructorNumber  = "sq.NewNumberField"
	FieldConstructorString  = "sq.NewStringField"
	FieldConstructorTime    = "sq.NewTimeField"
	FieldConstructorEnum    = "sq.NewEnumField"
	FieldConstructorArray   = "sq.NewArrayField"
	FieldConstructorBinary  = "sq.NewBinaryField"
	FieldConstructorUUID    = "sq.NewUUIDField"
)

// Go Types
const (
	GoTypeInterface    = "interface{}"
	GoTypeBool         = "bool"
	GoTypeInt          = "int"
	GoTypeFloat64      = "float64"
	GoTypeString       = "string"
	GoTypeTime         = "time.Time"
	GoTypeBoolSlice    = "[]bool"
	GoTypeIntSlice     = "[]int"
	GoTypeFloat64Slice = "[]float64"
	GoTypeStringSlice  = "[]string"
	GoTypeTimeSlice    = "[]time.Time"
	GoTypeByteSlice    = "[]byte"
)
