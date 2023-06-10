package reflection

import (
	"reflect"
	"testing"
)

func TestGetTypesImplementInterface(t *testing.T) {
	expected := map[string][]reflect.Type{
		"Group1": {reflect.TypeOf(&Test{})},
	}

	result := GetImplementInterfaceTypes[ITest]()

	if len(result) != len(expected) {
		t.Errorf("Unexpected number of groups. Expected: %d, Got: %d", len(expected), len(result))
	}

	for groupName, expectedTypes := range expected {
		resultTypes, ok := result[groupName]
		if !ok {
			t.Errorf("Group %s not found in the result map", groupName)
			continue
		}

		if len(resultTypes) != len(expectedTypes) {
			t.Errorf("Unexpected number of types in group %s. Expected: %d, Got: %d", groupName, len(expectedTypes), len(resultTypes))
			continue
		}

		for i, expectedType := range expectedTypes {
			resultType := resultTypes[i]
			if resultType != expectedType {
				t.Errorf("Type mismatch in group %s at index %d. Expected: %s, Got: %s", groupName, i, expectedType.Name(), resultType.Name())
			}
		}
	}
}

type Test struct {
	A int
}

type ITest interface {
	Method1()
}

func (t *Test) Method1() {
}
