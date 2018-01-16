package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const backlightDirectory = "/sys/class/backlight/intel_backlight"

var (
	maxBrightnessFilename = filepath.Join(backlightDirectory, "max_brightness")
	brightnessFilename    = filepath.Join(backlightDirectory, "brightness")
)

func main() {
	inc := flag.Int("inc", 0, "increment [%]")
	dec := flag.Int("dec", 0, "decrement [%]")
	flag.Parse()

	curr, err := readBrightness()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot read brightness value, err: %v\n", err)
		os.Exit(1)
	}
	max, err := readMaxBrightness()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot read max brightness value, err: %v\n", err)
		os.Exit(1)
	}

	switch {
	case *inc > 0:
		step := (*inc * max) / 100
		value := curr + step
		if value > max {
			value = max
		}
		if err := writeBrightness(value); err != nil {
			fmt.Fprintf(os.Stderr, "Cannot write brightness value, err: %v\n", err)
			os.Exit(1)
		}
	case *dec > 0:
		step := (*dec * max) / 100
		value := curr - step
		if value < 0 {
			value = 0
		}
		if err := writeBrightness(value); err != nil {
			fmt.Fprintf(os.Stderr, "Cannot write brightness value, err: %v\n", err)
			os.Exit(1)
		}
	}
}

func readMaxBrightness() (int, error) {
	return readBrightnessFile(maxBrightnessFilename)
}

func readBrightness() (int, error) {
	return readBrightnessFile(brightnessFilename)
}

func readBrightnessFile(filename string) (int, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(strings.TrimSpace(string(b)))
}

func writeBrightness(v int) error {
	return ioutil.WriteFile(brightnessFilename, []byte(strconv.Itoa(v)), 0644)
}
