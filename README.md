# Rudi

<p align="center">
  <img src="./docs/rudi-portrait.png" alt="">
</p>

<p align="center">
  <img src="https://img.shields.io/github/v/release/xrstf/rudi" alt="last stable release">

  <a href="https://goreportcard.com/report/go.xrstf.de/rudi">
    <img src="https://goreportcard.com/badge/go.xrstf.de/rudi" alt="go report card">
  </a>

  <a href="https://pkg.go.dev/go.xrstf.de/rudi">
    <img src="https://pkg.go.dev/badge/go.xrstf.de/rudi" alt="godoc">
  </a>
</p>

Rudi is a Lisp-like, embeddable programming language that focuses on transforming data structures
like those available in JSON (numbers, bools, objects, vectors etc.). A statement in Rudi looks like

```lisp
(set! .foo[0] (+ (len .users) 42))
```

## Contents

* [Features](#features)
* [Installation](#installation)
* [Documentation](#documentation)
  * [Language Description](docs/language.md)
  * [Type Handling](docs/coalescing.md)
  * [Standard Library](docs/functions/README.md)
* [Usage](#usage)
  * [Command Line](#command-line)
  * [Embedding](#embedding)
* [Alternatives](#alternatives)
* [Credits](#credits)
* [License](#license)

## Features

* **Safe** evaluation: Rudi is not Turing-complete and so Rudi programs are always guaranteed to
  complete in a reasonable time frame.
* **Lightweight**: Rudi comes with only a single dependency on `go-cmp`, nothing else.
* **Hackable**: Rudi tries to keep the language itself approachable, so that modifications are
  easier and newcomers have an easier time to get started.
* **Variables** can be pre-defined or set at runtime.
* **JSONPath** expressions are first-class citizens and make referring to the current JSON document
  a breeze.
* **Optional Type Safety**: Choose between pedantic, strict or humane typing for your programs.
  Strict allows nearly no type conversions, humane allows for things like `1` (int) turning into
  `"1"` (string) when needed.

## Installation

Rudi is primarily meant to be embedded into other Go programs, but a standalone CLI application,
`rudi`, is also available to test your scripts with. `rudi` can be installed using Git & Go. Rudi
requires **Go 1.18** or newer.

```bash
git clone https://github.com/xrstf/rudi
cd rudi
make build
```

Alternatively, you can download the [latest release](https://github.com/xrstf/rudi/releases/latest)
from GitHub.

## Documentation

Make yourself familiar with Rudi using the documentation:

* The [Language Description](docs/language.md) describes the Rudi syntax and semantics.
* All built-in functions are described in the [standard library](docs/functions/README.md).
* [Type Handling](docs/coalescing.md) describes how Rudi handles, converts and compares values.

## Usage

### Command Line

Rudi comes with a standalone CLI tool called `rudi`.

```
Usage of rudi:
      --coalesce string   Type conversion handling, choose one of strict, pedantic or humane. (default "strict")
      --debug-ast         Output syntax tree of the parsed script in non-interactive mode.
  -h, --help              Show help and documentation.
  -i, --interactive       Start an interactive REPL to run expressions.
  -p, --pretty            Output pretty-printed JSON.
  -s, --script string     Load Rudi script from file instead of first argument (only in non-interactive mode).
  -V, --version           Show version and exit.
  -y, --yaml              Output pretty-printed YAML instead of JSON.
```

`rudi` can run in one of two modes:

* **Interactive Mode** is enabled by passing `--interactive` (or `-i`). This will start a REPL
  session where Rudi scripts are read from stdin and evaluated against the loaded files.
* **Script Mode** is used the an Rudi script is passed either as the first argument or read from a
  file defined by `--script`. In this mode `rudi` will run all statements from the script and print
  the resulting value, then it exits.

    Examples:

    * `rudi '.foo' myfile.json`
    * `rudi '(set! .foo "bar") (set! .users 42) .' myfile.json`
    * `rudi --script convert.rudi myfile.json`

`rudi` has extensive help built right into it, try running `rudi help` to get started.

#### File Handling

The first loaded file is known as the "document". Its content is available via path expressions like
`.foo[0]`. All loaded files are also available via the `$files` variable (i.e. `.` is the same as
`$files[0]` for reading, but when writing data, there is a difference between both notations; refer
to the docs for `set` for more information).

### Embedding

Rudi is well suited to be embedded into Go applications. A clean and simple API makes it a breeze:

```go
package main

import (
   "fmt"
   "log"

   "go.xrstf.de/rudi"
)

const script = `(set! .foo 42) (+ $myvar 42 .foo)`

func main() {
   // Rudi programs are meant to manipulate a document (path expressions like
   // ".foo" resolve within that document). The document can be anything,
   // but is most often a JSON object.
   documentData := map[string]any{"foo": 9000}

   // parse the script (the name is used when generating error strings)
   program, err := rudi.ParseScript("myscript", script)
   if err != nil {
      log.Fatalf("The script is invalid: %v", err)
   }

   // evaluate the program;
   // this returns an evaluated value, which is the result of the last expression
   // that was evaluated, plus the final document state (the updatedData) after
   // the script has finished.
   updatedData, result, err := program.Run(
      documentData,
      // setup the set of variables available by default in the script
      rudi.NewVariables().Set("myvar", 42),
      // Likewise, setup the functions available (note that this includes
      // functions like "if" and "and", so running with an empty function set
      // is generally not advisable).
      rudi.NewBuiltInFunctions(),
      // Decide what kind of type strictness you would like; pedantic, strict
      // or humane; choose your own adventure (strict is default if you use nil
      // here; humane allows conversions like 1 == "1").
      coalescing.NewStrict(),
   )
   if err != nil {
      log.Fatalf("Script failed: %v", err)
   }

   fmt.Println(result)       // => 126
   fmt.Println(updatedData)  // => {"foo": 42}
}
```

## Alternatives

Rudi doesn't exist in a vacuum; there are many other great embeddable programming/scripting languages
out there, allbeit with slightly different ideas and goals than Rudi:

* [Anko](https://github.com/mattn/anko) – Go-like syntax and allows recursion, making it more
  dangerous and hard to learn for non-developers than I'd like.
* [ECAL](https://github.com/krotik/ecal) – Is an event-based system using rules which are triggered by
  events; comes with recursion as well and is therefore out.
* [Expr](https://github.com/antonmedv/expr), [GVal](https://github.com/PaesslerAG/gval),
  [CEL](https://github.com/google/cel-go) – Great languages for writing a single expression, but not
  suitable for transforming/mutating data structures.
* [Gentee](https://github.com/gentee/gentee) – Is similar to C/Python and allows recursion, so both
  to powerful/dangerous and not my preference in terms of syntax.
* [Jsonnet](https://github.com/google/go-jsonnet) – Probably one of the most obvious alternatives
  among this list. Jsonnet shines when constructing new elements and complexer configurations
  out of smaller pieces of information, less so when manipulating objects. Also I personally really
  am no fan of Jsonnet's syntax, plus: NIH.
* [Starlark](https://github.com/google/starlark-go) – Is the language behind Bazel and actually has
  an optional nun-Turing-complete mode. However I am really no fan of its syntax and have not
  investigated it further.

## Credits

Rudi has been named after my grandfather.

Thanks to [@embik](https://github.com/embik) and [@xmudrii](https://github.com/xmudrii) for enduring
my constant questions for feedback :smile:

Rudi has been made possible by the amazing [Pigeon](https://github.com/mna/pigeon) parser generator.

## License

MIT
