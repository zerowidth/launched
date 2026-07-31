package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLabel(t *testing.T) {
	for _, tc := range []struct {
		name  string
		label string
	}{
		{"backup", "launched.backup"},
		{"My Backup Job", "launched.my_backup_job"},
		{"Some\t whitespace", "launched.some_whitespace"},
		{"keeps.dots-and_underscores", "launched.keeps.dots-and_underscores"},
		// shell metacharacters must not survive into the install script
		{"x`id`", "launched.x_id"},
		{"x$(id)", "launched.x_id"},
		{"x; rm -rf /", "launched.x_rm_-rf"},
		{"x | tee /tmp/x", "launched.x_tee_tmp_x"},
		{"x`{touch,/tmp/pwned}`", "launched.x_touch_tmp_pwned"},
		{"../../etc/passwd", "launched..._.._etc_passwd"},
		{"x\ny", "launched.x_y"},
		{"", "launched.plist"},
		{"!!!", "launched.plist"},
	} {
		assert.Equal(t, tc.label, LaunchdPlist{Name: tc.name}.Label(), "name: %q", tc.name)
	}
}
