package catalog

import (
	"reflect"
	"testing"
)

func TestCards_Find(t *testing.T) {
	type args struct {
		title string
	}
	tests := []struct {
		name  string
		c     Cards
		args  args
		want  *Card
		want1 bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.c.Find(tt.args.title)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Find() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Find() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
