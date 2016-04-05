// Copyright 2016 tsuru-client authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/tsuru/gnuflag"
	"github.com/tsuru/tsuru/cmd"
	"github.com/tsuru/tsuru/errors"
)

type appSwap struct {
	cmd.Command
	force     bool
	cnameOnly bool
	fs        *gnuflag.FlagSet
}

func (s *appSwap) Info() *cmd.Info {
	return &cmd.Info{
		Name:  "app-swap",
		Usage: "app-swap <app1-name> <app2-name> [-f/--force] [-c/--cname-only]",
		Desc: `Swaps routing between two apps. This allows zero downtime and makes rollback
as simple as swapping the applications back.

Use [[--force]] if you want to swap applications with a different number of
units or different platform without confirmation.

Use [[--cname-only]] if you want to swap all cnames except the default 
cname of application`,
		MinArgs: 2,
	}
}

func (s *appSwap) Flags() *gnuflag.FlagSet {
	if s.fs == nil {
		s.fs = gnuflag.NewFlagSet("", gnuflag.ExitOnError)
		s.fs.BoolVar(&s.force, "force", false, "Force Swap among apps with different number of units or different platform.")
		s.fs.BoolVar(&s.force, "f", false, "Force Swap among apps with different number of units or different platform.")
		s.fs.BoolVar(&s.cnameOnly, "cname-only", false, "Swap all cnames except the default cname.")
		s.fs.BoolVar(&s.cnameOnly, "c", false, "Swap all cnames except the default cname.")
	}
	return s.fs
}

func (s *appSwap) Run(context *cmd.Context, client *cmd.Client) error {
	url, err := cmd.GetURL(fmt.Sprintf("/swap?app1=%s&app2=%s&force=%t&cnameOnly=%t", context.Args[0], context.Args[1], s.force, s.cnameOnly))
	if err != nil {
		return err
	}
	err = makeSwap(client, url)
	if err != nil {
		if e, ok := err.(*errors.HTTP); ok && e.Code == http.StatusPreconditionFailed {
			var answer string
			fmt.Fprintf(context.Stdout, "WARNING: %s.\nSwap anyway? (y/n) ", strings.TrimRight(e.Message, "\n"))
			fmt.Fscanf(context.Stdin, "%s", &answer)
			if answer == "y" || answer == "yes" {
				url, _ = cmd.GetURL(fmt.Sprintf("/swap?app1=%s&app2=%s&force=%t&cnameOnly=%t", context.Args[0], context.Args[1], true, s.cnameOnly))
				return makeSwap(client, url)
			}
			fmt.Fprintln(context.Stdout, "swap aborted.")
			return nil
		}
		return err
	}
	fmt.Fprintln(context.Stdout, "Apps successfully swapped!")
	return err
}

func makeSwap(client *cmd.Client, url string) error {
	request, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return err
	}
	_, err = client.Do(request)
	return err
}
