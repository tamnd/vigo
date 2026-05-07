package vio

import "testing"

func TestAttrWithWithout(t *testing.T) {
	a := NewAttr(White, Blue).With(ModBold).With(ModUnderline)
	if a.Mod&ModBold == 0 || a.Mod&ModUnderline == 0 {
		t.Fatalf("with: %+v", a)
	}
	b := a.Without(ModBold)
	if b.Mod&ModBold != 0 {
		t.Fatalf("without: %+v", b)
	}
	if a.Mod&ModBold == 0 {
		t.Fatal("Without mutated original")
	}
}

func TestAttrStyleAllModifiers(_ *testing.T) {
	mods := []Modifier{
		ModBold, ModItalic, ModUnderline, ModReverse, ModBlink, ModStrike,
	}
	for _, m := range mods {
		a := NewAttr(White, Black).With(m)
		// Style() must not panic and must produce something.
		_ = a.Style()
	}
}
