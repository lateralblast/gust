// 2>/dev/null; e=$(mktemp); go build -o $e "$0"; $e "$@" ; r=$?; rm $e; exit $r

/*
Name:         gust (Golang Universal Shell script Template)
Version:      0.0.1
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
  "runtime"
  "strings"
  "regexp"
  "bufio"
  "fmt"
  "os"
)

// Create a structure to manage commandline arguments/switches

type Argument struct {
  name      string
  long      string
  short     string
  info      string
  category  string
  value     string
}

// Initialize variables

var (
  // Create booleans for actions and options switches
  do_actions = false
  do_options = false
  // Create a map to store default booleans for options
  // This should contain a default value for each option created
  defaults = map[string]bool{
    "verbose": false,
    "force":   false,
    "dryrun":  false,
  }
  options = defaults
  // Create a map of Argument structs to store commanline argument information
  // This should include both the short forms (e.g. -V) and long version (e.g. --version) 
  arguments = map[string]Argument {
    "action": {
      info:     "Perform action",
      short:    "a",
      long:     "action",
      category: "switch",
      value:    "",
    },
    "a": {
      info:     "Perform action",
      short:    "a",
      long:     "action",
      category: "switch",
      value:    "",
    },
    "dryrun": {
      info:     "Enable dryrun mode",
      short:    "d",
      long:     "dryrun",
      category: "option",
      value:    "false",
    },
    "d": {
      info:     "Enable dryrun mode",
      short:    "d",
      long:     "dryrun",
      category: "option",
      value:    "false",
    },
    "help": {
      info:     "Print help information",
      short:    "h",
      long:     "help",
      category: "switch",
      value:    "",
    },
    "h": {
      info:     "Print help information",
      short:    "h",
      long:     "help",
      category: "switch",
      value:    "",
    },
    "verbose": {
      info:     "Enable verbose output",
      short:    "v",
      long:     "verbose",
      category: "option",
      value:    "false",
    },
    "v": {
      info:     "Enable verbose output",
      short:    "v",
      long:     "verbose",
      category: "option",
      value:    "false",
    },
    "version": {
      info:     "Print version information",
      short:    "V",
      long:     "version",
      category: "switch",
      value:    "",
    },
    "V": {
      info:     "Print version information",
      short:    "V",
      long:     "version",
      category: "switch",
      value:    "",
    },
  }
)

/*
Funtion:      verbose_message
Parameters:   message and formet
Description:  A routine to create consistently formatted output
*/

