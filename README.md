![Cat in wind](https://raw.githubusercontent.com/lateralblast/gust/master/gust.gif)

GUST
----

Golang Universal Shell script Template

Version
-------

Version 0.1.4

To get version:

``
./gust.go --version
``

Introduction
------------

This is loosely bast on just (Just a UNIX Shell Template):

https://github.com/lateralblast/just

This is intended to provide a template for writing golang shell scripts.
It gathers some arguably good and bad practices I've acquired over
the years to handle command line arguments and inline documentation.

This script is designed to have some more features than the standard flag
module that I like, and less complexity of a non standard module such as cobra.

This script utilises a special header (hashbang style) to be able to
be able to run as a script so it can be called directly, rather than
having to use "go run"

The header used is:

```
// 2>/dev/null; e=$(mktemp); go build -o $e "$0"; $e "$@" ; r=$?; rm $e; exit $r
```

Goals
-----

The goals of this script template are to:

Provide a command line processor in golang that:

- Can handle both short (e.g -h) and long command line arguments (e.g --help)
- Is able to process multiple actions and options so the script doesn't need to be called multiple times
  - Multiple actions and options are comma separated, or specified in multiple occurences of a switch
- Uses inline tags to semi-auto document script for the help/usage argument

Provide a base set of functionality that has:

- Help information
- Version information
- Debug and verbose command-line arguments/options to help with code debugging/quality
- Some additional base code checking capability (call out to shellcheck)
- Dry-run mode capability
- The ability to split larger scripts into modules which are loaded at run time

Choices
-------

This section details some of the choices taken in the template.

- Use true/false as values rather than 0/1 as I often generate YAML and other config files that use these as values
- Put environment variables and command line parameters into variables so that they are given more interpretable names
- Only run inline information gathering (version, help, etc) when called to reduce start up time
- Only use standard core modules

Requirements
------------

Standard modules required:

- os/exec
- unicode
- runtime
- strings
- regexp
- bufio
- fmt
- os

Workflow
--------

The script has the following workflow

- Get some environment variables
- If given no switches/parameters print help and exit
- Set defaults
- Parse command line switches/parameters
- If options switch has values process them
- Reset/updates defaults base on command line switches/parameters
- Perform action(s)

Defaults
--------

Defaults for options are set in map, e.g.

```
defaults = map[string]bool{
  "verbose": false,
  "force":   false,
  "dryrun":  false,
}
```

Options
-------

Options can be specified in two ways, as commandline arguments without values (e.g. --verbose),
or with the options commandline arguement (e.g. --option verbose)

Options are stored in a map, which is initiated from the defaults map specified in previous section.

Print options:

```
./gust.go --action printenv --options dryrun
Environment (Options):

dryrun:   true  (default = false)
verbose:  false (default = false)
force:    false (default = false)
```

Defaults
--------

As per previous section, defaults are stored in a map, and options inherits them at start.

Print defaults:

```
./gust.go --action printdefs
Defaults (Options):

force:    false
dryrun:   false
verbose:  false
```

Arguments
---------

Arguments are handled in a map of structs, e.g.

```
type Argument struct {
  info      string
  short     string
  long      string
  category  string
  value     string
  dynamic   bool
}

arguments = map[string]Argument {
  "action": {
    info:     "Perform action",
    short:    "a",
    long:     "action",
    category: "switch",
    value:    "",
    dynamic:  false,
  },
  "a": {
    info:     "Perform action",
    short:    "a",
    long:     "action",
    category: "switch",
    value:    "",
    dynamic:  false,
  },
  "version": {
    info:     "Print version information",
    short:    "V",
    long:     "version",
    category: "switch",
    value:    "",
    dynamic:  true,
  },
  "V": {
    info:     "Print version information",
    short:    "V",
    long:     "version",
    category: "switch",
    value:    "",
    dynamic:  false,
  },
  "printdefs": {
    info:     "Print Defaults",
    short:    "printdefs",
    long:     "",
    category: "action",
    value:    "",
    dynamic:  false,
  },
  "printenv": {
    info:     "Print Environment",
    short:    "printenv",
    long:     "",
    category: "action",
    value:    "",
    dynamic:  false,
  },
}
```

At the moment, both the short (e.g. -V) and the long version (e.g. --version) must be specified.

Switches (arguments that take values/parameters) are specified with the category "switch"

Options (arguments that don't take values/parameters) are specified with the category "option"

Actions (actions performed by switches) are specified with the category "action"

The info entry specifies the information that is given to the help routine.

The dynamic entry specifies whether a dynamic function/routine is available (one that can be called via its argument name)

Dynamic
-------

Dynamic routines/functions are ones that can be called via their action name, e.g. --version

In the arguments mapped struct, a entry is need that has the dynamic boolean set to true:

```
arguments = map[string]Argument {
  "version": {
    info:     "Print version information",
    short:    "V",
    long:     "version",
    category: "switch",
    value:    "",
    dynamic:  true,
  },
}
```

Then a function is required, e.g.:

```
func (dynamic Dynamic) Version() {
  // Code here
}

```

The function is then called as part of the commandline argument processing.

An example of code to call the dynamic routine:

```
instance := Dynamic{}
method   := reflect.ValueOf(instance).MethodByName("Verbose")
method.Call(nil)
```

Output
------

In general output is run through the verbose_message function, which handles formatting and outputs appropriate
prefixes, e.g. "Notice", "Warning", etc. This is designed to make sure out is consistently formatted.
It also allows output to be done only when the verbose mode is set so the script runs quietly if needed.

Check Values
------------

There is a check_values function that can be run with switches that take values.
This is a simple check to see if the value, if it starts with a "--" then it assumes that you haven't provided
a value and it's processing the next switch and exit. This can be overridden by using the --force switch.

Help
----

Help information is provided by thee argument struct, so any new arguments added will be available to help.

```
./gust.go --help
Usage (switch):

option, o:    Set option
action, a:    Perform action
version, V:   Print version information
help, h:      Print help information

Usage (option):

verbose, v:   Enable verbose output
dryrun, d:    Enable dryrun mode

Usage (action):

printdefs:    Print Defaults
linter:       Check script with linter
printenv:     Print Environment
```
Linter
------

The script has an action to run the golang linter against itself if installed:

```
./gust.go --action linter
```

Examples
--------

Run linter and printdefs as two actions in one command using two action switches:

```
./gust.go --action linter --action printdefs
```

Run linter and printdefs as two actions in one command using one action switch:

```
./gust.go --action linter,printdefs
```

The same can be done for options
