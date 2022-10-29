package clipboard

import (
	"context"

	"github.com/ntsd/cross-clipboard/pkg/config"
	"github.com/ntsd/cross-clipboard/pkg/device"
	"golang.design/x/clipboard"
)

// ClipboardManager struct for clipbaord manager
type ClipboardManager struct {
	config *config.Config

	ReadTextChannel   <-chan []byte
	ReadImageChannel  <-chan []byte
	clipboards        []Clipboard
	ClipboardsChannel chan []Clipboard
	receivedClipboard *Clipboard
}

// NewClipboardManager create new clipbaord manager
func NewClipboardManager(cfg *config.Config) *ClipboardManager {
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}

	textCh := clipboard.Watch(context.Background(), clipboard.FmtText)
	imgCh := clipboard.Watch(context.Background(), clipboard.FmtImage)

	return &ClipboardManager{
		config:            cfg,
		ReadTextChannel:   textCh,
		ReadImageChannel:  imgCh,
		ClipboardsChannel: make(chan []Clipboard),
		clipboards:        []Clipboard{},
	}
}

// limitAppend append and rotate when limit
func limitAppend[T any](limit int, slice []T, new T) []T {
	l := len(slice)
	if l >= limit {
		slice = slice[1:]
	}
	slice = append(slice, new)
	return slice
}

// WriteClipboard write os clipbaord
func (c *ClipboardManager) WriteClipboard(newClipboard Clipboard) {
	c.receivedClipboard = &newClipboard

	if newClipboard.IsImage {
		clipboard.Write(clipboard.FmtImage, newClipboard.Data)
		return
	}
	clipboard.Write(clipboard.FmtText, newClipboard.Data)
}

// UpdateClipboard add clipbaord to clipbaord history
func (c *ClipboardManager) UpdateClipboard(newClipboard Clipboard) {
	c.clipboards = limitAppend(c.config.MaxHistory, c.clipboards, newClipboard)
	c.ClipboardsChannel <- c.clipboards
}

func (c *ClipboardManager) IsReceivedClipboardFromDevice(dv *device.Device) bool {
	if c.receivedClipboard == nil {
		return false
	}

	if c.receivedClipboard.Device == nil {
		return false
	}

	return c.receivedClipboard.Device.AddressInfo.ID.Pretty() == dv.AddressInfo.ID.Pretty()
}
