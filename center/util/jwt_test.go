package util

import (
	"fmt"
	"testing"
)

func TestGenerateActivation(t *testing.T) {
	token, _ := GenerateActivation(1, 1)
	fmt.Println(ParseActivation(token))
	fmt.Println(GenerateActivation(1, 1))
	//	fmt.Println(GenerateActivation(1, 1))
	//	fmt.Println(GenerateActivation(1, 1))
	//}
}
