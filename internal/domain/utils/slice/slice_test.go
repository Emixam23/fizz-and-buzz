package slice

import "testing"

func TestIsStringInSlice(t *testing.T) {
	type params struct {
		a    string
		list []string
	}
	tests := []struct {
		name   string
		params params
		want   bool
	}{
		{
			name: "string in slice",
			params: params{
				a:    "test",
				list: []string{"plop", "test"},
			},
			want: true,
		},
		{
			name: "string not in slice",
			params: params{
				a:    "test",
				list: []string{"foo", "bar"},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsStringInSlice(tt.params.a, tt.params.list); got != tt.want {
				t.Errorf("IsStringInSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}
