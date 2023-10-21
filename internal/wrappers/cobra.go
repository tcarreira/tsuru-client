// Copyright © 2023 tsuru-client authors
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wrappers

import "github.com/spf13/cobra"

// ForceCallPreRun climbs the parent tree until it finds and runs one of:
//   - PreRun()
//   - PreRunE()
//   - PersistentPreRun()
//   - PersistentPreRunE()
func ForceCallPreRun(cmd *cobra.Command, args []string) {
	curr := cmd
	for curr != nil {
		if curr.PreRun != nil {
			curr.PreRun(cmd, args)
			return
		}
		if curr.PreRunE != nil {
			_ = curr.PreRunE(cmd, args)
			return
		}
		if curr.PersistentPreRun != nil {
			curr.PersistentPreRun(cmd, args)
			return
		}
		if curr.PersistentPreRunE != nil {
			_ = curr.PersistentPreRunE(cmd, args)
			return
		}

		if curr == curr.Parent() {
			return
		}
		curr = curr.Parent()
	}
}