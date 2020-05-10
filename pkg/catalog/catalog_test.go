package catalog

import "testing"

func Test_librarian_Write(t *testing.T) {
	type fields struct {
		catalog Options
	}
	type args struct {
		cards Cards
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := librarian{
				catalog: tt.fields.catalog,
			}
			if err := l.Write(tt.args.cards); (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
