package dataurl

import (
	"fmt"
	"io"
	"log"
	"mime"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/vincent-petithory/dataurl"

	"github.com/sagan/stool/cmd"
)

const MAX_SIZE = 2 * 1024 * 1024 // 2MiB

// command represents the base command when called without any subcommands
var command = &cobra.Command{
	Use:   "dataurl {file}",
	Short: "Generate Data URL for file",
	Long: `Generate Data URL for file
{file} is input filename, set to "-" to use stdin.

It outputs to stdout.`,
	RunE: dataurlcmd,
	Args: cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
}

var (
	performDecode bool
	asciiEncoding bool
	mimetype      string
)

func init() {
	command.Flags().BoolVarP(&performDecode, "decode", "", false, "Decode mode instead of default encode mode")
	command.Flags().BoolVar(&asciiEncoding, "ascii", false, "Force ascii encoding instead of base64")
	command.Flags().StringVarP(&mimetype, "mimetype", "", "", "Force mime type")
	cmd.RootCmd.AddCommand(command)
}

func dataurlcmd(cmd *cobra.Command, args []string) (err error) {
	encoding := dataurl.EncodingBase64
	if asciiEncoding {
		encoding = dataurl.EncodingASCII
	}
	filename := args[0]
	var in *os.File
	var detectedMimetype string
	if filename == "-" {
		in = os.Stdin
	} else {
		if stat, err := os.Stat(filename); err != nil {
			log.Fatalf("failed to stat file: %v", err)
		} else if stat.IsDir() {
			log.Fatalf("not a regular file")
		} else if stat.Size() >= MAX_SIZE {
			log.Printf("Warning: file is too large (%d)", stat.Size())
		}
		detectedMimetype = mime.TypeByExtension(filepath.Ext(filename))
		in, err = os.Open(filename)
		if err != nil {
			log.Fatalf("failed to open file: %v", err)
		}
		defer in.Close()
	}

	switch {
	case mimetype == "" && detectedMimetype == "":
		mimetype = "application/octet-stream"
	case mimetype == "" && detectedMimetype != "":
		mimetype = detectedMimetype
	}
	out := os.Stdout

	if performDecode {
		if err := decode(in, out); err != nil {
			log.Fatal(err)
		}
	} else {
		if err := encode(in, out, encoding, mimetype); err != nil {
			log.Fatal(err)
		}
	}
	return nil
}

func decode(in io.Reader, out io.Writer) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()

	du, err := dataurl.Decode(in)
	if err != nil {
		return
	}

	_, err = out.Write(du.Data)
	if err != nil {
		return
	}
	return
}

func encode(in io.Reader, out io.Writer, encoding string, mediatype string) (err error) {
	defer func() {
		if e := recover(); e != nil {
			var ok bool
			err, ok = e.(error)
			if !ok {
				err = fmt.Errorf("%v", e)
			}
			return
		}
	}()
	b, err := io.ReadAll(in)
	if err != nil {
		return
	}

	du := dataurl.New(b, mediatype)
	du.Encoding = encoding

	_, err = du.WriteTo(out)
	if err != nil {
		return
	}
	return
}
