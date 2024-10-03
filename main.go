package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/yqs112358/cross-clipboard/ui"
	"log"
	"os"
	"os/signal"

	"github.com/yqs112358/cross-clipboard/pkg/config"
	"github.com/yqs112358/cross-clipboard/pkg/crossclipboard"
	"github.com/yqs112358/cross-clipboard/pkg/device"
	"github.com/yqs112358/cross-clipboard/pkg/xerror"
)

func main() {
	tuiMode := flag.Bool("tui", false, "use terminal ui")
	flag.Parse()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	crossClipboard, err := crossclipboard.NewCrossClipboard(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if tuiMode != nil && *tuiMode {
		// TUI mode
		view := ui.NewView(crossClipboard)
		view.Start()
	} else {
		// Terminal mode
		exitSignal := make(chan os.Signal, 1)
		signal.Notify(exitSignal, os.Interrupt)

		for {
			select {
			case l := <-crossClipboard.LogChan:
				log.Println("log: ", l)
			case err := <-crossClipboard.ErrorChan:
				var fatalErr *xerror.FatalError
				if errors.As(err, &fatalErr) {
					log.Fatal(fmt.Errorf("fatal error: %w", fatalErr))
				}
				log.Println(fmt.Errorf("runtime error: %w", err))
			case <-crossClipboard.ClipboardManager.ClipboardsHistoryUpdated:
				// log.Printf("clipboard history updated, history size %d", len(crossClipboard.ClipboardManager.ClipboardsHistory))
			case <-crossClipboard.DeviceManager.DevicesUpdated:
				for _, dv := range crossClipboard.DeviceManager.Devices {
					if dv.Status == device.StatusPending {
						fmt.Printf("device %s wanted to connect (Y/n)", dv.Name)
						var input string
						fmt.Scanln(&input)
						if input == "n" {
							dv.Block()
						} else {
							err = dv.Trust()
							if err != nil {
								log.Println(fmt.Errorf("can not trust device: %w", err))
							}
						}
						crossClipboard.DeviceManager.UpdateDevice(dv)
					}
				}
			case exit := <-exitSignal:
				log.Printf("got %s signal. aborting...\n", exit)
				err := crossClipboard.Stop()
				if err != nil {
					log.Panicln(fmt.Errorf("error to graceful eixt: %w", err))
				}
				os.Exit(0)
			}
		}
	}
}
