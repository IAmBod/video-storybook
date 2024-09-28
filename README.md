# Video Thumbnail Sprite Generator

Generate thumbnail sprites from videos using ffmpeg.

## Requirements

- ffmpeg

# Example

```go
metadata, err := GetMetadata(fileName)

if err != nil {
    log.Fatalln("Error while reading video metadata: " + err.Error())
}

frameCount := metadata.Duration / interval
frameWidth, frameHeight, err := CalculateFrameDimensions(metadata.Width, metadata.Height, maxWidth, maxHeight)

if err != nil {
    log.Fatalln("Error calculating sprite frame dimensions: " + err.Error())
}

gridColumns := min(frameCount, maxColumns)
gridRows := int(math.Ceil(float64(frameCount) / float64(maxColumns)))
spriteBuffer, err := GenerateSprite(fileName, interval, gridColumns, gridRows, frameWidth, frameHeight)

if err != nil {
    log.Fatalln("Error calculating creating sprite: " + err.Error())
}

err = os.WriteFile(outputFileName, spriteBuffer.Bytes(), 0777)
```

The sprite was generated from [some Blender sample video](https://files.vidstack.io/sprite-fight/720p.mp4).

![](./.assets/output.jpg)