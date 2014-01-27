// © 2014 Steve McCoy. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"time"
	"os"

	"svisor"
)

func main() {
	s := svisor.New(os.Stderr)
	s.Stdout = os.Stdout
	s.Stderr = os.Stderr
	go s.Supervise()
	s.Add("./hi.sh")
	s.Add("./bye.sh")
	s.Add("./dne")

	time.Sleep(1 * time.Second)
	s.Remove("./hi.sh")

	time.Sleep(1 * time.Second)
	s.Remove("./bye.sh")

	time.Sleep(1 * time.Second)
	s.Stop()
}
