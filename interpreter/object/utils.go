package object

func GetJValueType(arg JValue) string {
	switch arg.(type) {
	case *JNumber:
		return Number
	case *JString:
		return String
	case *JList:
		return List
	case *JFunction:
		return Function
	case *JBuiltInFunction:
		return BuiltInFunction
	}

	return Unknow
}

func CanHashed(arg JValue) bool {
	argType := GetJValueType(arg)
	if argType == String || argType == Number {
		return true
	}

	return false
}
