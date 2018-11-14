package main

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/dtylman/connan/db"
	"github.com/gobwas/glob"
)

//Condition a condition in which analyzer will work
type Condition struct {
	Field      string `json:"field"`
	Expression string `json:"expression"`
	g          glob.Glob
}

//Match matches the doc with the field glob expression
func (c *Condition) Match(doc *db.Document) (bool, error) {
	var err error
	if c.g == nil {
		c.g, err = glob.Compile(c.Expression)
		if err != nil {
			return false, err
		}
	}
	return c.g.Match(doc.GetField(c.Field)), nil
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
			return fmt.Errorf("Condition '%v' on field '%v' failed on doc '%v': '%v'", cond.Expression, cond.Field,
				doc.Path, err)
		}
		if !m {
			log.Printf("Document '%v' does not match condition '%v' on field '%v'", doc.Path, cond.Expression, cond.Field)
			return nil
		}
	}
	val := ca.execute(path)
	doc.SetField(ca.PopulateField, val)
	return nil
}

//Name ...
func (ca CommandAnalyzer) Name() string {
	return ca.AnalyzerName
}

func (ca *CommandAnalyzer) execute(path string) string {
	line := strings.Replace(ca.Command, "[path]", path, -1)
	log.Printf("Executing '%v'", line)

	command := strings.Fields(line)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(ca.Timeout)*time.Second)
	defer cancel() // The cancel should be deferred so resources are cleaned up

	cmd := exec.CommandContext(ctx, command[0], command[1:]...)

	out, err := cmd.Output()
	if err != nil {
		log.Printf("'%v' failed: '%v'", line, err)
	}
	if ctx.Err() == context.DeadlineExceeded {
		log.Printf("'%v' timeout: '%v'", line, ctx.Err())
	}
	return string(out)
}
