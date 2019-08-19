package laddermethods

import "log"

// From method name as a string return ladder method
func MethodFromName(methodName string) LadderMethod {
	var method LadderMethod
	switch methodName {
	case "elo":
		method = Elo{StartingPoints: 1000, ScaleFactor: 32}
	default:
		log.Fatalf("trying to user a ladder method that doesn't exist, %v", methodName)
	}
	return method
}
