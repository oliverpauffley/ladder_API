package laddermethods

// From method name as a string return ladder method
func MethodFromName(methodName string) LadderMethod {
	var method LadderMethod
	switch methodName {
	case "elo":
		method = Elo{StartingPoints: 1000}
	default:
		panic("Other ladder methods not yet implemented")
	}
	return method
}
