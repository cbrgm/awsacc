package main

import (
	"reflect"
	"testing"
)

func TestMapIDsToName(t *testing.T) {
	type args struct {
		c *Config
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "convert",
			args: args{
				c: &Config{
					Accounts: []Account{
						{
							ID:   "123456789012",
							Name: "production",
						},
						{
							ID:   "223456789012",
							Name: "integration",
						},
						{
							ID:   "323456789012",
							Name: "development",
						},
					},
				},
			},
			want: map[string]string{
				"123456789012": "production",
				"223456789012": "integration",
				"323456789012": "development",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MapIDsToName(tt.args.c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MapIDsToName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMapNameToIDs(t *testing.T) {
	type args struct {
		c *Config
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "convert",
			args: args{
				c: &Config{
					Accounts: []Account{
						{
							ID:   "123456789012",
							Name: "production",
						},
						{
							ID:   "223456789012",
							Name: "integration",
						},
						{
							ID:   "323456789012",
							Name: "development",
						},
					},
				},
			},
			want: map[string]string{
				"production":  "123456789012",
				"integration": "223456789012",
				"development": "323456789012",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MapNameToIDs(tt.args.c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MapNameToIDs() = %v, want %v", got, tt.want)
			}
		})
	}
}
