package video_thumbnail_sprite_generator

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"math"
	"strconv"
)

type VideoMetadata struct {
	Width    int
	Height   int
	Duration float64
}

type StoryboardMetadata struct {
	Url        string                   `json:"url"`
	TileWidth  int                      `json:"tile_width"`
	TileHeight int                      `json:"tile_height"`
	Duration   float64                  `json:"duration"`
	Tiles      []StoryboardMetadataTile `json:"tiles"`
}

type StoryboardMetadataTile struct {
	Start float64 `json:"start"`
	X     int     `json:"x"`
	Y     int     `json:"y"`
}

func GetMetadata(fileName string) (VideoMetadata, error) {
	probe, err := ffmpeg.Probe(fileName)

	if err != nil {
		return VideoMetadata{}, fmt.Errorf("error fetching metadata: %s", err.Error())
	}

	var metadata map[string]interface{}
	err = json.Unmarshal([]byte(probe), &metadata)

	if err != nil {
		return VideoMetadata{}, fmt.Errorf("error parsing metadata: %s", err.Error())
	}

	duration, err := strconv.ParseFloat(metadata["format"].(map[string]interface{})["duration"].(string), 64)

	if err != nil {
		return VideoMetadata{}, fmt.Errorf("error parsing duration: %s", err.Error())
	}

	streams := metadata["streams"].([]interface{})

	for _, stream := range streams {
		streamInfo := stream.(map[string]interface{})

		if streamInfo["codec_type"] == "video" {
			return VideoMetadata{
				Width:    int(streamInfo["width"].(float64)),
				Height:   int(streamInfo["height"].(float64)),
				Duration: duration,
			}, nil
		}
	}

	return VideoMetadata{}, errors.New("could not find video stream")
}

func CalculateTileDimensions(width int, height int, maxWidth int, maxHeight int) (int, int, error) {
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

	return 0, 0, errors.New("could not find divisor within maxWidth and maxHeight")
}

func GenerateStoryboardImage(fileName string, interval int, columns int, rows int, tileWidth int, tileHeight int) (*bytes.Buffer, error) {
	buffer := bytes.NewBuffer(nil)

	err := ffmpeg.
		Input(fileName).
		Filter("fps", ffmpeg.Args{
			fmt.Sprintf("1/%d", interval),
		}).
		Filter("scale", ffmpeg.Args{strconv.Itoa(tileWidth), strconv.Itoa(tileHeight)}).
		Filter("tile", ffmpeg.Args{}, ffmpeg.KwArgs{
			"layout": fmt.Sprintf("%dx%d", columns, rows),
		}).
		Output("pipe:", ffmpeg.KwArgs{
			"format":   "image2",
			"qscale:v": 2,
			"vframes":  1,
		}).
		WithOutput(buffer).
		Silent(false).
		Run()

	return buffer, err
}

func GenerateStoryboardMetadata(url string, interval int, columns int, tileWidth int, tileHeight int, duration float64) StoryboardMetadata {
	tiles := make([]StoryboardMetadataTile, 0)

	tileCount := int(math.Ceil(duration / float64(interval)))

	for i := 0; i < tileCount; i++ {
		start := interval * i

		column := i % columns
		row := i / columns

		x := column * tileWidth
		y := row * tileHeight

		tiles = append(tiles, StoryboardMetadataTile{
			Start: float64(start),
			X:     x,
			Y:     y,
		})
	}

	return StoryboardMetadata{
		Url:        url,
		TileWidth:  tileWidth,
		TileHeight: tileHeight,
		Duration:   duration,
		Tiles:      tiles,
	}
}
