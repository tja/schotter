package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	svg "github.com/ajstarks/svgo"
	"github.com/fatih/color"
	"github.com/spf13/pflag"
)

// Main entry point.
func main() {
	// Command line
	var (
		row     = pflag.IntP("columns", "c", 10, "Number of columns")
		col     = pflag.IntP("rows", "r", 20, "Number of rows")
		size    = pflag.Float64P("size", "s", 100.0, "Size of each square")
		margin  = pflag.Float64P("margin", "m", 100.0, "Margin around all squares")
		order   = pflag.Float64P("orderliness", "o", 8.0, "Amount of orderlines, higher ist more ordered")
		help    = pflag.BoolP("help", "h", false, "Help for schotter")
		version = pflag.Bool("version", false, "Version of schotter")
	)

	pflag.Parse()

	if *help {
		printBanner()
		fmt.Println("Generate digital line art, inspired by Georg Nees' \"Schotter\".")
		fmt.Println("")
		printHelp()
		os.Exit(0)
	}

	if *version {
		fmt.Println("schotter version 1.1.0")
		os.Exit(0)
	}

	if len(pflag.Args()) != 1 {
		fmt.Println("Error: Output filename missing.")
		os.Exit(2)
	}

	printBanner()

	// Open output file
	out, err := os.Create(pflag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}

	defer out.Close()

	// Create SVG canvas
	canvas := svg.New(out)
	canvas.Start(
		int(2.0*(*margin)+float64(*col)*(*size)),
		int(2.0*(*margin)+float64(*row)*(*size)),
	)

	defer canvas.End()

	// Generate random Schotter art
	rand.Seed(time.Now().UnixNano())

	for y := 0; y < (*row); y++ {
		factor := float64(y*y+y) / (*order)

		for x := 0; x < (*col); x++ {
			// Create randomness
			var (
				offsetX = (rand.Float64()*2 - 1) * factor / 100.0 // Random horizontal offset
				offsetY = (rand.Float64()*2 - 1) * factor / 100.0 // Random vertical offset
				rotate  = (rand.Float64()*2 - 1) * factor / 2.0   // Random rotation
			)

			// Draw rotated and translated square
			canvas.TranslateRotate(
				int((*margin)+(float64(x)+offsetX)*(*size)),
				int((*margin)+(float64(y)+offsetY)*(*size)),
				rotate,
			)

			canvas.Rect(
				0,
				0,
				int(*size),
				int(*size),
				"fill:none;stroke:black",
			)

			canvas.Gend()
		}
	}

	fmt.Println("Successfully created output file", pflag.Arg(0))
}

// Print banner
func printBanner() {
	color.HiCyan("             __          __   __              ")
	color.HiCyan(".-----.-----|  |--.-----|  |_|  |_.-----.-.--.")
	color.HiCyan("|__ --|  ---|  .  |  -  |   _|   _|  -__|  .-'")
	color.HiCyan("|_____|_____|__|__|_____|____|____|_____|__|  ")
	color.HiCyan("                                              ")
}

// Print help screen
func printHelp() {
	fmt.Println("Usage:")
	fmt.Println("  schotter [flags] file")
	fmt.Println("")
	fmt.Println("Flags:")

	pflag.PrintDefaults()
}
