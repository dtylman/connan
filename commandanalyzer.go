package main

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/dtylman/connan/db"
)

//Condition a condition in which analyzer will work
type Condition struct {
	Field   string `json:"field"`
	Pattern string `json:"regexp"`
}

//Match returns true if document matches the condition
func (c *Condition) Match(doc *db.Document) (bool, error) {
	return regexp.MatchString(c.Pattern, doc.GetField(c.Field))
}

const defaultTimeout = 60

//CommandAnalyzer is an analyzer that executes a cli
type CommandAnalyzer struct {
	AnalyzerName  string      `json:"name"`
	Command       string      `json:"command"`
	Timeout       int         `json:"timeout"`
	PopulateField string      `json:"popluates"`
	Conditions    []Condition `json:"conditions"`
}

//Process ...
func (ca CommandAnalyzer) Process(path string, doc *db.Document) error {
	log.Printf("%v: processing '%v'", ca.AnalyzerName, path)
	if ca.Timeout == 0 {
		ca.Timeout = defaultTimeout
	}
	for _, cond := range ca.Conditions {
		m, err := cond.Match(doc)
		if err != nil {
			return fmt.Errorf("Condition '%v' failed on doc '%v': '%v'", cond, doc, err)
		}
		if !m {
			log.Printf("Document '%v' does not match condition '%v'", doc.Path, cond)
			return nil
		}
	}
	val, err := ca.execute(path)
	log.Printf("Error:%v output: %v", err, val)
	if err != nil {
		return err
	}
	doc.SetField(ca.PopulateField, val)
	return nil
}

//Name ...
func (ca CommandAnalyzer) Name() string {
	return ca.AnalyzerName
}

func (ca *CommandAnalyzer) execute(path string) (string, error) {
	command := strings.Fields(ca.Command)
	for i := range command {
		if command[i] == "[path]" {
			command[i] = path
		}
	}
	line := strings.Join(command, " ")
	var out []byte
	var err error
	log.Printf("exec: %v: (%v) '%v'", line, string(out), err)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(ca.Timeout)*time.Second)
	defer cancel() // The cancel should be deferred so resources are cleaned up

	cmd := exec.CommandContext(ctx, command[0], command[1:]...)

	out, err = cmd.Output()
	if err != nil {
		return "", fmt.Errorf("'%v' failed: '%v'", line, err)
	}
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("'%v' timeout: '%v'", line, ctx.Err())
	}

	if cmd.ProcessState != nil && !cmd.ProcessState.Success() {
		return "", fmt.Errorf("Process failed: %v", string(out))
	}
	return string(out), nil
}
