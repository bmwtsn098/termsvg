package export

import (
	"bytes"
	"os"

	"github.com/mrmarble/termsvg/internal/svg"
	"github.com/mrmarble/termsvg/pkg/asciicast"
	"github.com/rs/zerolog/log"
	"github.com/tdewolff/minify/v2"
	msvg "github.com/tdewolff/minify/v2/svg"
)

type Cmd struct {
	Input  string `short:"i" type:"existingfile" help:"asciicast file to export"`
	Output string `optional:"" short:"o" type:"path" help:"where to save the file. Defaults to <input>.svg"`
	Mini   bool   `name:"minify" optional:"" short:"m" help:"minify output file. May be slower"`
}

func (cmd *Cmd) Run() error {
	output := cmd.Output
	if output == "" {
		output = cmd.Input + ".svg"
	}

	err := export(cmd.Input, output, cmd.Mini)
	if err != nil {
		return err
	}

	log.Info().Str("output", output).Msg("svg file saved.")

	return nil
}

func export(input, output string, mini bool) error {
	inputFile, err := os.ReadFile(input)
	if err != nil {
		return err
	}

	cast, err := asciicast.Unmarshal(inputFile)
	if err != nil {
		return err
	}

	outputFile, err := os.Create(output)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	if mini {
		out := new(bytes.Buffer)
		svg.Export(*cast, out)

		m := minify.New()
		m.AddFunc("image/svg+xml", msvg.Minify)

		b, err := m.Bytes("image/svg+xml", out.Bytes())
		if err != nil {
			return err
		}

		_, err = outputFile.Write(b)
		if err != nil {
			return err
		}
	} else {
		svg.Export(*cast, outputFile)
	}

	return nil
}
