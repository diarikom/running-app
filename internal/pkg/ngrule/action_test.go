package ngrule

import (
	"encoding/json"
	"testing"
)

func renderActionTest(t *testing.T, input, expected string) {
	a := ActionArray{}
	err := json.Unmarshal([]byte(input), &a)
	if err != nil {
		t.Errorf("FAIL: unable to parse action (error=%s)", err)
		return
	}

	actual, err := a.Render("Var")
	if err != nil {
		t.Errorf("FAIL: unable to render action (error=%s)", err)
		return
	}

	if actual != expected {
		t.Errorf("FAIL:\nexpected = %s\nactual   = %s", expected, actual)
		return
	}

	t.Logf("Result: %s", actual)
}

func TestRenderAction(t *testing.T) {
	inputJson := `[{"type":"add","options":{"target":"credit","value":1,"value_type":"int"}}]`
	expected := "Var.Add(\"credit\", 1);"
	renderActionTest(t, inputJson, expected)

	inputJson = `[{"type":"subtract","options":{"target":"credit","value":1,"value_type":"int"}}]`
	expected = "Var.Subtract(\"credit\", 1);"
	renderActionTest(t, inputJson, expected)

	inputJson = `[{"type":"multiply","options":{"target":"credit","value":1,"value_type":"int"}}]`
	expected = "Var.Multiply(\"credit\", 1);"
	renderActionTest(t, inputJson, expected)

	inputJson = `[{"type":"divide","options":{"target":"credit","value":1,"value_type":"int"}}]`
	expected = "Var.Divide(\"credit\", 1);"
	renderActionTest(t, inputJson, expected)
}
