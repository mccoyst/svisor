// © 2014 Steve McCoy. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "time"

func main() {
	s := New()
	go s.Supervise()
	s.Add("./hi.sh")
	s.Add("./bye.sh")

	time.Sleep(1 * time.Second)
	s.Remove("./hi.sh")

	time.Sleep(1 * time.Second)
	s.Remove("./bye.sh")

	time.Sleep(1 * time.Second)
	s.Stop()
}
