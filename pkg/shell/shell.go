package shell

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"

	"github.com/kballard/go-shellquote"
)

// Command describes the command object
type Command struct {
	*exec.Cmd
	StdOutLog bool
}

// Option sets the state of an option for a command
type Option func(*Command)

// WithCommandLog prints the full command
func WithCommandLog() Option {
	return func(cmd *Command) {
		fmt.Println(cmd.String())
	}
}

// WithStdOutLog makes the stdout of a command visible
func WithStdOutLog() Option {
	return func(cmd *Command) {
		cmd.StdOutLog = true
	}
}

// NewCommand builds and returns a command to be executed
func NewCommand(command string, options ...Option) *Command {
	args, err := shellquote.Split(command)
	if err != nil {
		log.Fatalln(err)
	}
	if len(args) == 0 {
		log.Fatalf("No args for command %s\n", command)
	}

	cmd := Command{exec.Command(args[0], args[1:]...), false}

	for _, o := range options {
		o(&cmd)
	}

	return &cmd
}

// Run executes the command and streams the results of stdout and stderr into separated channels
func (cmd *Command) Run() error {
	outputBuff := strings.Builder{}
	stdoutChan := make(chan string)
	stderrChan := make(chan string)
	errChan := make(chan error)

	go cmd.run(stdoutChan, stderrChan, errChan)

outputCapture:
	for {
		select {
		case s, ok := <-stdoutChan:
			if cmd.StdOutLog {
				if !ok {
					break outputCapture
				}
				fmt.Print(s)
			} else {
				outputBuff.WriteString(s)
			}
		case s, ok := <-stderrChan:
			if !ok {
				break outputCapture
			}
			if cmd.StdOutLog {
				fmt.Print(s)
			} else {
				outputBuff.WriteString(s)
			}
		case err := <-errChan:
			return err
		}
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	statusCode := cmd.ProcessState.ExitCode()
	if statusCode != 0 {
		return fmt.Errorf("%s failed with status [%d], output: %s", strings.Join(cmd.Args, " "), statusCode, outputBuff.String())
	}

	return nil
}

func (cmd *Command) run(stdoutChan, stderrChan chan string, errChan chan error) {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		errChan <- err
		return
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		errChan <- err
		return
	}

	if err := cmd.Start(); err != nil {
		errChan <- err
		return
	}

	go read(stdout, stdoutChan, errChan)
	go read(stderr, stderrChan, errChan)
}

func read(pipe io.ReadCloser, outputChan chan string, errChan chan error) {
	buf := bufio.NewReader(pipe)

	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				close(outputChan)
				return
			}

			errChan <- err
			close(outputChan)
			return
		}

		outputChan <- line
	}
}

// PipeCommands executes commands in pipe ("|") passing the stdout of one process into the stdin of the following one
func PipeCommands(stack ...*Command) (*bytes.Buffer, *bytes.Buffer, error) {
	stdout := bytes.Buffer{}
	stderr := bytes.Buffer{}
	pipeStack := make([]*io.PipeWriter, len(stack)-1)

	i := 0
	for ; i < len(stack)-1; i++ {
		inPipe, outPipe := io.Pipe()
		stack[i].Stdout = outPipe
		stack[i].Stderr = &stderr
		stack[i+1].Stdin = inPipe
		pipeStack[i] = outPipe
	}
	stack[i].Stdout = &stdout
	stack[i].Stderr = &stderr

	if err := pipe(stack, pipeStack); err != nil {
		return &stdout, &stderr, err
	}

	return &stdout, &stderr, nil
}

func pipe(stack []*Command, pipes []*io.PipeWriter) (err error) {
	if stack[0].Process == nil {
		if err = stack[0].Start(); err != nil {
			return err
		}
	}

	if len(stack) > 1 {
		if err = stack[1].Start(); err != nil {
			return err
		}
		defer func() {
			if err == nil {
				if err := pipes[0].Close(); err != nil {
					return
				}
				err = pipe(stack[1:], pipes[1:])
			}
		}()
	}

	return stack[0].Wait()
}
