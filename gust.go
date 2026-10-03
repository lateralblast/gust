// 2>/dev/null; e=$(mktemp); go build -o $e "$0"; $e "$@" ; r=$?; rm $e; exit $r

/*
Name:         gust (Golang Universal Shell script Template)
Version:      0.3.6
Release:      1
License:      CC BY-NC-SA (Creative Commons Attribution-NonCommercial-ShareAlike)
              https://creativecommons.org/licenses/by-nc-sa/4.0/legalcode
Group:        System
Source:       N/A
URL:          https://github.com/lateralblast/just
Distribution: UNIX
Vendor:       UNIX
Packager:     Richard Spindler <richard@lateralblast.com.au>
Description:  A template for writing golang shell scripts
*/

package main

// Import modules

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

// scriptVersion is the version printed by --version.
// Keep it in sync with the Version in the header, the README and the CHANGELOG.
const scriptVersion = "0.3.6"

// Argument describes a commandline argument/switch.
// Each argument is registered under both its long and short name, see addArgument.
type Argument struct {
	info     string
	short    string
	long     string
	category string // "switch", "option" or "action"
	function func()
}

var (
	// defaults holds the default value for each option.
	// It must contain an entry for each option created.
	defaults = map[string]string{
		"verbose":   "false",
		"force":     "false",
		"dryrun":    "false",
		"doactions": "false",
		"dooptions": "false",
		"help":      "all",
	}
	// options holds the current option values, initialised from defaults.
	options = map[string]string{}
	// arguments holds the commandline argument information.
	// It is populated by populateArguments.
	arguments = map[string]Argument{}
)

// capitalize upper cases the first letter of each word in a sentence.
func capitalize(sentence string) string {
	output := []rune{}
	isWord := true
	for _, val := range sentence {
		if isWord && unicode.IsLetter(val) {
			output = append(output, unicode.ToUpper(val))
			isWord = false
		} else {
			if !unicode.IsLetter(val) {
				isWord = true
			}
			output = append(output, val)
		}
	}
	return string(output)
}

// tabs returns the tabs needed to align a column, given the width of the text in it.
func tabs(text string, width int) string {
	if len(text) < width {
		return "\t\t"
	}
	return "\t"
}

// messageHeader converts a message format (e.g. enable) into a header (e.g. Enabling).
func messageHeader(format string) string {
	format = capitalize(strings.ToLower(format))
	switch {
	case strings.HasSuffix(format, "ing"):
		return format
	case strings.HasSuffix(format, "s"), strings.HasSuffix(format, "n"):
		return format + "ing"
	case strings.HasSuffix(format, "t"):
		return format + "ting"
	case strings.HasSuffix(format, "e"):
		return strings.TrimSuffix(format, "e") + "ing"
	case format == "Info":
		return "Information"
	}
	return format
}

// verboseMessage prints a consistently formatted message when verbose mode is enabled.
func verboseMessage(message, format string) {
	if verbose, _ := strconv.ParseBool(options["verbose"]); !verbose {
		return
	}
	header := messageHeader(format)
	fmt.Printf("%s:%s%s\n", header, tabs(header, 15), message)
}

// warningMessage displays a warning, overriding non verbose mode if needed.
func warningMessage(message string) {
	verbose := options["verbose"]
	options["verbose"] = "true"
	verboseMessage(message, "warn")
	options["verbose"] = verbose
}

// usageError displays a warning and the help information, then exits with an error code.
func usageError(message string) {
	warningMessage(message)
	options["help"] = "all"
	printHelp()
	os.Exit(1)
}

// addArgument registers an argument under both its long and short name.
// The function is optional, arguments without one (e.g. those that take a value)
// are handled in parseArguments.
func addArgument(info, short, long, category string, function func()) {
	argument := Argument{
		info:     info,
		short:    short,
		long:     long,
		category: category,
		function: function,
	}
	arguments[long] = argument
	arguments[short] = argument
}

// runArgument runs the function of an argument, given its long or short name.
// An argument that does not exist, or has no function, is a usage error.
func runArgument(name string) {
	argument, exists := arguments[name]
	if !exists {
		usageError("Commandline argument " + name + " does not exist")
	}
	if argument.function == nil {
		usageError("Commandline argument " + name + " can not be used here")
	}
	argument.function()
}

