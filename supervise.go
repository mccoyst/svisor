// © 2014 Steve McCoy. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"os"
	"os/exec"
)

type S struct {
	add    chan string
	remove chan string
	deaths chan string
	stop   chan bool
	kids   map[string]*exec.Cmd
}

func New() *S {
	return &S{
		add:    make(chan string),
		remove: make(chan string),
		deaths: make(chan string),
		stop:   make(chan bool),
		kids:   make(map[string]*exec.Cmd),
	}
}

func (s *S) Supervise() {
	for {
		select {
		case prog := <-s.deaths:
			if s.kids[prog] != nil {
				s.spawn(prog)
			}
		case prog := <-s.add:
			s.spawn(prog)
		case prog := <-s.remove:
			delete(s.kids, prog)
		case <-s.stop:
			return
		}
	}
}

func (s *S) Add(prog string) error {
	if prog == "" {
		return errors.New("svisor: program name must be nonempty")
	}
	s.add <- prog
	return nil
}

func (s *S) Remove(prog string) {
	s.remove <- prog
}

func (s *S) Stop() {
	s.stop <- true
}

func (s *S) spawn(prog string) {
	c := exec.Command(prog)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	s.kids[prog] = c

	go func() {
		err := c.Run()
		if err != nil {
			os.Stderr.WriteString(c.Path + " died: " + err.Error() + "\n")
		} else {
			// os.Stderr.WriteString(c.Path + " exited normally\n")
		}
		s.deaths <- prog
	}()
}
