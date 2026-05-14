package skills

import (
	"reflect"
	"testing"
)

func TestFilterSkillsByIDSubstring(t *testing.T) {
	sk := []Skill{{ID: "pr-review"}, {ID: "plan-review"}, {ID: "launch-ready"}}

	tests := []struct {
		name   string
		needle string
		want   []string
	}{
		{name: "empty needle passthrough", needle: "", want: []string{"pr-review", "plan-review", "launch-ready"}},
		{name: "whitespace needle passthrough", needle: "   ", want: []string{"pr-review", "plan-review", "launch-ready"}},
		{name: "substring case insensitive", needle: "PR", want: []string{"pr-review"}},
		{name: "multiple matches", needle: "review", want: []string{"pr-review", "plan-review"}},
		{name: "no matches", needle: "zzz", want: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FilterSkillsByIDSubstring(sk, tt.needle)
			var ids []string
			for _, s := range got {
				ids = append(ids, s.ID)
			}
			if !reflect.DeepEqual(ids, tt.want) {
				t.Fatalf("got %#v want %#v", ids, tt.want)
			}
		})
	}
}
