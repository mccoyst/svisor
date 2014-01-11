// © 2014 Steve McCoy.

package main

import (
	"os"
	"os/exec"
)

func supervise(progs []string) {
	cmds := make([]*exec.Cmd, 0, len(progs))
	for _, prog := range progs {
		if len(prog) == 0 {
			continue
		}
		cmds = append(cmds, createCmd(prog))
	}
	deaths := make(chan *exec.Cmd)
	for _, cmd := range cmds {
		spawn(cmd, deaths)
	}
	for {
		select {
		case cmd := <-deaths:
			cmd = &exec.Cmd{
				Path: cmd.Path,
				Stdout: os.Stdout,
				Stderr: os.Stderr,
			}
			spawn(cmd, deaths)
		}
	}
}

func createCmd(prog string) *exec.Cmd {
	c := exec.Command(prog)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c
}

func spawn(c *exec.Cmd, deaths chan *exec.Cmd) {
	go func(){
		err := c.Run()
		if err != nil {
			os.Stderr.WriteString(c.Path + " died: " + err.Error() + "\n")
		} else {
			os.Stderr.WriteString(c.Path + " exited normally\n")
		}
		deaths <- c
	}()
}
