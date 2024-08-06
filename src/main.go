package main

import (
	"fmt"
	"net"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func scanPort(protocol, hostname string, port int) bool {
	address := fmt.Sprintf("%s:%d", hostname, port)
	conn, err := net.DialTimeout(protocol, address, 1*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func main() {
	a := app.New()
	w := a.NewWindow("Port Scanner")
	stop_scan := make(chan bool)

	ipEntry := widget.NewEntry()
	ipEntry.SetPlaceHolder("Enter IP Address")

	portEntry := widget.NewEntry()
	portEntry.SetPlaceHolder("Enter Port Range (e.g., 1-1024)")

	result := widget.NewMultiLineEntry()

	result.SetPlaceHolder("Scan results will appear here...")

	result.Wrapping = fyne.TextWrapWord
	result.Resize(fyne.NewSize(400, 200))

	actualScanLabel := widget.NewLabel("")

	w.SetFixedSize(true)
	scanButton := widget.NewButton("Scan", nil)

	stopScanButton := widget.NewButton("Stop Scan", func() {
		scanButton.Enable()
		scanButton.SetText("Scan")
		actualScanLabel.SetText("Scan Stopped")

		stop_scan <- true

	})

	stopScanButton.Disable()

	scanButton = widget.NewButton("Scan", func() {
		go onBtnClick(ipEntry, portEntry, result, scanButton, actualScanLabel, stop_scan, stopScanButton)
	})

	content := container.NewVBox(
		widget.NewLabel("Simple Port Scanner"),
		ipEntry,
		portEntry,
		scanButton,
		result,
		actualScanLabel,
		stopScanButton,
	)

	w.SetContent(content)
	w.Resize(fyne.NewSize(500, 500))
	w.ShowAndRun()
}

func onBtnClick(ipEntry, portEntry *widget.Entry, result *widget.Entry, scanButton *widget.Button, actualScanLabel *widget.Label, stop_scan chan bool, stopScanButton *widget.Button) {
	ip := ipEntry.Text
	portRange := portEntry.Text
	result.SetText("")

	if ip == "" || portRange == "" {
		stopScanButton.Disable()
		actualScanLabel.SetText("Please enter both IP address and port range")
		actualScanLabel.TextStyle = fyne.TextStyle{Bold: true, Italic: true}
		return
	}

	var startPort, endPort int
	fmt.Sscanf(portRange, "%d-%d", &startPort, &endPort)

	stopScanButton.Enable()
	for port := startPort; port <= endPort; port++ {
		select {
		case <-stop_scan:
			// stop_scan <- false
			stopScanButton.Disable()
			return

		default:

			scanButton.Disable()
			scanButton.SetText("Scanning...")
			actualScanLabel.SetText(fmt.Sprintf("%s - Scanning port %d", ip, port))
			if scanPort("tcp", ip, port) {
				result.SetText(result.Text + fmt.Sprintf("Port %d is open\n", port))
				// } else {
				// 	// result.SetText(result.Text + fmt.Sprintf("Port %d is closed\n", port))
				// 	println(fmt.Sprintf("Port %d is closed\n", port))
				// }

				scanButton.SetText("Scan")
				scanButton.Enable()

			}

			scanButton.SetText("Scan")
			scanButton.Enable()

		}
	}

	scanButton.SetText("Scan")

	scanButton.Enable()

	stopScanButton.Disable()

}
