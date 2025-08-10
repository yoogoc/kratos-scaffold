package field

type TypeField int

const (
	TypeDouble TypeField = iota
	TypeFloat
	TypeInt32
	TypeInt64
	TypeUint32
	TypeUint64
	TypeBool
	TypeString
	TypeText
	TypeTime
	TypeDate
	TypeJson
)

const DefaultPrimaryFieldType = TypeInt64

var (
	types = map[string]TypeField{
		"float64": TypeDouble,
		"float32": TypeFloat,
		"int32":   TypeInt32,
		"int64":   TypeInt64,
		"uint32":  TypeUint32,
		"uint64":  TypeUint64,
		"bool":    TypeBool,
		"string":  TypeString,
		"text":    TypeText,
		"time":    TypeTime,
		"date":    TypeDate,
		"json":    TypeJson,
	}

	typeNames = [...]string{
		TypeDouble: "float64",
		TypeFloat:  "float32",
		TypeInt32:  "int32",
		TypeInt64:  "int64",
		TypeUint32: "uint32",
		TypeUint64: "uint64",
		TypeBool:   "bool",
		TypeString: "string",
		TypeText:   "string",
		TypeTime:   "time.Time",
		TypeDate:   "time.Time",
		TypeJson:   "*structpb.Struct",
	}

	typeEntSchemaNames = [...]string{
		TypeDouble: "Float",
		TypeFloat:  "Float32",
		TypeInt32:  "Int32",
		TypeInt64:  "Int64",
		TypeUint32: "Uint32",
		TypeUint64: "Uint64",
		TypeBool:   "Bool",
		TypeString: "String",
		TypeText:   "String",
		TypeTime:   "Time",
		TypeDate:   "Time",
		TypeJson:   "JSON",
	}

	typeMysqlNames = [...]string{
		TypeDouble: "numeric",
		TypeFloat:  "numeric",
		TypeInt32:  "int",
		TypeInt64:  "bigint",
		TypeUint32: "int",
		TypeUint64: "bigint",
		TypeBool:   "tinyint",
		TypeString: "varchar(255)",
		TypeText:   "text",
		TypeTime:   "timestamp",
		TypeDate:   "timestamp",
		TypeJson:   "jsonb",
	}

	typePramNames = [...]string{
		TypeDouble: "*float64",
		TypeFloat:  "*float32",
		TypeInt32:  "*int32",
		TypeInt64:  "*int64",
		TypeUint32: "*uint32",
		TypeUint64: "*uint64",
		TypeBool:   "*bool",
		TypeString: "*string",
		TypeTime:   "*time.Time",
		TypeDate:   "*time.Time",
	}

	MaybeGoPackages = []string{
		"time",
	}

	typeProtoNames = [...]string{
		TypeDouble: "double",
		TypeFloat:  "float",
		TypeInt32:  "int32",
		TypeInt64:  "int64",
		TypeUint32: "uint32",
		TypeUint64: "uint64",
		TypeBool:   "bool",
		TypeString: "string",
		TypeText:   "string",
		TypeTime:   "string",
		TypeDate:   "string",
		TypeJson:   "google.protobuf.Struct",
	}

	typeParamProtoNames = [...]string{
		TypeDouble: "optional double",
		TypeFloat:  "optional float",
		TypeInt32:  "optional int32",
		TypeInt64:  "optional int64",
		TypeUint32: "optional uint32",
		TypeUint64: "optional uint64",
		TypeBool:   "optional bool",
		TypeString: "optional string",
		TypeText:   "optional string",
		TypeTime:   "optional string",
		TypeDate:   "optional string",
	}

	typeBiz2Proto = [...]string{
		TypeTime: "%s.Format(time.DateTime)",
		TypeDate: "%s.Format(time.DateOnly)",
	}
)

func (t TypeField) String() string {
	return typeNames[t]
}

func (t TypeField) StringByType(http bool) string {
	if http && (t == TypeTime || t == TypeDate) {
		return "string"
	}
	return typeNames[t]
}

func (t TypeField) IsTime() bool {
	return t == TypeTime
}

func (t TypeField) IsJson() bool {
	return t == TypeJson
}

func (t TypeField) StringParam() string {
	return typePramNames[t]
}

func (t TypeField) StringParamByType(http bool) string {
	if http && (t == TypeTime || t == TypeDate) {
		return "*string"
	}
	return typePramNames[t]
}

func (t TypeField) StringEnt() string {
	return typeEntSchemaNames[t]
}

func (t TypeField) StringMysql() string {
	return typeMysqlNames[t]
}

func (t TypeField) StringProto() string {
	return typeProtoNames[t]
}

func (t TypeField) StringProtoParam() string {
	return typeParamProtoNames[t]
}

func (t TypeField) Biz2Proto() string {
	return typeBiz2Proto[t]
}