func verbose_message(message, format string) {
  var header = ""
  format = strings.ToLower(format)
  format = strings.Title(format)
  matches, _ := regexp.MatchString("verbose", format)
  if matches {
    fmt.Println(message) 
  } else {
    if (options["verbose"]) {
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
Funtion:      warning_message
Parameters:   message
Description:  A routine to display a warning, overriding non verbose mode if needed
*/

func warning_message(message string) {
  if (options["verbose"]) {
    verbose_message(message, "warn")
  } else {
    options["verbose"] = true
    verbose_message(message, "warn")
    options["verbose"] = false
  }
}

/*
Funtion:      print_help_category
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
        if len(argument.long) < 15 {
          fmt.Printf("%s, %s:\t\t%s\n", argument.long, argument.short, argument.info)
        } else {
          fmt.Printf("%s, %s:\t%s\n", argument.long, argument.short, argument.info)
        }
      }
    }
  }
  fmt.Println("")
}


/*
Funtion:      print_help
Parameters:   help_flags
Description:  A routine to print help information
*/

func print_help(help_flags string) {
  switch help_flags {
    case "option", "options":
      print_help_category("option")
    case "switch", "switches":
      print_help_category("switch")
    case "all":
      print_help_category("switch")
      print_help_category("option")
  }
  os.Exit(0)
}

/*
Funtion:      print_version
Parameters:   script_file
Description:  A routine to print version information
*/

func print_version(script_file string) {
  open_file, file_error := os.Open(script_file)
  if file_error != nil {
      fmt.Println(file_error)
  }
  defer open_file.Close()
  scanner := bufio.NewScanner(open_file)
  for scanner.Scan() {
    line := scanner.Text()
    if strings.Contains(line, "Version:") {
      matches, _ := regexp.MatchString("[0-9]", line)
      if matches {
        fields := regexp.MustCompile("[^\\s]+").FindAllString(line, -1)
        fmt.Println(fields[1])
      }
    }
  }
  os.Exit(0)
}

/*
Funtion:      handle_options 
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
  for number := 0 ;  number < len(parameters) ; number++ {
    parameter := parameters[number]
    matches, _ := regexp.MatchString("^no", parameter)
    format := ""
    if matches {
      format = "disable"
      parameter = strings.Split(parameter, "no")[1]
      options[parameter] = false
    } else {
      format = "enable"
      options[parameter] = true
    }
    verbose_message(parameter, format)
  }
  return
}

/*
Funtion:      check_value
Parameters:   arg_num
Description:  A routine to handle argument values
*/

func check_value(arg_num int) {
  parameter := os.Args[arg_num]
  if arg_num == len(os.Args)-1 {
    message := "No value given for " + parameter
    switch parameter {
      case "--help", "-h":
        print_help("all")
      default:
        verbose_message(message, "warn")
        print_help("all")
    }
    os.Exit(1)
  }
  check_value := os.Args[arg_num+1] 
  matches, _ := regexp.MatchString("^-", check_value)
  if matches {
    message := "No value given for " + parameter
    options["verbose"] = true
    verbose_message(message, "warn")
    os.Exit(1)
  } else {
    message := "Value given for " + parameter + " " + check_value
    verbose_message(message, "info")
    matches, _ := regexp.MatchString("help|h", parameter)
    if matches {
      value := os.Args[arg_num+1]
      print_help(value)
    }
  }
}

// Main function

func main() {
  // Create arrays to store actions or options
  action_flags := []string{}
  option_flags := []string{}
  // Save CLI arguments and check for verbose option
  cli_args := strings.Join([]string(os.Args), " ")
  matches, _ := regexp.MatchString("noverbose", cli_args)
  if matches {
    options["verbose"] = false
  } else {
    matches, _ := regexp.MatchString("verbose", cli_args)
    if matches {
      options["verbose"] = true
    }
  }
  // Get script file
  _, script_file, _,  _ := runtime.Caller(0)
  // If we have no arguments print help information
  if len(os.Args) < 2 {
    help_flags := "all"
    print_help(help_flags)
  }
  // loop through command line arguments and handle them
  for arg_num := 1 ; arg_num < len(os.Args) ; arg_num++ {
    arg_name := os.Args[arg_num]
    // Convert plural arguments to non plural
    arg_name = strings.Replace(arg_name, "options", "option", -1)
    arg_name = strings.Replace(arg_name, "actions", "action", -1)
    // Check if we have a -abc style switch and process
    matches, _ := regexp.MatchString("^-[a-z,A-Z][a-z,Z]", arg_name)
    if matches {
      // Strip -
      arg_names := strings.Split(arg_name, "-")[1]
      // Step though each command line arguement, e.g. -abc > a, b, c,
      letters   := strings.Split(arg_names, "")
      for num :=0 ; num < len(letters) ; num++ {
        letter := letters[num] 
        _, exists := arguments[letter]
        if (exists) {
          // Check that an argument structure exists and grab the long version
          matches, _ := regexp.MatchString("option", arguments[letter].category)
          if matches {
            long_name := arguments[letter].long
            handle_options(long_name)
          }
        } else {
          // Print help if there is no argument structure
          message := "Commandline argument "+letter+" does not exist"
          warning_message(message)
          print_help("all")
        }
      }
    } else {
      // Strip -
      arg_name = strings.Replace(arg_name, "-", "", -1)
      // Check argument structure exists
      _, exists := arguments[arg_name]
      if exists {
        // If argument structure exists check if it is an option and handle
        matches, _ := regexp.MatchString("option", arguments[arg_name].category)
        long_name := arguments[arg_name].long
        if matches {
          handle_options(long_name)
        } else {
          // If argument is not an option, handle appropriatle
          switch long_name {
            case "action":
              check_value(arg_num)
              action_flags = append(action_flags, os.Args[arg_num+1])
              do_actions = true
            case "option":
              check_value(arg_num)
              option_flags = append(option_flags, os.Args[arg_num+1])
              do_options = true
            case "version":
              print_version(script_file)
            case "help":
              check_value(arg_num)
            default:
              print_help("all")
          }
        }
      } else {
        // check if argument is a negative option, e.g. noverbose and handle
        matches, _ := regexp.MatchString("^no", arg_name)
        if matches {
          parameter := strings.Split(arg_name, "no")[1]
          matches, _ := regexp.MatchString("option", arguments[parameter].category)
          if matches {
            handle_options(arg_name)
          } else {
            // If argument structure does exist warn and print help
            message := "Commandline argument "+arg_name+" does not exist"
            warning_message(message)
            print_help("all")
          }
        } else {
          // If argument structure does exist warn and print help
          message := "Commandline argument "+arg_name+" does not exist"
          warning_message(message)
          print_help("all")
        }
      }
    }
  } 
  // If we have option(s) handle each
  if do_options {
    for number := 0 ; number < len(option_flags) ; number++ {
      values := option_flags[number]
      handle_options(values)
    }
  }
  // If we have action(s) handle each
  if do_actions {
    for number := 0 ; number < len(action_flags) ; number++ {
      actions := []string{}
      action := action_flags[number]
      matches, _ := regexp.MatchString(",", action)
      if matches {
        actions = strings.Split(action, ",")
      } else {
        actions = append(actions, action)
      }
      for act_num := 0 ; act_num < len(actions) ; act_num++ {
        message := "action flag " +actions[act_num]  
        verbose_message(message, "process")
      }
    }
  }
  os.Exit(0)
}
