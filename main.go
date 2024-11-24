package main

import (
	"flag"
	"fmt"
	"os"
)

func qn(n int) int {
	if 1 <= n && n <= 6 {
		return n + 1
	} else {
		return n - 5
	}
}

func calc_digits(input []int) int {
	sum := 0
	for idx, pn := range input {
		sum = sum + (pn * qn(idx+1))
	}

	if remain := sum % 11; remain < 2 {
		return 0
	} else {
		return 11 - remain
	}
}

func generate_mynumber(fn string, from int, to int, debug bool) error {
	if from < 0 || from >= 100_000_000_000 ||
		to < 0 || to >= 100_000_000_000 ||
		from >= to {
		return nil
	}

	file, err := os.Create(fn)
	if err != nil {
		return fmt.Errorf("Failed: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	for target := from; target < to; target++ {
		digit_string := fmt.Sprintf("%011d", target)
		var digits [11]int

		for idx, char := range digit_string {
			digits[10-idx] = int(char) - '0'
		}
		mynumber := fmt.Sprintf("%s%d\n", digit_string, calc_digits(digits[:]))

		if debug {
			fmt.Print(mynumber)
		}
		file.Write([]byte(mynumber))
	}
	return nil
}

func main() {
	var (
		write_file_name = flag.String("o", "mynumber-list.txt", "output file name (default mynumber-list.txt)")
		from            = flag.Int("f", 0, "from number (default 0)")
		to              = flag.Int("t", 100, "to number (default 100, limit 99,999,999,999)")
		debug           = flag.Bool("d", false, "tee debug")
	)
	flag.Parse()

	if err := generate_mynumber(*write_file_name, *from, *to, *debug); err != nil {
		return
	}
}
