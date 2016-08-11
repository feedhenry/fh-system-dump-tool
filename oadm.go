package main

import (
	"os/exec"
)

// OadmData is a task factory for tasks that fetch the output of oc admin diagnostics.
// The task uses outFor and errOutFor to get io.Writers to write, respectively, the
// stdout output and stderr output.
func OadmData(outFor, errOutFor platformResourceWriterCloserFactory) Task {
	return oadmData(func() *exec.Cmd {
		return exec.Command("oc", "adm", "diagnostics")
	}, outFor, errOutFor)
}

// A getOcAdmDiagnosticsCmdFactory generates commands to get the raw oc admin diagnostics
// data.
type getOcAdmDiagnosticsCmdFactory func() *exec.Cmd

// oadmData will dump the output of 'oc admin diagnostics' into the stdout writer any stderr
// output is sent to the stderr writer
func oadmData(cmdFactory getOcAdmDiagnosticsCmdFactory, outFor, errOutFor platformResourceWriterCloserFactory) Task {
	return func() error {
		stdOut, stdOutCloser, err := outFor("diagnostics")
		if err != nil {
			return err
		}
		defer stdOutCloser.Close()

		stdErr, stdErrCloser, err := errOutFor("diagnostics")
		if err != nil {
			return err
		}
		defer stdErrCloser.Close()

		cmd := cmdFactory()

		if err := runCmdCaptureOutput(cmd, stdOut, stdErr); err != nil {
			stdErr.Write([]byte(err.Error()))
			return err
		}

		return nil
	}
}
