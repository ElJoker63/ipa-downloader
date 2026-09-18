//go:build !windows

package update

import "os/exec"

func setDetachFlags(cmd *exec.Cmd) {
}
