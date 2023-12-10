// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package datetime

import (
	"time"

	"go.xrstf.de/rudi/pkg/eval/functions"
	"go.xrstf.de/rudi/pkg/eval/types"
)

var (
	Functions = types.Functions{
		"now": functions.NewBuilder(nowFunction).WithDescription("returns the current date & time (UTC), formatted like a Go date").Build(),
	}
)

func nowFunction(format string) (any, error) {
	return time.Now().Format(format), nil
}
