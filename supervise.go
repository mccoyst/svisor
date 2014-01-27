// © 2014 Steve McCoy. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package svisor

import (
	"errors"
	"io"
	"os/exec"
)

type S struct {
	add    chan string
	remove chan string
	deaths chan string
	stop   chan bool
	kids   map[string]*exec.Cmd
	log	io.Writer
	Stdout	io.Writer
	Stderr	io.Writer
	Stdin	io.Reader
}

// New returns an initialized S, which will log to w.
// Stdout, Stderr, and Stdin are used to set the
// the same-named variables in supervised subcommands,
// and all of them default to nil.
func New(w io.Writer) *S {
	return &S{
		add:    make(chan string),
		remove: make(chan string),
		deaths: make(chan string),
		stop:   make(chan bool),
		kids:   make(map[string]*exec.Cmd),
		log: w,
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
	c.Stdout = s.Stdout
	c.Stderr = s.Stderr
	c.Stdin = s.Stdin
	s.kids[prog] = c

	go func() {
		err := c.Run()
		if err != nil {
			io.WriteString(s.log, c.Path + " died: " + err.Error() + "\n")
		} else {
			io.WriteString(s.log, c.Path + " exited normally\n")
		}
		s.deaths <- prog
	}()
}
