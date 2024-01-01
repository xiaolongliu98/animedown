package util

import (
	"testing"
)

func TestRemoveHTMLElem(t *testing.T) {
	type args struct {
		content string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "T1",
			args: args{
				content: `<td nowrap="nowrap" align="center">>CONTENT-A span class="btl_1">CONTENT-B</span></td>CONTENT-C`,
			},
			want: `>CONTENT-A span class="btl_1">CONTENT-BCONTENT-C`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveHTMLElem(tt.args.content); got != tt.want {
				t.Errorf("RemoveHTMLElem() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetTerminalSize(t *testing.T) {
	h, w, err := GetTerminalSize()
	if err != nil {
		t.Error(err)
	}
	t.Logf("h: %d, w: %d", h, w)
}
