package vgui

import (
    "bufio"
    "os"
    "strings"
)

type FileSourceProvider interface {
    ReadAllText(file string) (string, error)
    ReadAllLines(file string) ([]string, error)
}

type HudFileSourceProvider struct {
    root string
}

func (p *HudFileSourceProvider) ReadAllText(path string) (string, error) {
    lines, err := p.ReadAllLines(path)
    if err != nil {
        return "", err
    }

    return strings.Join(lines, "\n"), nil
}

func (p *HudFileSourceProvider) ReadAllLines(path string) ([]string, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close() // TODO: file close failure?

    scanner := bufio.NewScanner(file)
    scanner.Split(bufio.ScanLines)
    var text []string

    for scanner.Scan() {
        text = append(text, scanner.Text())
    }

    return text, nil
}