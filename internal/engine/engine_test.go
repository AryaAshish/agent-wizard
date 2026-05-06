package engine

import "testing"

func TestAddUnique(t *testing.T) {
	got := AddUnique([]string{"a"}, "b")
	if len(got) != 2 {
		t.Fatalf("len = %d, want 2", len(got))
	}
	got = AddUnique(got, "b")
	if len(got) != 2 {
		t.Fatalf("duplicate should not be added")
	}
}

func TestRemoveValue(t *testing.T) {
	got := RemoveValue([]string{"a", "b", "a"}, "a")
	if len(got) != 1 || got[0] != "b" {
		t.Fatalf("got %#v, want [b]", got)
	}
}

func TestValidateICP(t *testing.T) {
	valid := []string{"solo", "team", "enterprise"}
	for _, v := range valid {
		if err := ValidateICP(v); err != nil {
			t.Fatalf("ValidateICP(%s) error = %v", v, err)
		}
	}
	if err := ValidateICP("invalid"); err == nil {
		t.Fatal("expected error for invalid icp")
	}
}
