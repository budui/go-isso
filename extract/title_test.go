package extract

import (
	"io"
	"strings"
	"testing"
)

var (
	allHTML = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>isso fake</title>
</head>
<body>
	<div id=isso-thread data-title=isso data-isso-id=/new/ >
</body>
</html>
	`
	withoutDIDHTML = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>isso fake</title>
</head>
<body>
	<div id=isso-thread data-title=isso ></div>
</body>
</html>
	`
	withoutDTHTML = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>isso</title>
</head>
<body>
	<div id=isso-thread></div>
</body>
</html>
	`
	withoutTitleHTML = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
</head>
<body>
	<div id=isso-thread></div>
</body>
</html>
	`
	withoutRootHTML = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
</head>
<body>
</body>
</html>
	`
	InvalidHTML = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>isso</title>
</head>
<body>
	<div id=isso-thread></div>
	<div
	<p></span>
</body>
</html>
	`
)

func TestTitleAndThreadURI(t *testing.T) {
	type args struct {
		body         io.Reader
		defaultTitle string
		defaultURI   string
	}
	tests := []struct {
		name      string
		args      args
		wantTitle string
		wantUri   string
		wantErr   bool
	}{
		{"all", args{strings.NewReader(allHTML), "Untitled", "/"}, "isso", "/new/", false},
		{"withoutDID", args{strings.NewReader(withoutDIDHTML), "Untitled", "/"}, "isso", "/", false},
		{"withoutDT", args{strings.NewReader(withoutDTHTML), "Untitled", "/"}, "isso", "/", false},
		{"withoutTitle", args{strings.NewReader(withoutTitleHTML), "Untitled", "/"}, "Untitled", "/", false},
		{"withoutRoot", args{strings.NewReader(withoutRootHTML), "Untitled", "/"}, "Untitled", "/", true},
		{"Invalid", args{strings.NewReader(InvalidHTML), "Untitled", "/"}, "isso", "/", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTitle, gotUri, err := TitleAndThreadURI(tt.args.body, tt.args.defaultTitle, tt.args.defaultURI)
			if (err != nil) != tt.wantErr {
				t.Errorf("TitleAndThreadURI() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotTitle != tt.wantTitle {
				t.Errorf("TitleAndThreadURI() gotTitle = %v, want %v", gotTitle, tt.wantTitle)
			}
			if gotUri != tt.wantUri {
				t.Errorf("TitleAndThreadURI() gotUri = %v, want %v", gotUri, tt.wantUri)
			}
		})
	}
}
