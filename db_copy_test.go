package main

import "testing"

func Test_resolvePath(t *testing.T) {
	tests := []struct {
		postgresVersion string
		tool            string
		want            string
	}{
		{
			postgresVersion: "12",
			tool:            "pg_restore",
			want:            "/usr/lib/postgresql/12/bin/pg_restore",
		},
		{
			postgresVersion: "12",
			tool:            "psql",
			want:            "/usr/lib/postgresql/12/bin/psql",
		},
	}

	for _, tt := range tests {
		t.Run(tt.tool, func(t *testing.T) {
			got := resolvePath(binPathPattern, tt.postgresVersion, tt.tool)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_normalizeTableName(t *testing.T) {
	tests := []struct {
		tableName string
		want      string
	}{
		{
			tableName: "standard_table",
			want:      "pm_standard_table_bundle",
		},
		{
			tableName: "marketing-tags",
			want:      "pm_marketing_x002d_tags_bundle",
		},
		{
			tableName: "ab-testing",
			want:      "pm_ab_x002d_testing_bundle",
		},
		{
			tableName: "magnolia-mgnlVersion",
			want:      "pm_mgnlversion_bundle",
		},
		{
			tableName: "magnolia_conf_sec-mgnlVersion",
			want:      "version_bundle",
		},
	}

	for _, tt := range tests {
		t.Run(tt.tableName, func(t *testing.T) {
			got := normalizeTableName(tt.tableName)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
