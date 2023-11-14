// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package debug

import (
	"io"

	"go.xrstf.de/otto/pkg/lang/ast"
)

var Indent = "  "

const doNotIndent = -1

func Dump(p *ast.Program, out io.Writer) error {
	return dumpProgram(p, out, 0)
}

func DumpSingleline(p *ast.Program, out io.Writer) error {
	return dumpProgram(p, out, -1)
}

func writeString(out io.Writer, str string) error {
	_, err := out.Write([]byte(str))
	return err
}
