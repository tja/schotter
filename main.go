package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/ajstarks/svgo"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Main entry point.
func main() {
	// Print banner
	color.HiCyan("             __          __   __              ")
	color.HiCyan(".-----.-----|  |--.-----|  |_|  |_.-----.-.--.")
	color.HiCyan("|__ --|  ---|  .  |  -  |   _|   _|  -__|  .-'")
	color.HiCyan("|_____|_____|__|__|_____|____|____|_____|__|  ")
	color.HiCyan("                                              ")

	// Command definition
	cmd := &cobra.Command{
		Use:     "schotter",
		Long:    "Generate digital line art, inspired by Georg Nees' \"Schotter\".",
		Args:    cobra.ExactArgs(1),
		Version: "0.0.1",
		Run:     schotter,
	}

	cmd.Flags().IntP("columns", "c", 10, "Number of columns")
	cmd.Flags().IntP("rows", "r", 20, "Number of rows")
	cmd.Flags().Float64P("size", "s", 100.0, "Size of each square")
	cmd.Flags().Float64P("margin", "m", 100.0, "Margin around all squares")
	cmd.Flags().Float64P("orderliness", "o", 8.0, "Amount of orderlines, higher ist more ordered")

	// Read viper config
	viper.BindPFlags(cmd.Flags())

	viper.SetConfigName("config")
	viper.AddConfigPath("$HOME/.config/schotter")
	viper.AddConfigPath(".")

	viper.ReadInConfig()

	// Run command
	cmd.Execute()
}

// schotter is called when all arguments are properly set.
func schotter(cmd *cobra.Command, args []string) {
	// Read configuration
	var (
		row    = viper.GetInt("rows")
		col    = viper.GetInt("columns")
		size   = viper.GetFloat64("size")
		margin = viper.GetFloat64("margin")
		order  = viper.GetFloat64("orderliness")
	)

	// Open output file
	out, err := os.Create(args[0])
	if err != nil {
		log.Fatal(err)
	}

	defer out.Close()

	// Create SVG canvas
	canvas := svg.New(out)
	canvas.Start(
		int(2.0*margin+float64(col)*size),
		int(2.0*margin+float64(row)*size),
	)

	defer canvas.End()

	// Generate random Schotter art
	rand.Seed(time.Now().UnixNano())

	for y := 0; y < row; y++ {
		factor := float64(y*y+y) / order

		for x := 0; x < col; x++ {
			// Create randomness
			var (
				offsetX = (rand.Float64()*2 - 1) * factor / 100.0 // Horizontal random offset
				offsetY = (rand.Float64()*2 - 1) * factor / 100.0 // Vertical random offset
				rotate  = (rand.Float64()*2 - 1) * factor / 2.0   // Random rotation
			)

			// Draw rotated and translated square
			canvas.TranslateRotate(
				int(margin+(float64(x)+offsetX)*size),
				int(margin+(float64(y)+offsetY)*size),
				rotate,
			)

			canvas.Rect(0, 0, int(size), int(size),
				"fill:none;stroke:black",
			)

			canvas.Gend()
		}
	}
}
