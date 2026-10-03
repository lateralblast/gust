# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [0.3.6] - 2026-10-03

### Changed

- Linter action prints stderr as well as stdout and exits with the exit code of golangci-lint, rather than always 0

## [0.3.5] - 2026-10-03

### Changed

- Split main into smaller functions: scanVerbose, parseArguments, handleCombinedSwitch and runActions

## [0.3.4] - 2026-10-03

### Changed

- Replaced the repeated short and long Argument entries with an addArgument helper, so each argument is declared once
- Help output is printed once per argument

## [0.3.3] - 2026-10-03

### Changed

- Removed the regexp module, replaced with the strings package and direct comparisons
- Simplified verboseMessage, handleOptions and printenv, and removed dead code (the unreachable "verbose" format and the no-long-name help branch)
- Warning and error messages no longer modify the verbose option permanently

## [0.3.2] - 2026-10-03

### Changed

- Replaced block comments with godoc style doc comments

## [0.3.1] - 2026-10-03

### Changed

- Renamed functions and the version constant from snake_case to camelCase, as per Go naming conventions
- Renamed variables that shadowed the error builtin and the regexp package

## [0.3.0] - 2026-10-03

### Changed

- Formatted with gofmt (tab indentation, sorted imports)

## [0.2.9] - 2026-10-03

### Fixed

- Unknown help category (e.g. `--help bogus`) printed nothing, now prints all help

## [0.2.8] - 2026-10-03

### Fixed

- Help, printenv and printdefs output order changed between runs, now sorted (replaces the bufio import with sort)

## [0.2.7] - 2026-10-03

### Fixed

- Version output depended on the source file still existing at its build path, now uses a `script_version` constant

## [0.2.6] - 2026-10-03

### Fixed

- check_value treated any switch containing "h" as the help switch

## [0.2.5] - 2026-10-03

### Fixed

- The verbose pre-scan matched "verbose" anywhere in the arguments, now only matches the verbose option
- Negated options (e.g. `--noverbose`) were split on every "no" in the name, now only the prefix is removed

## [0.2.4] - 2026-10-03

### Fixed

- Stray debug output (the argument name) printed before the warning for an unknown letter in a `-abc` style switch

## [0.2.3] - 2026-10-03

### Fixed

- Options were not validated, `--option foo` silently set an unknown option, now it is an error

## [0.2.2] - 2026-10-03

### Fixed

- Errors (unknown arguments, missing values, unknown actions or options) exited with code 0, now exit with code 1

## [0.2.1] - 2026-10-03

### Fixed

- Linter action never ran as check_command called the shell builtin `command`, now uses exec.LookPath

## [0.2.0] - 2026-10-03

### Fixed

- `-d` printed help instead of enabling dryrun because of a duplicate argument entry
- `-l`, `-D` and `-E` did nothing because their short entries had no long name, and `-e` was not defined
- Short names are now `-l` (linter), `-D` (printdefs) and `-E` (printenv)

## [0.1.9] - 2026-10-03

### Fixed

- Nil pointer panic for `-a` or `-o` in a combined switch (e.g. `-av`), now a usage error

## [0.1.8] - 2026-10-03

### Changed

- Changed license from CC BY to CC BY-NC-SA 4.0 and added LICENSE file
- Converted changelog to CHANGELOG.md (Keep a Changelog format)

### Added

- CLAUDE.md guidance for Claude Code

## [0.1.7] - 2024-12-30

### Changed

- Moved from dynamic functions to functions in structs

## [0.1.6] - 2024-12-28

### Changed

- Converted more routines to dynamic

## [0.1.5] - 2024-12-26

### Fixed

- Fixes recommended by linter

## [0.1.4] - 2024-12-26

### Added

- Code for dynamic functions

## [0.1.3] - 2024-12-26

### Changed

- Updated documentation

## [0.1.2] - 2024-12-26

### Added

- Code to check OS command exists

## [0.1.1] - 2024-12-26

### Fixed

- Tabulation of help output

## [0.1.0] - 2024-12-26

### Changed

- Updated documentation

## [0.0.9] - 2024-12-26

### Fixed

- Regular expressions in loops, as per linter recommendations

## [0.0.8] - 2024-12-26

### Added

- Action category to arguments struct map

## [0.0.7] - 2024-12-26

### Fixed

- Options map

## [0.0.6] - 2024-12-26

### Added

- Routine to print defaults

### Changed

- Updated code to print environment

## [0.0.5] - 2024-12-25

### Fixed

- Version code

## [0.0.4] - 2024-12-25

### Changed

- Updated documentation

## [0.0.3] - 2024-12-25

### Changed

- Cleanup based on golangci-lint recommendations

## [0.0.2] - 2024-12-25

### Added

- printenv function

### Changed

- Code cleanup

## [0.0.1]

### Added

- Initial working edition with some features and documentation

