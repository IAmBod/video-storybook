package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"math"
	"os"
	"strconv"
)

type Metadata struct {
	Width    int
	Height   int
	Duration int
}

func main() {
	fileName := os.Args[1]
	interval, err := strconv.Atoi(os.Args[2])

	if err != nil {
		panic(err)
	}

	maxWidth, err := strconv.Atoi(os.Args[3])

	if err != nil {
		panic(err)
	}

	maxHeight, err := strconv.Atoi(os.Args[4])

	if err != nil {
		panic(err)
	}

	maxColumns, err := strconv.Atoi(os.Args[5])

	if err != nil {
		panic(err)
	}

	outputFileName := os.Args[6]

	metadata, err := GetMetadata(fileName)

	if err != nil {
		panic(err)
	}

	frameCount := metadata.Duration / interval
	frameWidth, frameHeight, err := Calculate(metadata.Width, metadata.Height, maxWidth, maxHeight)

	if err != nil {
		panic(err)
	}

	spriteBuffer, err := CreateSprite(fileName, interval, frameCount, frameWidth, frameHeight, maxColumns)

	if err != nil {
		panic(err)
	}

	err = os.WriteFile(outputFileName, spriteBuffer.Bytes(), 0777)

	if err != nil {
		panic("Error writing file: " + err.Error())
	}
}

func GetMetadata(fileName string) (Metadata, error) {
	probe, err := ffmpeg.Probe(fileName)

	if err != nil {
		return Metadata{}, err
	}

	var metadata map[string]interface{}
	err = json.Unmarshal([]byte(probe), &metadata)

	if err != nil {
		return Metadata{}, err
	}

	duration, err := strconv.ParseFloat(metadata["format"].(map[string]interface{})["duration"].(string), 64)

	if err != nil {
		return Metadata{}, err
	}

	streams := metadata["streams"].([]interface{})

	for _, stream := range streams {
		streamInfo := stream.(map[string]interface{})

		if streamInfo["codec_type"] == "video" {
			return Metadata{
				Width:    int(streamInfo["width"].(float64)),
				Height:   int(streamInfo["height"].(float64)),
				Duration: int(duration),
			}, nil
		}
	}

	return Metadata{}, errors.New("could not find video stream")
}

func Calculate(width int, height int, maxWidth int, maxHeight int) (int, int, error) {
	highestDivisor := min(width, height) / 2

	for i := 2; i < highestDivisor; i++ {
		if width%i != 0 {
			continue
		}

		if height%i != 0 {
			continue
		}

		if width/i > maxWidth {
			continue
		}

		if height/i > maxHeight {
			continue
		}

		return width / i, height / i, nil
	}

	return 0, 0, errors.New("could not find divisor")
}

func CreateSprite(fileName string, interval int, frameCount int, frameWidth int, frameHeight int, maxColumn int) (*bytes.Buffer, error) {
	gridColumns := min(frameCount, maxColumn)
	gridRows := int(math.Ceil(float64(frameCount) / float64(maxColumn)))

	img := image.NewRGBA(image.Rect(0, 0, gridColumns*frameWidth, gridRows*frameHeight))
	draw.Draw(img, img.Bounds(), &image.Uniform{C: color.Black}, image.Pt(0, 0), draw.Src)

	for i := 0; i < frameCount; i++ {
		buffer, err := ReadFrame(fileName, i*interval, frameWidth, frameHeight)

		if err != nil {
			return nil, err
		}

		reader := bytes.NewReader(buffer.Bytes())
		frame, err := jpeg.Decode(reader)

		if err != nil {
			return nil, err
		}

		column := i % maxColumn
		row := i / maxColumn

		draw.Draw(img, img.Bounds().Add(image.Pt(column*frameWidth, row*frameHeight)), frame, image.Pt(0, 0), draw.Over)
	}

	outputBuffer := bytes.NewBuffer(nil)
	err := jpeg.Encode(outputBuffer, img, nil)

	if err != nil {
		return nil, err
	}

	return outputBuffer, nil
}

func ReadFrame(fileName string, seconds int, frameWidth int, frameHeight int) (*bytes.Buffer, error) {
	buffer := bytes.NewBuffer(nil)

	err := ffmpeg.
		Input(fileName, ffmpeg.KwArgs{"ss": seconds}).
		Output("pipe:", ffmpeg.KwArgs{
			"format":  "image2",
			"s":       fmt.Sprintf("%dx%d", frameWidth, frameHeight),
			"vcodec":  "mjpeg",
			"vframes": 1,
		}).
		WithOutput(buffer).
		Silent(false).
		Run()

	return buffer, err
}
