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
	TypeTime
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
		"time":    TypeTime,
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
		TypeTime:   "time.Time",
	}

	typeEntSchemaNames = [...]string{
		TypeDouble: "Float64",
		TypeFloat:  "Float32",
		TypeInt32:  "Int32",
		TypeInt64:  "Int64",
		TypeUint32: "Uint32",
		TypeUint64: "Uint64",
		TypeBool:   "Bool",
		TypeString: "String",
		TypeTime:   "Time",
	}

	typePramNames = [...]string{
		TypeDouble: "*wrapperspb.DoubleValue",
		TypeFloat:  "*wrapperspb.FloatValue",
		TypeInt32:  "*wrapperspb.Int32Value",
		TypeInt64:  "*wrapperspb.Int64Value",
		TypeUint32: "*wrapperspb.UInt32Value",
		TypeUint64: "*wrapperspb.UInt64Value",
		TypeBool:   "*wrapperspb.BoolValue",
		TypeString: "*wrapperspb.StringValue",
		TypeTime:   "time.Time",
	}

	MaybeGoPackages = []string{
		"google.golang.org/protobuf/types/known/timestamppb",
		"google.golang.org/protobuf/types/known/wrapperspb",
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
		TypeTime:   "google.protobuf.Timestamp",
	}

	typeParamProtoNames = [...]string{
		TypeDouble: "google.protobuf.DoubleValue",
		TypeFloat:  "google.protobuf.FloatValue",
		TypeInt32:  "google.protobuf.Int32Value",
		TypeInt64:  "google.protobuf.Int64Value",
		TypeUint32: "google.protobuf.UInt32Value",
		TypeUint64: "google.protobuf.UInt64Value",
		TypeBool:   "google.protobuf.BoolValue",
		TypeString: "google.protobuf.StringValue",
		TypeTime:   "google.protobuf.Timestamp",
	}

	typeProtoPackages = [...]string{
		TypeDouble: "",
		TypeFloat:  "",
		TypeInt32:  "",
		TypeInt64:  "",
		TypeUint32: "",
		TypeUint64: "",
		TypeBool:   "",
		TypeString: "",
		TypeTime:   "google/protobuf/timestamp.proto",
	}

	typeProtoParamPackages = [...]string{
		TypeDouble: "google/protobuf/wrappers.proto",
		TypeFloat:  "google/protobuf/wrappers.proto",
		TypeInt32:  "google/protobuf/wrappers.proto",
		TypeInt64:  "google/protobuf/wrappers.proto",
		TypeUint32: "google/protobuf/wrappers.proto",
		TypeUint64: "google/protobuf/wrappers.proto",
		TypeBool:   "google/protobuf/wrappers.proto",
		TypeString: "google/protobuf/wrappers.proto",
		TypeTime:   "google/protobuf/timestamp.proto",
	}

	typeBiz2Proto = [...]string{
		TypeDouble: "wrapperspb.Double(%s)",
		TypeFloat:  "wrapperspb.Float(%s)",
		TypeInt32:  "wrapperspb.Int32(%s)",
		TypeInt64:  "wrapperspb.Int64(%s)",
		TypeUint32: "wrapperspb.UInt32(%s)",
		TypeUint64: "wrapperspb.UInt64(%s)",
		TypeBool:   "wrapperspb.Bool(%s)",
		TypeString: "wrapperspb.String(%s)",
		TypeTime:   "timestamppb.New(%s)",
	}
	typeProto2Biz = [...]string{
		TypeDouble: "",
		TypeFloat:  "",
		TypeInt32:  "",
		TypeInt64:  "",
		TypeUint32: "",
		TypeUint64: "",
		TypeBool:   "",
		TypeString: "",
		TypeTime:   "%s.AsTime()",
	}
)

func (t TypeField) String() string {
	return typeNames[t]
}

func (t TypeField) IsTime() bool {
	return t == TypeTime
}

func (t TypeField) StringParam() string {
	return typePramNames[t]
}

func (t TypeField) StringEnt() string {
	return typeEntSchemaNames[t]
}

func (t TypeField) StringProto() string {
	return typeProtoNames[t]
}

func (t TypeField) StringProtoParam() string {
	return typeParamProtoNames[t]
}

func (t TypeField) ImportProto() string {
	return typeProtoPackages[t]
}

func (t TypeField) ImportProtoParam() string {
	return typeProtoParamPackages[t]
}

func (t TypeField) NeedBiz2Proto() bool {
	return typeBiz2Proto[t] != ""
}

func (t TypeField) NeedProto2Biz() bool {
	return typeProto2Biz[t] != ""
}

func (t TypeField) Biz2Proto() string {
	return typeBiz2Proto[t]
}

func (t TypeField) Proto2Biz() string {
	return typeProto2Biz[t]
}
