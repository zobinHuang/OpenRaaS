package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
)

/*
	func: RunShell
	description: run a target shell with parameters on the host
	return: a channel instance used for closing the terminal
*/
func RunShell(execCmd string, params []string) chan struct{} {
	// fmt.Printf("pwd: %s\n", os.Getenv("PWD"))
	fmt.Println("\n==================\nRunShell: ", execCmd, params)

	cmd := exec.Command(execCmd, params...)

	cmd.Env = os.Environ()
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			fmt.Printf("%s", scanner.Text())
		}
	}()

	err = cmd.Start()
	if err != nil {
		fmt.Printf("ERROR: cmd %s fail, %v\n", execCmd, err)
		return nil
	}
	fmt.Printf("Done running cmd %s\n", execCmd)
	err = cmd.Wait()
	if err != nil {
		fmt.Printf("ERROR: cmd %s fail, %v\n", execCmd, err)
		return nil
	}

	done := make(chan struct{})
	// clean up func
	go func() {
		<-done // kill process when data coming
		err := cmd.Process.Kill()
		fmt.Printf("Kill cmd %s error: %v\n==================\n", execCmd, err)
	}()
	return done
}

/*
	func: RunShellWithReturn
	description: run a target shell with parameters on the host, and wait until it finished
	return: logs from the shell
*/
func RunShellWithReturn(execCmd string, params []string) string {
	// fmt.Printf("pwd: %s\n", os.Getenv("PWD"))
	fmt.Println("\n==================\nRunShellWithReturn: ", execCmd, params)

	cmd := exec.Command(execCmd, params...)
	cmd.Env = os.Environ()

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run() // cmd.Run() will block until it finishes
	if err != nil {
		fmt.Printf("ERROR: cmd %s fail, %v\n", execCmd, err)
		return ""
	}

	ret := stdout.String()
	fmt.Println(ret)

	fmt.Printf("Done running cmd %s with return\n==================\n", execCmd)

	return ret
}
