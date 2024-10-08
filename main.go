package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
)

func generateRuleset(rule int) map[byte]byte {
	br := byte(rule)

	if rule > 255 {
		fmt.Printf("Rule exceeds byte capacity - truncated\n")
	}

	ruleset := make(map[byte]byte, 8)

	// 1 0 0 0 0 0 0 0
	// Then we divide by 2 to get each bit
	checkvalue := 128

	// The number of bits we left shift
	lsft := 7

	// We create a map:
	// key = value of 3 bits
	// value = the value of the bit in the original rule that is found at position `key`
	for i := 7; i >= 0; i-- {
		ruleset[byte(i)] = br & byte(checkvalue) >> byte(lsft)
		checkvalue /= 2
		lsft--
	}

	return ruleset
}

func initBuf(size int, config string, hatstart bool, showFor1 string) []byte {
	conflen := len(config)

	if conflen != size && config != "" {
		fmt.Printf("Setting the buffer size to %d (config value size) instead of %d (passed value)\n\n", conflen, size)
		size = conflen
	}

	buf := make([]byte, size)
	defer printBuf(buf, showFor1)

	// If we do a hat start, we return early
	if hatstart {
		buf[size/2] = 1
		return buf
	}

	for i := 0; i < size; i++ {
		if config != "" {
			// We get rid of the ASCII value and transform it into number
			buf[i] = config[i] - 48
		} else {
			// We set random value between 0 and 1
			buf[i] = byte(rand.Intn(2))
		}
	}

	return buf
}

func updateBuf(buf []byte, ruleset map[byte]byte, showFor1 string) {
	bufsize := len(buf)

	// We create the simulated buffer
	updated := make([]byte, bufsize)

	// The 3 bits
	l := byte(0)
	c := byte(0)
	r := byte(0)

	bitval := byte(0)

	for i := 0; i < bufsize; i++ {
		switch i {
		case 0:
			l = buf[bufsize-1]
			c = buf[i]
			r = buf[i+1]
		case bufsize - 1:
			l = buf[bufsize-2]
			c = buf[bufsize-1]
			r = buf[0]
		default:
			l = buf[i-1]
			c = buf[i]
			r = buf[i+1]
		}

		// We get the 3-bits value
		bitval = byte(l<<2 | c<<1 | r)

		// We look for value in the map
		updated[i] = ruleset[bitval]
	}

	// We update the original buffer
	for i := 0; i < bufsize; i++ {
		buf[i] = updated[i]
	}

	printBuf(buf, showFor1)

}

// The terminal codes to clear the screen
func clearScreen(stackPrint bool) {
	if !stackPrint {
		fmt.Print("\033[H\033[2J")
	}
}

func printBuf(buf []byte, showFor1 string) {
	bufsize := len(buf)

	for i := 0; i < bufsize; i++ {
		if showFor1 != "" {
			if buf[i] != 0 {
				fmt.Printf("%s", showFor1)
			} else {
				fmt.Printf(" ")
			}
		} else {
			fmt.Printf("%d ", buf[i])
		}
	}
	fmt.Printf("\n")
}

func main() {
	rule := flag.Int("r", 0, "The rule of the Elementary Cellular Automaton (default 0)")
	bufsize := flag.Int("w", 11, "The size of the buffer")
	stackPrint := flag.Bool("s", false, "If passed, it will not clear the previous buffer output")
	showFor1 := flag.String("S", "", "Print this value instead of 1, and hide all 0")
	config := flag.String("c", "", "The initial config of the cells (default random)")
	tbp := flag.Int("t", 2000, "The time in milliseconds taken between each print")
	hatstart := flag.Bool("h", false, "A configuration where the middle cell is 1, the rest are 0 - requires an odd buffer size")
	hatrule := flag.Bool("H", false, "The ruleset stated by the Hat Rule (rule 18) - it will overwrite the rule")

	flag.Parse()

	if *hatrule {
		*rule = 18
	}

	if *tbp <= 0 {
		fmt.Printf("Cannot have less than 0 ms of time between prints")
	}

	if *bufsize <= 0 {
		fmt.Printf("Cannot have a buffer of size 0 or less\n")
		return
	}

	if *rule <= 0 {
		fmt.Printf("Cannot have a rule smaller or equal to 0\n")
		return
	}

	if *hatstart && *bufsize%2 == 0 {
		fmt.Printf("Cannot have a Hat start without a buffer of odd size\n")
		return
	}

	rs := generateRuleset(*rule)
	buf := initBuf(*bufsize, *config, *hatstart, *showFor1)

	for {
		clearScreen(*stackPrint)
		updateBuf(buf, rs, *showFor1)
		time.Sleep(time.Duration(*tbp) * time.Millisecond)
	}

}
