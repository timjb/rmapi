package annotations

import (
	"io/ioutil"

	"github.com/jung-kurt/gofpdf"
	"github.com/juruen/rmapi/archive"
	"github.com/juruen/rmapi/encoding/rm"
)

const (
	RmX      = rm.Width
	RmY      = rm.Height
	A4X      = 210
	A4Y      = 297
	ratioA4X = float32(A4X) / float32(RmX)
	ratioA4Y = float32(A4Y) / float32(RmY)
)

type PdfGenerator struct {
	zipName        string
	outputFilePath string
}

func CreatePdfGenerator(zipName, outputFilePath string) PdfGenerator {
	return PdfGenerator{zipName: zipName, outputFilePath: outputFilePath}
}

func (p PdfGenerator) Generate() error {
	reader, err := archive.OpenReader(p.zipName)
	if err != nil {
		return err
	}

	pdf := gofpdf.New("P", "mm", "A4", "")

	for _, page := range reader.Pages {
		if page.Data == nil {
			continue
		}

		f, err := page.Data.Open()
		if err != nil {
			return err
		}
		defer f.Close()

		fb, err := ioutil.ReadAll(f)
		if err != nil {
			return err
		}

		pageRm := rm.New()
		err = pageRm.UnmarshalBinary(fb)
		if err != nil {
			return err
		}

		pdf.AddPage()

		for _, layer := range pageRm.Layers {
			for _, line := range layer.Lines {

				if len(line.Points) < 1 {
					continue
				}

				for i := 1; i < len(line.Points); i++ {
					s := line.Points[i-1]
					x1 := s.X * ratioA4X
					y1 := s.Y * ratioA4Y

					s = line.Points[i]
					x2 := s.X * ratioA4X
					y2 := s.Y * ratioA4Y

					pdf.Line(float64(x1), float64(y1), float64(x2), float64(y2))
				}
			}
		}
	}

	return pdf.OutputFileAndClose(p.outputFilePath)
}
