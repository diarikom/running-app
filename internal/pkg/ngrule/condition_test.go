package ngrule

import (
	"encoding/json"
	"testing"
)

func renderConditionTest(t *testing.T, input, expected string) {
	c := ConditionArray{}
	err := json.Unmarshal([]byte(input), &c)
	if err != nil {
		t.Errorf("FAIL: unable to parse conditions (error=%s)", err)
		return
	}

	actual, err := c.Render("Var")
	if err != nil {
		t.Errorf("FAIL: unable to render conditions (error=%s)", err)
		return
	}

	if actual != expected {
		t.Errorf("FAIL:\nexpected = %s\nactual   = %s", expected, actual)
		return
	}

	t.Logf("Result: %s", actual)
}

func TestRenderIntCondition(t *testing.T) {
	input := `[{"type":"int","options":{"param":"Count","operator":"gte","ref_value":12000}}]`
	expected := "Var.GetInt(\"Count\") >= 12000"
	renderConditionTest(t, input, expected)
}

func TestRenderFloatCondition(t *testing.T) {
	input := `[{"type":"float","options":{"param":"Percent","operator":"lte","ref_value":2.5}}]`
	expected := "Var.GetFloat(\"Percent\") <= 2.500000"
	renderConditionTest(t, input, expected)
}

func TestRenderStringCondition(t *testing.T) {
	input := `[{"type":"string","options":{"param":"Access","operator":"eq","ref_value":"GRANTED"}}]`
	expected := "Var.GetString(\"Access\") == \"GRANTED\""
	renderConditionTest(t, input, expected)
}

func TestRenderBetweenIntCondition(t *testing.T) {
	input := `[{"type":"between_int","options":{"param_prefix":"period","ref_value":1589303132}}]`
	expected := "(Var.GetInt(\"period_ref\") >= Var.GetInt(\"period_start\") && Var.GetInt(\"period_ref\") <= Var.GetInt(\"period_end\"))"
	renderConditionTest(t, input, expected)
}

func TestRenderMultipleCondition(t *testing.T) {
	input := `[{"type":"between_int","options":{"param_prefix":"period","next_op":"and"}},{"type":"int","options":{"param":"distance","operator":">=","ref_value":10000}}]`
	expected := "(Var.GetInt(\"period_ref\") >= Var.GetInt(\"period_start\") && Var.GetInt(\"period_ref\") <= Var.GetInt(\"period_end\")) && Var.GetInt(\"distance\") >= 10000"
	renderConditionTest(t, input, expected)

	input = `[{"type":"string","options":{"param":"UserId","operator":"eq","ref_value":"12345","next_op":"and"}},{"type":"int","options":{"param":"Level","operator":"gte","ref_value":2,"next_op":"or"}},{"type":"float","options":{"param":"Completion","operator":"lte","ref_value":99.9}}]`
	expected = "Var.GetString(\"UserId\") == \"12345\" && Var.GetInt(\"Level\") >= 2 || Var.GetFloat(\"Completion\") <= 99.900000"
	renderConditionTest(t, input, expected)
}

func TestRenderGroupCondition(t *testing.T) {
	input := `[{"type":"group","options":{"conditions":[{"type":"string","options":{"param":"UserId","operator":"eq","ref_value":"12345","next_op":"and"}},{"type":"int","options":{"param":"Level","operator":"gte","ref_value":2}}],"next_op":"or"}},{"type":"float","options":{"param":"Completion","operator":"lte","ref_value":99.9}}]`
	expected := "(Var.GetString(\"UserId\") == \"12345\" && Var.GetInt(\"Level\") >= 2) || Var.GetFloat(\"Completion\") <= 99.900000"
	renderConditionTest(t, input, expected)
}

func BenchmarkRenderMultipleCondition(t *testing.B) {
	for n := 0; n < t.N; n++ {
		input := `[{"type":"string","options":{"param":"UserId","operator":"eq","ref_value":"12345","next_op":"and"}},{"type":"int","options":{"param":"Level","operator":"gte","ref_value":2,"next_op":"or"}},{"type":"float","options":{"param":"Completion","operator":"lte","ref_value":99.9}}]`
		expected := "Var.GetString(\"UserId\") == \"12345\" && Var.GetInt(\"Level\") >= 2 || Var.GetFloat(\"Completion\") <= 99.900000"

		var c ConditionArray
		err := json.Unmarshal([]byte(input), &c)
		if err != nil {
			t.Errorf("FAIL: unable to parse conditions (error=%s)", err)
			return
		}

		actual, err := c.Render("Var")
		if err != nil {
			t.Errorf("FAIL: unable to render conditions (error=%s)", err)
			return
		}

		if actual != expected {
			t.Errorf("FAIL:\nexpected = %s\nactual= %s", expected, actual)
			return
		}
	}
}
