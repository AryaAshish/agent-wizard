package model

import "testing"

func TestParseSkillRef(t *testing.T) {
	r, err := ParseSkillRef("foo")
	if err != nil || r.SourceAlias != "" || r.ID != "foo" {
		t.Fatalf("bare ref: %#v err=%v", r, err)
	}
	r2, err := ParseSkillRef("org/foo-bar")
	if err != nil || r2.SourceAlias != "org" || r2.ID != "foo-bar" {
		t.Fatalf("qualified ref: %#v err=%v", r2, err)
	}
}
