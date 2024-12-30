// 2>/dev/null; e=$(mktemp); go build -o $e "$0"; $e "$@" ; r=$?; rm $e; exit $r

/*
Name:         gust (Golang Universal Shell script Template)
Version:      0.1.7
Release:      1
License:      CC-BA (Creative Commons By Attribution)
              http://creativecommons.org/licenses/by/4.0/legalcode
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
  "strconv"
  "os/exec"
  "unicode"
  "runtime"
  "strings"
  "regexp"
  "bufio"
  "fmt"
  "os"
)

// Create a structure to manage commandline arguments/switches

type Argument struct {
  info        string
  short       string
  long        string
  category    string
  function    func()
}

// Initialize variables

var (
  // Create a map to store default values for options
  // This should contain a default value for each option created
  defaults = map[string]string{
    "verbose":    "false",
    "force":      "false",
    "dryrun":     "false",
    "doactions":  "false",
    "dooptions":  "false",
    "help":       "all",
  }
  // Create options map
  options = map[string]string{}
  // Create a map of Argument structs to store commanline argument information
  // This gets populated in the populate_arguments function
  arguments = map[string]Argument{}
)


/*
Function:     capitalize 
Parameters:   sentence
Description:  A routine to capitalize a sentence
*/

func capitalize(sentence string) string {
  var output []rune    //create an output slice
  isWord := true
  for _, val := range sentence {
    if isWord && unicode.IsLetter(val) {  //check if character is a letter convert the first character to upper case
      output = append(output, unicode.ToUpper(val))
      isWord = false
    } else if !unicode.IsLetter(val) {
      isWord = true
      output = append(output, val)
    } else {
      output = append(output, val)
    }
  }
  sentence = string(output)
  return sentence
}

/*
Function:     verbose_message
Parameters:   message and formet
Description:  A routine to create consistently formatted output
*/

func verbose_message(message, format string) {
  var header string
  format = strings.ToLower(format)
  format = capitalize(format)
  matches, _ := regexp.MatchString("verbose", format)
  if matches {
    fmt.Println(message) 
  } else {
    matches, _ = strconv.ParseBool(options["verbose"])
    if matches {
      matches, _ := regexp.MatchString("ing$", format)
      if matches {
        header = format
      } else {
        matches, _ := regexp.MatchString("s$|n$", format)
        if matches {
          header = format+"ing"
        } else {
          matches, _ := regexp.MatchString("t$", format)
          if matches {
            header = format+"ting"
          } else {
            matches, _ := regexp.MatchString("e$", format)
            if matches {
              header = string(format[:len(format)-1])
              header = header+"ing"
            } else {
              matches, _ := regexp.MatchString("^Info$", format)
              if matches {
                header = "Information"
              } else {
                header = format
              }
            }
          }
        }
      }
      if len(header) < 15 {
        fmt.Printf("%s:\t\t%s\n", header, message)
      } else {
        fmt.Printf("%s:\t%s\n", header, message)
      }
    }
  }
}

/*
Function:     warning_message
Parameters:   message
Description:  A routine to display a warning, overriding non verbose mode if needed
*/

func warning_message(message string) {
  matches, _ := strconv.ParseBool(options["verbose"])
  if matches {
    verbose_message(message, "warn")
  } else {
    options["verbose"] = "true"
    verbose_message(message, "warn")
    options["verbose"] = "false"
  }
}

/*
Function:     check_command
Parameters:   command
Description:  A routine to check that a shell command exists
*/

func check_command(command string) bool {
  exists := false
  shell  := exec.Command("command", "-v", command)
  stdout, _ := shell.Output()
  output := string(stdout)
  matches, _ := regexp.MatchString(command, output)
  if matches {
    exists = true
  } else {
    exists = false
  }
  return exists
}

/*
Function:     linter 
Parameters:   script_file
Description:  A routine to run linter over script
*/

func linter() {
  command := "golangci-lint"
  exists  := check_command(command)
  if exists {
    fmt.Println("Linter output:")
    script_file := options["script"]
    shell := exec.Command(command, "run", script_file)
    stdout, _ := shell.Output()
    output := string(stdout)
    fmt.Println(output)
  } else {
    warning_message("No linter found")
  }
  os.Exit(0)
}

/*
Function:     print_help_category
Parameters:   category
Description:  A routine to print help information for a specific category
*/

