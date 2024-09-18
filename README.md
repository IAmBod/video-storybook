# Video Thumbnail Sprite Generator

Generate thumbnail sprites from videos using ffmpeg.

## Requirements

- ffmpeg

## Build

```shell
go build -o generator
```

## Usage

```shell
./generator <input> <interval> <maxWidth> <maxHeight> <maxColumns> <output>
```

### Options

| **Name**        | **Description**                                       |
|-----------------|-------------------------------------------------------|
| _\<input>_      | Path to input video file                              |
| _\<interval>_   | Frame interval in seconds                             |
| _\<maxWidth>_   | Maximum width of a single frame in the output sprite  |
| _\<maxHeight>_  | Maximum height of a single frame in the output sprite |
| _\<maxColumns>_ | Maximum number of columns in output sprite            |
| _\<output>_     | Path to output sprite file                            |

# Example

The sprite was generated from [some Blender sample video](https://files.vidstack.io/sprite-fight/720p.mp4).

![](./.assets/output.jpg)