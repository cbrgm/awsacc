package main

import (
	"reflect"
	"testing"
)

func TestLookupIDForName(t *testing.T) {
	type fields struct {
		idsToName  map[string]string
		nameToIDs  map[string]string
		print      printer
		strictMode bool
	}
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "Must find name",
			fields: fields{
				idsToName: map[string]string{
					"123456789012": "production",
					"223456789012": "integration",
					"323456789012": "development",
				},
				nameToIDs: map[string]string{
					"production":  "123456789012",
					"integration": "223456789012",
					"development": "323456789012",
				},
				print:      newPrinter(true, true),
				strictMode: false,
			},
			args: args{
				name: "production",
			},
			want: "123456789012",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AccTool{
				idsToName:  tt.fields.idsToName,
				nameToIDs:  tt.fields.nameToIDs,
				print:      tt.fields.print,
				strictMode: tt.fields.strictMode,
			}
			if got := a.lookupIDForName(tt.args.name); got != tt.want {
				t.Errorf("lookupIDForName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLookupNameForID(t *testing.T) {
	type fields struct {
		idsToName  map[string]string
		nameToIDs  map[string]string
		print      printer
		strictMode bool
	}
	type args struct {
		id string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "Must find name",
			fields: fields{
				idsToName: map[string]string{
					"123456789012": "production",
					"223456789012": "integration",
					"323456789012": "development",
				},
				nameToIDs: map[string]string{
					"production":  "123456789012",
					"integration": "223456789012",
					"development": "323456789012",
				},
				print:      newPrinter(true, true),
				strictMode: false,
			},
			args: args{
				id: "123456789012",
			},
			want: "production",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AccTool{
				idsToName:  tt.fields.idsToName,
				nameToIDs:  tt.fields.nameToIDs,
				print:      tt.fields.print,
				strictMode: tt.fields.strictMode,
			}
			if got := a.lookupNameForID(tt.args.id); got != tt.want {
				t.Errorf("lookupNameForID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFindAllAccountIDs(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "must find id in arn",
			args: args{
				line: "arn:aws:iam::123456789012:role/somebody",
			},
			want: []string{
				"123456789012",
			},
		},
		{
			name: "must find three ids",
			args: args{
				line: "this 123456789012 is 223456789012 a 323456789012 string",
			},
			want: []string{
				"123456789012",
				"223456789012",
				"323456789012",
			},
		},
		{
			name: "must find no id",
			args: args{
				line: "123456789012E",
			},
			want: []string{},
		},
		{
			name: "must find no id",
			args: args{
				line: "1234567890",
			},
			want: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findAllAccountIDs(tt.args.line); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findAllAccountIDs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasAccountID(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "must be true",
			args: args{
				line: "this 123456789012 is 223456789012 a 323456789012 string",
			},
			want: true,
		},
		{
			name: "must be false",
			args: args{
				line: "this is a string",
			},
			want: false,
		},
		{
			name: "must be false",
			args: args{
				line: "this is a string 123456789012E",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasAccountID(tt.args.line); got != tt.want {
				t.Errorf("hasAccountID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsOnlyDigits(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "must be true",
			args: args{
				line: "1234567890",
			},
			want: true,
		},
		{
			name: "must be false",
			args: args{
				line: "1234567890e",
			},
			want: false,
		},
		{
			name: "must be false",
			args: args{
				line: "e1234567890",
			},
			want: false,
		},
		{
			name: "must be false",
			args: args{
				line: "12345e67890",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isOnlyDigits(tt.args.line); got != tt.want {
				t.Errorf("isOnlyDigits() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToRegexp(t *testing.T) {
	type args struct {
		pattern string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "must be false",
			args: args{
				pattern: "*foobar*",
			},
			want: ".*foobar.*",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toRegexp(tt.args.pattern); got != tt.want {
				t.Errorf("toRegexp() = %v, want %v", got, tt.want)
			}
		})
	}
}
