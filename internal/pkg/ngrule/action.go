package ngrule

import (
	"encoding/json"
	"sort"
)

var actionTC = map[string]ActionConstructor{
	AssignActionType: func() Action {
		return &FactMapAction{Operator: AssignActionType}
	},
	AddActionType: func() Action {
		return &FactMapAction{Operator: AddActionType}
	},
	SubtractActionType: func() Action {
		return &FactMapAction{Operator: SubtractActionType}
	},
	MultiplyActionType: func() Action {
		return &FactMapAction{Operator: MultiplyActionType}
	},
	DivideActionType: func() Action {
		return &FactMapAction{Operator: DivideActionType}
	},
}

type Action interface {
	GetTargetName() string
	GetValue() interface{}
	Render(varName string) string
	GetTargets() []FactParam
}

type ActionArray struct {
	coll []Action
}

type ActionConstructor func() Action

func (a *ActionArray) Render(varName string) (string, error) {
	// Count condition
	end := len(a.coll) - 1

	// Iterate conditions
	var result string
	for k, v := range a.coll {
		// Render condition rule
		result += v.Render(varName) + ";"

		// If not end of collection, added new line
		if k < end {
			result += "\n    "
		}
	}

	return result, nil
}

func (a *ActionArray) UnmarshalJSON(data []byte) error {
	// Parse as a temporary Condition object placeholder
	var p []struct {
		Type    string          `json:"type"`
		Options json.RawMessage `json:"options"`
	}
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	// PrepareFacts collection
	a.coll = make([]Action, len(p))

	// Iterate Action item array
	for k, v := range p {
		// Parse by type
		actionConstructor, ok := actionTC[v.Type]
		if !ok {
			return NewError(ErrActionUnsupported, "unsupported Action type")
		}

		// Construct item
		item := actionConstructor()

		// Parse json
		err = json.Unmarshal(v.Options, item)
		if err != nil {
			return err
		}

		// Push to a
		a.coll[k] = item
	}

	return nil
}

func (a *ActionArray) AddType(name string, cFn ActionConstructor) {
	actionTC[name] = cFn
}

func (a *ActionArray) RemoveType(name string) {
	// If type is available, delete
	_, ok := actionTC[name]
	if ok {
		delete(actionTC, name)
	}
}

func (a *ActionArray) GetTargets() []FactParam {
	// Get parameters from each condition
	params := make([]FactParam, 0)
	for _, v := range a.coll {
		// Get item parameters
		itemParams := v.GetTargets()
		// Merge parameters
		params = MergeParams(params, itemParams)
	}
	return params
}

func (a *ActionArray) FindByTargetName(targetName string) Action {
	arr := ActionSorter(a.coll)

	// Get index
	idx := sort.Search(len(arr), func(i int) bool {
		return arr[i].GetTargetName() >= targetName
	})

	if idx == -1 {
		return nil
	}

	return a.coll[idx]
}

type ActionSorter []Action

func (f ActionSorter) Len() int {
	return len(f)
}

func (f ActionSorter) Less(i, j int) bool {
	return f[i].GetTargetName() < f[j].GetTargetName()
}

func (f ActionSorter) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}
