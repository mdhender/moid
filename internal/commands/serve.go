// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package commands

import (
	"context"
	"net/http"
)

type Serve struct{}

func (c *Serve) Execute(scheme, host, port string, mux http.Handler, ctx context.Context) error {
	panic("this has been deprecated")
}
