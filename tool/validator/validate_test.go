package validator

import "testing"

type st1 struct {
	Author  string  `validate:"required,gte=1,lte=15"`
	Email   *string `validate:"omitempty,email"`
	Website string  `validate:"omitempty,url"`
}

func TestValidate(t *testing.T) {
	email := "a@m.com"
	notemail := "a@m"
	type args struct {
		data interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"ok", args{st1{"aw", &email, "https://google.com"}}, false},
		{"ok2", args{st1{"dfa", nil, "https://google.com"}}, false},
		{"not ok: author", args{st1{"awfsfsadfsdfasdfasdfasdfasdfasdfa", &email, "https://google.com"}}, true},
		{"not ok: email", args{st1{"awfsfs", &notemail, "https://google.com"}}, true},
		{"not ok: website", args{st1{"awfsfs", &email, "google.com"}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Validate(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