func print_help_category(category string) {
  fmt.Printf("Usage (%s):\n", category)
  fmt.Println("")
  for key, argument := range arguments {
    matches, _ := regexp.MatchString(category, argument.category)
    if matches {
      if len(key) > 1 {
        if len(argument.long) <1 {
          if len(argument.short) < 7 {
            fmt.Printf("%s:\t\t\t%s\n", argument.short, argument.info)
          } else {
            fmt.Printf("%s:\t\t%s\n", argument.short, argument.info)
          }
        } else {
          if len(argument.long) < 15 {
            fmt.Printf("%s, %s:\t\t%s\n", argument.long, argument.short, argument.info)
          } else {
            fmt.Printf("%s, %s:\t%s\n", argument.long, argument.short, argument.info)
          }
        }
      }
    }
  }
  fmt.Println("")
}

/*
Function:     print_help
Parameters:   options
Description:  A routine to print help information
*/

func help() {
  switch options["help"] {
    case "option", "options":
      print_help_category("option")
    case "switch", "switches":
      print_help_category("switch")
    case "action", "actions":
      print_help_category("action")
    case "all":
      print_help_category("switch")
      print_help_category("option")
      print_help_category("action")
  }
  os.Exit(0)
}

/*
Function:     version
Parameters:   options 
Description:  A routine to print version information
*/

func version() {
  script_file := options["script"]
  open_file, file_error := os.Open(script_file)
  if file_error != nil {
      fmt.Println(file_error)
  }
  defer open_file.Close()
  regexp  := regexp.MustCompile("[0-9]")
  scanner := bufio.NewScanner(open_file)
  for scanner.Scan() {
    line := scanner.Text()
    if strings.Contains(line, "Version:") {
      matches := regexp.MatchString(line)
      if matches {
        fmt.Println(line)
      }
    }
  }
  os.Exit(0)
}

/*
Function:     handle_options 
Parameters:   values
Description:  A routine to handle otions
              e.g. --verbose sets the verbose option to true
              e.g. --noverbose sets the verbose option to false
*/

func handle_options(values string) {
  parameters := []string{}
  matches, _ := regexp.MatchString(",", values)
  if matches {
    parameters = strings.Split(values, ",")
  } else {
    parameters = append(parameters, values)
  }
  regexp := regexp.MustCompile("^no")
  for number := 0 ;  number < len(parameters) ; number++ {
    parameter := parameters[number]
    matches   := regexp.MatchString(parameter)
    format := ""
    if matches {
      format    = "disable"
      parameter = strings.Split(parameter, "no")[1]
      options[parameter] = "false"
    } else {
      format = "enable"
      options[parameter] = "true"
    }
    verbose_message(parameter, format)
  }
}

/*
Function:     check_value
Parameters:   arg_num
Description:  A routine to handle argument values
*/

func check_value(arg_num int) {
  parameter := os.Args[arg_num]
  if arg_num == len(os.Args)-1 {
    message := "No value given for " + parameter
    switch parameter {
      case "--help", "-h":
        options["help"] = "all"
        help()
      default:
        verbose_message(message, "warn")
        options["help"] = "all"
        help()
    }
    os.Exit(1)
  }
  check_value := os.Args[arg_num+1] 
  matches, _ := regexp.MatchString("^-", check_value)
  if matches {
    message := "No value given for " + parameter
    options["verbose"] = "true"
    verbose_message(message, "warn")
    os.Exit(1)
  } else {
    message := "Value given for " + parameter + " " + check_value
    verbose_message(message, "info")
    matches, _ := regexp.MatchString("help|h", parameter)
    if matches {
      value := os.Args[arg_num+1]
      options["help"] = value
      help()
    }
  }
}

/*
Function:     printenv
Parameters:   none
Description:  A routine to print environment variables (options)
*/

func printenv() {
  fmt.Println("Environment (Options):")
  fmt.Println()
  regexp := regexp.MustCompile("script")
  for key, value := range options {
    matches := regexp.MatchString(key)
    if (!matches) {
      def := defaults[key]
      if len(key) < 7 {
        fmt.Printf("%s:\t\t%s\t(default = %s)\n", key, value, def)
      } else {
        fmt.Printf("%s:\t%s\t(default = %s)\n", key, value, def)
      }
    }
  }
  fmt.Println()
}

/*
Function:     printdefs 
Parameters:   none
Description:  A routine to print default environment variables (options)
*/

