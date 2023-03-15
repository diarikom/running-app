package ngrule

import (
	"encoding/json"
	"testing"
)

func renderRuleTest(t *testing.T, input, expected string) {
	var r Rule
	err := json.Unmarshal([]byte(input), &r)
	if err != nil {
		t.Errorf("FAIL: unable to parse rule (error=%s)", err)
		return
	}
	r.VariableName = "Var"

	actual, err := r.Render()
	if err != nil {
		t.Errorf("FAIL: unable to render rule (error=%s)", err)
		return
	}

	if actual != expected {
		t.Errorf("FAIL:\nexpected = %s\nactual   = %s", expected, actual)
		return
	}

	t.Logf("Result: %s", actual)
}

func TestRenderRule(t *testing.T) {
	input := `{"code":"ChallengeLevel1","description":"User reach 10km of accumulated run distance during milestone will get 1 credit","priority":0,"params":[{"name":"challenge_id","type":"string"}],"conditions":[{"type":"between_int","options":{"param_prefix":"period","next_op":"and"}},{"type":"int","options":{"param":"distance","operator":">=","ref_value":10000}}],"actions":[{"type":"add","options":{"target":"credit","value":1}}]}`
	expected := "rule ChallengeLevel1 \"User reach 10km of accumulated run distance during milestone will get 1 credit\" salience 0 {\n  when\n    (Var.GetInt(\"period_ref\") >= Var.GetInt(\"period_start\") && Var.GetInt(\"period_ref\") <= Var.GetInt(\"period_end\")) && Var.GetInt(\"distance\") >= 10000\n  then\n    Var.Add(\"credit\", 1);\n    Retract(\"ChallengeLevel1\");\n}"
	renderRuleTest(t, input, expected)
}
