package conf

import "testing"

func TestDurationSecond_UnmarshalText(t *testing.T) {
	tests := []struct {
		name    string
		text    []byte
		durS    DurationSecond
		wantErr bool
	}{
		{"c1", []byte("1m40s"), DurationSecond(100), false},
		{"c2", []byte("20s5h"), DurationSecond(5*3600+20), false},
		{"c3", []byte("168h20s"), DurationSecond(7*24*3600+20), false},
		{"w1", []byte("168h20 s"), DurationSecond(0), true},
		{"NotSupportedDuration0", []byte("1m40d"), DurationSecond(0), true},
		{"NotSupportedDuration1", []byte("3d"), DurationSecond(0), true},
		{"NotSupportedDuration2", []byte("1w2d"), DurationSecond(0), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ds DurationSecond
			err := ds.UnmarshalText(tt.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("DurationSecond.UnmarshalText() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && (ds != tt.durS) {
				t.Errorf("DurationSecond.UnmarshalText() expected = %v, get = %v", tt.durS, ds)
			}
		})
	}
}