func printdefs() {
  fmt.Println("Defaults (Options):")
  fmt.Println()
  for key, value := range defaults {
    if len(key) < 7 {
      fmt.Printf("%s:\t\t%s\n", key, value)
    } else {
      fmt.Printf("%s:\t%s\n", key, value)
    }
  }
  fmt.Println()
}

func populate_arguments() {
  arguments["action"] = Argument{
    info:     "Perform action",
    short:    "a",
    long:     "action",
    category: "switch",
  }
  arguments["a"] = Argument{
    info:     "Perform action",
    short:    "a",
    long:     "action",
    category: "switch",
  }
  arguments["option"] = Argument{
    info:     "Set option",
    short:    "o",
    long:     "option",
    category: "switch",
  }
  arguments["o"] = Argument{
    info:     "Set option",
    short:    "o",
    long:     "option",
    category: "switch",
  }
  arguments["dryrun"] = Argument{
    info:     "Enable dryrun mode",
    short:    "d",
    long:     "dryrun",
    category: "option",
  }
  arguments["d"] = Argument{
    info:     "Enable dryrun mode",
    short:    "d",
    long:     "dryrun",
    category: "option",
  }
  arguments["d"] = Argument{
    info:     "Print help information",
    short:    "h",
    long:     "help",
    category: "switch",
  }
  arguments["verbose"] = Argument{
    info:     "Enable verbose output",
    short:    "v",
    long:     "verbose",
    category: "option",
  }
  arguments["v"] = Argument{
    info:     "Enable verbose output",
    short:    "v",
    long:     "verbose",
    category: "option",
  }
  arguments["help"] = Argument{
    info:     "Print help information",
    short:    "h",
    long:     "help",
    category: "action",
    function: func() {
      help()
    },
  }
  arguments["h"] = Argument{
    info:     "Print help information",
    short:    "h",
    long:     "help",
    category: "action",
  }
  arguments["linter"] = Argument{
    info:     "Check script with linter",
    short:    "l",
    long:     "linter",
    category: "action",
    function: func() {
      linter()
    },
  }
  arguments["l"] = Argument{
    info:     "Check script with linter",
    short:    "linter",
    long:     "",
    category: "action",
  }
  arguments["printdefs"] = Argument{
    info:     "Print Defaults",
    short:    "d",
    long:     "printdefs",
    category: "action",
    function: func() {
      printdefs()
    },
  }
  arguments["D"] = Argument{
    info:     "Print Defaults",
    short:    "printdefs",
    long:     "",
    category: "action",
  }
  arguments["printenv"] = Argument{
    info:     "Print Environment",
    short:    "e",
    long:     "printenv",
    category: "action",
    function: func() {
      printenv()
    },
  }
  arguments["E"] = Argument{
    info:     "Print Environment",
    short:    "printenv",
    long:     "",
    category: "action",
  }
  arguments["version"] = Argument{
    info:     "Print version information",
    short:    "V",
    long:     "version",
    category: "switch",
    function: func() {
      version()
    },
  }
  arguments["V"] = Argument{
    info:     "Print version information",
    short:    "V",
    long:     "version",
    category: "switch",
  }
}

// Main function

