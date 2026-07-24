package analyzer

import "testing"

func TestRank_ExportedTypesAboveUnexportedFuncs(t *testing.T) {
	symbols := []Symbol{
		{Kind: "func", Name: "internalHelper", File: "a.go", Line: 1},
		{Kind: "type", Name: "Config", File: "b.go", Line: 1},
		{Kind: "func", Name: "PublicFunc", File: "c.go", Line: 1},
	}

	ranked := Rank(symbols)

	// Config (type+exported=15) > PublicFunc (func+exported=13) > internalHelper (func=3)
	if ranked[0].Name != "Config" {
		t.Errorf("expected Config first, got %s", ranked[0].Name)
	}
	if ranked[1].Name != "PublicFunc" {
		t.Errorf("expected PublicFunc second, got %s", ranked[1].Name)
	}
	if ranked[2].Name != "internalHelper" {
		t.Errorf("expected internalHelper last, got %s", ranked[2].Name)
	}
}

func TestRank_SameScorePreservesRelativeOrder(t *testing.T) {
	symbols := []Symbol{
		{Kind: "func", Name: "Alpha", File: "a.go", Line: 1},
		{Kind: "func", Name: "Beta", File: "b.go", Line: 1},
		{Kind: "func", Name: "Gamma", File: "c.go", Line: 1},
	}

	ranked := Rank(symbols)

	for i := 0; i < len(ranked)-1; i++ {
		if ranked[i].Name != symbols[i].Name {
			t.Errorf("order changed at index %d: expected %s, got %s", i, symbols[i].Name, ranked[i].Name)
		}
	}
}

func TestRank_Empty(t *testing.T) {
	ranked := Rank(nil)
	if len(ranked) != 0 {
		t.Errorf("expected empty result, got %d items", len(ranked))
	}
}

func TestRank_InterfaceScore(t *testing.T) {
	symbols := []Symbol{
		{Kind: "func", Name: "Run", File: "a.go", Line: 1},
		{Kind: "interface", Name: "Runner", File: "b.go", Line: 1},
	}

	ranked := Rank(symbols)

	// interface+exported=15 > func+exported=13
	if ranked[0].Name != "Runner" {
		t.Errorf("expected Runner first, got %s", ranked[0].Name)
	}
}

func TestRank_UnexportedTypeAboveUnexportedFunc(t *testing.T) {
	symbols := []Symbol{
		{Kind: "func", Name: "helper", File: "a.go", Line: 1},
		{Kind: "type", Name: "config", File: "b.go", Line: 1},
	}

	ranked := Rank(symbols)

	// type(5) > func(3)
	if ranked[0].Name != "config" {
		t.Errorf("expected config first, got %s", ranked[0].Name)
	}
}
