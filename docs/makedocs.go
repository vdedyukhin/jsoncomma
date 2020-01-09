package main

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"os"
	"os/exec"
)

// +build ignore

// builds the documentation

func main() {
	vars := map[string]string{
		"Version":             getVersion(),
		"JsonCommaServerHelp": getHelp(),
	}
	tmpl := template.Must(template.New("index.html.template").ParseFiles("./docs/index.html.template"))
	f, err := os.Create("./docs/index.html")
	if err != nil {
		log.Fatal(err)
	}
	if err := tmpl.Execute(f, vars); err != nil {
		log.Fatal(err)
	}
}

func getVersion() string {
	log.Printf("getting version from git")
	output, err := exec.Command("git", "describe", "--tags").Output()

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		log.Print("Stderr:")
		fmt.Println(string(exitErr.Stderr))
		log.Fatal(err)
	} else if err != nil {
		log.Fatal(err)
	}

	return string(bytes.TrimRight(output, "\n"))
}

func getHelp() string {
	log.Printf("getting help message from jsoncomma")
	output, err := exec.Command("./jsoncomma", "server", "-help").Output()

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		if exitErr.ExitCode() == 2 {
			output = exitErr.Stderr
		} else {
			log.Print("Stderr:")
			fmt.Println(string(exitErr.Stderr))
			log.Fatal(err)
		}
	} else if err != nil {
		log.Fatal(err)
	}

	return string(output)
}