// checkCommand checks that a shell command exists.
func checkCommand(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

// linter runs golangci-lint over the script and exits with its exit code.
func linter() {
	command := "golangci-lint"
	if !checkCommand(command) {
		warningMessage("No linter found")
		os.Exit(0)
	}
	fmt.Println("Linter output:")
	output, err := exec.Command(command, "run", options["script"]).CombinedOutput()
	fmt.Println(string(output))
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			os.Exit(exitError.ExitCode())
		}
		warningMessage("Failed to run linter: " + err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}

// sortedKeys returns the keys of a map in sorted order.
// This keeps output consistent, as Go randomises map iteration order.
func sortedKeys[V any](values map[string]V) []string {
	keys := []string{}
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

// printHelpCategory prints help information for a specific category.
func printHelpCategory(category string) {
	fmt.Printf("Usage (%s):\n", category)
	fmt.Println()
	for _, key := range sortedKeys(arguments) {
		argument := arguments[key]
		// Each argument is registered under its long and short name, only print once
		if argument.category != category || key != argument.long {
			continue
		}
		fmt.Printf("%s, %s:%s%s\n", argument.long, argument.short, tabs(argument.long, 15), argument.info)
	}
	fmt.Println()
}

// printHelp prints the help information selected by the help option, without exiting.
func printHelp() {
	switch options["help"] {
	case "option", "options":
		printHelpCategory("option")
	case "switch", "switches":
		printHelpCategory("switch")
	case "action", "actions":
		printHelpCategory("action")
	default:
		printHelpCategory("switch")
		printHelpCategory("option")
		printHelpCategory("action")
	}
}

// help prints the help information and exits.
func help() {
	printHelp()
	os.Exit(0)
}

// version prints the version information and exits.
func version() {
	fmt.Println("Version:      " + scriptVersion)
	os.Exit(0)
}

// handleOptions handles a comma separated list of options.
//
//	e.g. verbose sets the verbose option to true
//	e.g. noverbose sets the verbose option to false
func handleOptions(values string) {
	for _, parameter := range strings.Split(values, ",") {
		original := parameter
		format := "enable"
		value := "true"
		// A name that is a valid option is enabled, otherwise check for a no prefix, e.g. noverbose
		_, exists := defaults[parameter]
		if !exists && strings.HasPrefix(parameter, "no") {
			format = "disable"
			value = "false"
			parameter = strings.TrimPrefix(parameter, "no")
			_, exists = defaults[parameter]
		}
		if !exists {
			usageError("Option " + original + " does not exist")
		}
		options[parameter] = value
		verboseMessage(parameter, format)
	}
}

// checkValue checks that the argument at argNum has a value, and handles the help argument.
func checkValue(argNum int) {
	parameter := os.Args[argNum]
	isHelp := strings.TrimLeft(parameter, "-") == "help" || parameter == "-h"
	if argNum == len(os.Args)-1 {
		if !isHelp {
			usageError("No value given for " + parameter)
		}
		options["help"] = "all"
		help()
	}
	value := os.Args[argNum+1]
	if strings.HasPrefix(value, "-") {
		warningMessage("No value given for " + parameter)
		os.Exit(1)
	}
	verboseMessage("Value given for "+parameter+" "+value, "info")
	if isHelp {
		options["help"] = value
		help()
	}
}

// printEnv prints the environment (options).
func printEnv() {
	fmt.Println("Environment (Options):")
	fmt.Println()
	for _, key := range sortedKeys(options) {
		if key == "script" {
			continue
		}
		fmt.Printf("%s:%s%s\t(default = %s)\n", key, tabs(key, 7), options[key], defaults[key])
	}
	fmt.Println()
}

// printDefs prints the default values of the environment (options).
func printDefs() {
	fmt.Println("Defaults (Options):")
	fmt.Println()
	for _, key := range sortedKeys(defaults) {
		fmt.Printf("%s:%s%s\n", key, tabs(key, 7), defaults[key])
	}
	fmt.Println()
}

// populateArguments registers the commandline arguments.
func populateArguments() {
	addArgument("Perform action", "a", "action", "switch", nil)
	addArgument("Set option", "o", "option", "switch", nil)
	addArgument("Enable dryrun mode", "d", "dryrun", "option", nil)
	addArgument("Enable verbose output", "v", "verbose", "option", nil)
	addArgument("Print help information", "h", "help", "action", help)
	addArgument("Check script with linter", "l", "linter", "action", linter)
	addArgument("Print Defaults", "D", "printdefs", "action", printDefs)
	addArgument("Print Environment", "E", "printenv", "action", printEnv)
	addArgument("Print version information", "V", "version", "switch", version)
}

// scanVerbose checks for the verbose option, so it is active while parsing.
// The last occurrence wins.
func scanVerbose(args []string) {
	for number, arg := range args {
		names := []string{}
		if strings.HasPrefix(arg, "-") {
			names = append(names, strings.TrimLeft(arg, "-"))
			if arg == "-v" {
				names = []string{"verbose"}
			}
		} else if number > 0 {
			switch args[number-1] {
			case "-o", "--option", "--options":
				names = strings.Split(arg, ",")
			}
		}
		for _, name := range names {
			switch name {
			case "verbose":
				options["verbose"] = "true"
			case "noverbose":
				options["verbose"] = "false"
			}
		}
	}
}

// isCombinedSwitch checks for a -abc style argument, i.e. a single dash and at least two letters.
func isCombinedSwitch(arg string) bool {
	letters := []rune(arg)
	return len(letters) > 2 && letters[0] == '-' && unicode.IsLetter(letters[1]) && unicode.IsLetter(letters[2])
}

// handleCombinedSwitch handles each letter of a -abc style argument, e.g. -abc > a, b, c.
func handleCombinedSwitch(arg string) {
	for _, letter := range strings.TrimPrefix(arg, "-") {
		name := string(letter)
		argument, exists := arguments[name]
		if !exists {
			usageError("Commandline argument " + name + " does not exist")
		}
		if argument.category == "option" {
			handleOptions(argument.long)
		} else {
			runArgument(name)
		}
	}
}

// parseArguments loops through the commandline arguments and handles them.
// Options and actions that take values are returned, so they can be handled once all arguments are parsed.
func parseArguments() (actionFlags, optionFlags []string) {
	for argNum := 1; argNum < len(os.Args); argNum++ {
		argName := os.Args[argNum]
		// Convert plural arguments to non plural
		argName = strings.ReplaceAll(argName, "options", "option")
		argName = strings.ReplaceAll(argName, "actions", "action")
		if isCombinedSwitch(argName) {
			handleCombinedSwitch(argName)
			continue
		}
		// Arguments that don't start with a dash are values, they are handled with their switch
		if !strings.HasPrefix(argName, "-") {
			continue
		}
		argName = strings.ReplaceAll(argName, "-", "")
		argument, exists := arguments[argName]
		if !exists {
			// Check if the argument is a negative option, e.g. noverbose, handleOptions checks it exists
			if !strings.HasPrefix(argName, "no") {
				usageError("Commandline argument " + argName + " does not exist")
			}
			handleOptions(argName)
			continue
		}
		if argument.category == "option" {
			handleOptions(argument.long)
			continue
		}
		// Handle arguments that take values, anything else runs its function
		switch argument.long {
		case "action":
			checkValue(argNum)
			actionFlags = append(actionFlags, os.Args[argNum+1])
			options["doactions"] = "true"
		case "option":
			checkValue(argNum)
			optionFlags = append(optionFlags, os.Args[argNum+1])
			options["dooptions"] = "true"
		case "help":
			checkValue(argNum)
		default:
			runArgument(argument.long)
		}
	}
	return actionFlags, optionFlags
}

// runActions runs each action, an action flag can be a comma separated list of actions.
func runActions(actionFlags []string) {
	for _, actionFlag := range actionFlags {
		for _, action := range strings.Split(actionFlag, ",") {
			verboseMessage("action flag "+action, "process")
			runArgument(action)
		}
	}
}

func main() {
	// Get script file
	_, scriptFile, _, _ := runtime.Caller(0)
	options["script"] = scriptFile
	populateArguments()
	// Copy defaults to options map
	for key, value := range defaults {
		options[key] = value
	}
	scanVerbose(os.Args[1:])
	// If we have no arguments print help information
	if len(os.Args) < 2 {
		options["help"] = "all"
		help()
	}
	actionFlags, optionFlags := parseArguments()
	// If we have option(s) handle each
	if doOptions, _ := strconv.ParseBool(options["dooptions"]); doOptions {
		for _, values := range optionFlags {
			handleOptions(values)
		}
	}
	// If we have action(s) handle each
	if doActions, _ := strconv.ParseBool(options["doactions"]); doActions {
		runActions(actionFlags)
	}
	os.Exit(0)
}
