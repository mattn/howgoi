package main

import (
	"flag"
	"fmt"
	"github.com/mattn/howgoi"
	"os"
)

var pos = flag.Int("p", 1, "select answer in specified question (default: 1)")
var all = flag.Bool("a", false, "display the full text of the answer")
var num = flag.Int("n", 1, "number of answers to return")
var link = flag.Bool("l", false, "display only the answer link")

func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, `
usage: howgoi [-h] [-p POS] [-a] [-l] [-c] [-n NUM_ANSWERS] [-C]
              [QUERY [QUERY ...]]

instant coding answers via the command line

positional arguments:
  QUERY                 the question to answer

optional arguments:
  -h                  show this help message and exit
  -p POS              select answer in specified position (default: 1)
  -a                  display the full text of the answer
  -l                  display only the answer link
  -n NUM_ANSWERS      number of answers to return
`[1:])
	}
	flag.Parse()

	if flag.NArg() == 0 || *pos < 1 {
		flag.Usage()
		os.Exit(0)
	}

	answers, err := howgoi.QueryN(*num, flag.Args()...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err.Error())
		os.Exit(1)
	}
	if len(answers) == 0 {
		fmt.Fprintln(os.Stderr, "Sorry, couldn't find any help with that topic")
		os.Exit(1)
	}
	if *all == false {
		n := *pos - 1
		if n >= len(answers) || n < 0 {
			n = 0
		}
		answers = answers[n : n+1]
	}
	for _, answer := range answers {
		if *link {
			fmt.Println(answer.Link)
		} else {
			fmt.Print(answer.Code)
			if len(answers) > 1 {
				fmt.Println()
			}
		}
	}
}
