package main

import (
	"flag"
	"fmt"
	"os"
)

const MYNUMBER_MAX = 100_000_000_000

type Config struct {
	Output          string
	From, To, Count int
	Debug, Quiet    bool
}

func (c Config) GetEnd() int {
	if c.Count != 0 {
		return c.From + c.Count
	}

	return c.To
}

func (c Config) validate() error {
	begin := c.From
	end := c.GetEnd()

	if begin < 0 || begin >= MYNUMBER_MAX ||
		end < 0 || end >= MYNUMBER_MAX ||
		begin >= end {
		return fmt.Errorf("%s", "out of range")
	}

	return nil
}

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

func generate_mynumber(conf Config) error {
	file, err := os.Create(conf.Output)
	if err != nil {
		return fmt.Errorf("Failed: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	limit := conf.GetEnd()
	for target := conf.From; target < limit; target++ {
		digit_string := fmt.Sprintf("%011d", target)
		var digits [11]int

		for idx, char := range digit_string {
			digits[10-idx] = int(char) - '0'
		}

		mynumber := fmt.Sprintf("%s%d\n", digit_string, calc_digits(digits[:]))
		file.Write([]byte(mynumber))

		if conf.Debug {
			fmt.Print(mynumber)
		}

		if !conf.Quiet && target%10000 == 0 {
			fmt.Print("\rGenerated: ", target, " / ", limit)
		}
	}
	fmt.Println()
	return nil
}

func getConfig() (*Config, error) {
	var (
		write_file_name = flag.String("o", "mynumber-list.txt", "output file name (default mynumber-list.txt)")
		from            = flag.Int("f", 0, "from number (default 0)")
		to              = flag.Int("t", 100, "to number (default 100, limit 99,999,999,999)")
		count           = flag.Int("c", 0, "generate count (default null, limit 99,999,999,999)")
		debug           = flag.Bool("d", false, "debug")
		quiet           = flag.Bool("q", false, "quiet")
	)
	flag.Parse()

	conf := new(Config)
	conf.Debug = *debug
	conf.Quiet = *quiet
	conf.From = *from
	if *count == 0 {
		conf.To = *to
	} else {
		conf.Count = *count
		conf.To = 0
	}
	conf.Output = *write_file_name
	if err := conf.validate(); err != nil {
		return nil, fmt.Errorf("%s", "Invalid option combination")
	}
	return conf, nil
}

func main() {
	conf, err := getConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := generate_mynumber(*conf); err != nil {
		fmt.Println("Generation failed.")
		return
	}
}
