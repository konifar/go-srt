# go-srt
[![](https://github.com/konifar/go-srt/workflows/CI/badge.svg)](https://github.com/konifar/go-srt/actions)

Parser for the simple SubRip (.srt) file.

## Description

go-srt just supports simple srt file format like below.

```
1
00:00:00,000 --> 00:00:00,000
Don-don donuts! Let's go nuts!
```

it does not support html decoration, subtitle rectangle and orientation. 

## Installation

```shell
$ go get github.com/konifar/go-srt
```

## Usage

### Parse from file name

```go
subtitles, err := gosrt.ReadFile("your_sub_rip_file.srt")
```

### Parse from io.Reader

```go
var f *os.File
os.Open(filepath.Clean(fileName))
if err != nil {
    err = errors.Wrapf(err, "failed to open file:%s", fileName)
    return []Subtitle{}, err
}
defer func() {
    err = f.Close()
    if err != nil {
        log.Fatalln(err)
    }
}()

gosrt.ReadSubtitles(f)
```