package main

import (
	// "github.com/spf13/cobra"

	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	// "os"
	"flag"
)

var debug = false

func main() {
	debugFlag := flag.Bool("debug", false, "Display debugging messages")

	flag.Parse()

	debug = *debugFlag

	// read shortcuts and remove old full dates
	//
	// var outbuf bytes.Buffer
	alloutput, err := readShortcuts()
	if err != nil {
		log.Fatal(fmt.Errorf("Error reading shortcuts: %w", err))
	}
	//alloutput := outbuf.String()

	if debug {
		//fmt.Printf("outbuf: %s\n", outbuf)
		fmt.Printf("alloutput: %s\n", alloutput)
	}

	if alloutput == "" {
		log.Fatal(fmt.Errorf("Error reading shortcuts, empty output"))
	}

	// add new shortcuts for today
	//
	now := time.Now()
	today_formatted := now.Format("2006-01-02")

	if debug {
		fmt.Printf("today_formatted: %s\n", today_formatted)
	}

	lines := strings.Split(alloutput, "\n")
	// fmt.Printf("lines: %v\n", lines)
	// fmt.Printf("len(lines): %v\n", len(lines))
	// 1. "input" -> "output"
	// we can ignore nested quotes since we only care about date shortcuts
	deletedShortcutCount := 0
	for _, line := range lines {
		tokens := strings.Split(line, "\"")
		if debug {
			fmt.Printf("tokens: %v\n", tokens)
			fmt.Printf("len(tokens): %v\n", len(tokens))
		}
		if len(tokens) == 5 {
			input := tokens[1]
			// output := tokens[3]
			if debug {
				log.Printf("input: %v\n", input)
			}
			// log.Printf("output: %v\n", output)
			// cheap pattern check
			if input > "0000-00-00" && input < "9999-99-99" && input != today_formatted {
				deletedShortcutCount += 1
				log.Printf("Deleting shortcuts: %v\n", line)
				if _, err := deleteShortcut(input); err != nil {
					log.Fatal(fmt.Errorf("Error deleting date shortcut: %w", err))
				}
			}
		}
	}
	if debug {
		log.Printf("Deleted %d shortcuts\n", deletedShortcutCount)
	}
	// datestringLength := len(datestring)
	// shortcutString := datestring[:datestringLength-1]
	// log.Printf("now: %s\n", now)
	// log.Printf("datestring: %s\n", datestring)
	// log.Printf("datestringLength: %d\n", datestringLength)
	// log.Printf("shortcutString: %s\n", shortcutString)

	shortcuts := []string{"dth", "Dth", ""}

	shortcuts[2] = today_formatted

	for _, shortcut := range shortcuts {
		if debug {
			log.Printf("Updating shortcut %v to %v\n", shortcut, today_formatted)
		}
		if _, err := updateShortcut(shortcut, today_formatted); err != nil {
			log.Fatal(fmt.Errorf("Error updating shortcut: %w", err))
		}
	}

}

func deleteShortcut(input string) (string, error) {
	// don't care about output
	return runShortcuts("delete", input, "")
}
func updateShortcut(input, output string) (string, error) {
	return runShortcuts("update", input, output)
}
func readShortcuts() (string, error) {
	return runShortcuts("read", "", "")
}
func runShortcuts(runtype, input, output string) (string, error) {
	args := make([]string, 1, 3)
	args[0] = runtype
	if input != "" {
		args = append(args, input)
	}
	if output != "" {
		args = append(args, output)
	}
	if debug {
		fmt.Printf("args: %v\n", args)
		fmt.Printf("len(args): %v\n", len(args))
	}
	cmd := exec.Command("/opt/homebrew/bin/shortcuts", args...)
	stdoutStderr, err := cmd.CombinedOutput()
	if debug {
		log.Printf("cmd: %s\n", cmd)
		log.Printf("stdoutStderr: %s\n", stdoutStderr)
	}
	return string(stdoutStderr), err
}
