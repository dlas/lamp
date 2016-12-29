
package main

import "fmt"
import "os"
import "hw"
import "strconv"


func main() {
	h, _ := hw.NewHW();
	r, _ := strconv.Atoi(os.Args[1]);
	g, _ := strconv.Atoi(os.Args[2]);
	b, _ := strconv.Atoi(os.Args[3]);

	h.SetLEDs(r,g,b);

	inputs := h.ReadButtons();

	fmt.Fprintf(os.Stdout, "INPUT: %v\n", inputs);



}

