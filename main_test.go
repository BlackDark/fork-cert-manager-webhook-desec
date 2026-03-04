package main

import (
	"testing"

	extapi "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

func TestGetSubName(t *testing.T) {
	tests := []struct {
		name       string
		fqdn       string
		domainName string
		wantSub    string
		wantErr    bool
	}{
		{name: "normal subdomain", fqdn: "_acme-challenge.example.com", domainName: "example.com", wantSub: "_acme-challenge"},
		{name: "multi-label subdomain", fqdn: "_acme-challenge.sub.example.com", domainName: "example.com", wantSub: "_acme-challenge.sub"},
		{name: "apex domain equals domain", fqdn: "example.com", domainName: "example.com", wantSub: ""},
		{name: "empty fqdn", fqdn: "", domainName: "example.com", wantErr: true},
		{name: "empty domain", fqdn: "_acme-challenge.example.com", domainName: "", wantErr: true},
		{name: "fqdn not in domain", fqdn: "_acme-challenge.other.com", domainName: "example.com", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getSubName(tt.fqdn, tt.domainName)
			if (err != nil) != tt.wantErr {
				t.Errorf("getSubName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.wantSub {
				t.Errorf("getSubName() = %q, want %q", got, tt.wantSub)
			}
		})
	}
}

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "valid config",
			input:   `{"apiTokenSecretRef":{"name":"my-secret","key":"token"}}`,
			wantErr: false,
		},
		{
			name:    "nil config",
			input:   "",
			wantErr: true,
		},
		{
			name:    "missing name",
			input:   `{"apiTokenSecretRef":{"name":"","key":"token"}}`,
			wantErr: true,
		},
		{
			name:    "missing key",
			input:   `{"apiTokenSecretRef":{"name":"my-secret","key":""}}`,
			wantErr: true,
		},
		{
			name:    "invalid json",
			input:   `not-json`,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cfgJSON *extapi.JSON
			if tt.input != "" {
				cfgJSON = &extapi.JSON{Raw: []byte(tt.input)}
			}
			_, err := loadConfig(cfgJSON)
			if (err != nil) != tt.wantErr {
				t.Errorf("loadConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