func main() {
  // Get script file
  _, script_file, _, _ := runtime.Caller(0)
  options["script"] = script_file
  populate_arguments()
  // Copy defaults to options map
  for key, value := range defaults {
    options[key] = value
  }
  // Create arrays to store actions or options
  action_flags := []string{}
  option_flags := []string{}
  // Save CLI arguments and check for verbose option
  cli_args := strings.Join([]string(os.Args), " ")
  matches, _ := regexp.MatchString("noverbose", cli_args)
  if matches {
    options["verbose"] = "false"
  } else {
    matches, _ := regexp.MatchString("verbose", cli_args)
    if matches {
      options["verbose"] = "true"
    }
  }
  // If we have no arguments print help information
  if len(os.Args) < 2 {
    options["help"] = "all"
    help()
  }
  regexp1 := regexp.MustCompile("^-[a-z,A-Z][a-z,A-Z]")
  regexp2 := regexp.MustCompile("^-")
  regexp3 := regexp.MustCompile("option")
  regexp4 := regexp.MustCompile(",")
  regexp5 := regexp.MustCompile("^no")
  regexp6 := regexp.MustCompile("action|switch")
//  regexp6 := regexp.MustCompile("action")
  // loop through command line arguments and handle them
  for arg_num := 1 ; arg_num < len(os.Args) ; arg_num++ {
    arg_name := os.Args[arg_num]
    // Convert plural arguments to non plural
    arg_name = strings.Replace(arg_name, "options", "option", -1)
    arg_name = strings.Replace(arg_name, "actions", "action", -1)
    // Check if we have a -abc style switch and process
    matches := regexp1.MatchString(arg_name)
    if matches {
      // Strip -
      arg_names := strings.Split(arg_name, "-")[1]
      // Step though each command line arguement, e.g. -abc > a, b, c,
      letters   := strings.Split(arg_names, "")
      for num :=0 ; num < len(letters) ; num++ {
        letter := letters[num] 
        _, exists := arguments[letter]
        if (exists) {
          long_name := arguments[letter].long
          // Check that an argument structure exists and grab the long version
          matches := regexp3.MatchString(arguments[letter].category)
          if matches {
            handle_options(long_name)
          } else {
            arguments[long_name].function()
          }
        } else {
          fmt.Println(arg_name)
          // Print help if there is no argument structure
          message := "Commandline argument "+letter+" does not exist"
          warning_message(message)
          options["help"] = "all"
          help()
        }
      }
    } else {
      matches := regexp2.MatchString(arg_name)
      if matches {
        // Strip -
        arg_name = strings.Replace(arg_name, "-", "", -1)
        // Check argument structure exists
        _, exists := arguments[arg_name]
        if exists {
          // If argument structure exists check if it is an option and handle
          long_name := arguments[arg_name].long
          matches   := regexp3.MatchString(arguments[long_name].category)
          if matches {
            handle_options(long_name)
          } else {
            _, exists := arguments[long_name]
            if exists {
              // If argument is not an option, handle appropriatly
              switch long_name {
                case "action":
                  check_value(arg_num)
                  action_flags = append(action_flags, os.Args[arg_num+1])
                  options["doactions"] = "true"
                case "option":
                  check_value(arg_num)
                  option_flags = append(option_flags, os.Args[arg_num+1])
                  options["dooptions"] = "true"
                case "help":
                  check_value(arg_num)
                default:
                  matches := regexp6.MatchString(arguments[long_name].category)
                  if matches {
                    arguments[long_name].function()
                  }else {
                    options["help"] = "all"
                    help()
                  }
              }
            }
          }
        } else {
          // check if argument is a negative option, e.g. noverbose and handle
          matches := regexp5.MatchString(arg_name)
          if matches {
            parameter  := strings.Split(arg_name, "no")[1]
            matches   := regexp3.MatchString(arguments[parameter].category)
            if matches {
              handle_options(arg_name)
            } else {
              // If argument structure does exist warn and print help
              message := "Commandline argument "+arg_name+" does not exist"
              warning_message(message)
              options["help"] = "all"
              help()
            }
          } else {
            long_name := arguments[arg_name].long
            matches   := regexp6.MatchString(arguments[long_name].category)
            if matches {
              arguments[long_name].function()
            } else {
              // If argument structure does exist warn and print help
              message := "Commandline argument "+arg_name+" does not exist"
              warning_message(message)
              options["help"] = "all"
              help()
            }
          }
        }
      }
    }
  } 
  // If we have option(s) handle each
  do_options, _ := strconv.ParseBool(options["dooptions"])
  if do_options {
    for number := 0 ; number < len(option_flags) ; number++ {
      values := option_flags[number]
      handle_options(values)
    }
  }
  // If we have action(s) handle each
  do_actions, _ := strconv.ParseBool(options["doactions"])
  if do_actions {
    for number := 0 ; number < len(action_flags) ; number++ {
      action_list := []string{}
      action_name := action_flags[number]
      matches     := regexp4.MatchString(action_name)
      if matches {
        action_list = strings.Split(action_name, ",")
      } else {
        action_list = append(action_list, action_name)
      }
      for act_num := 0 ; act_num < len(action_list) ; act_num++ {
        parameter := action_list[act_num]
        message   := "action flag " +parameter
        verbose_message(message, "process")
        _, exists := arguments[parameter]
        if exists {
          matches := regexp6.MatchString(arguments[parameter].category)
          if (matches) {
            arguments[parameter].function()
          } else {
            options["help"] = "all"
            help()
          }
        } else {
          options["help"] = "all"
          help()
        }
      }
    }
  }
  os.Exit(0)
}
