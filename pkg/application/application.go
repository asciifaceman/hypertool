package application

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/k3a/html2text"
)

var appVersion = "0.1.0"

type Manager struct {
	Application  fyne.App
	WStatus      binding.String
	supertoolURL string
	servedHTML   string
	log          binding.String
	title        string
}

func (m *Manager) Init() error {
	m.title = fmt.Sprintf("HyperTool v%s", appVersion)
	m.Application = app.New()
	m.WStatus = binding.NewString()
	m.log = binding.NewString()
	m.WStatus.Set("[not started]")
	m.Splash()
	m.Run()
	return nil
}

func (m *Manager) Splash() error {
	w := m.Application.NewWindow(m.title)
	welcome := widget.NewLabel("Welcome to HyperTool.\n\nEnter the supertool URL you wish to view below\nand then click Go!")
	urlEntry := widget.NewEntry()
	urlEntry.SetPlaceHolder("Enter Supertool URL here")

	w.SetContent(container.NewVBox(
		welcome,
		urlEntry,
		widget.NewButton("Go!", func() {
			m.Go(urlEntry.Text)
			w.Close()
		}),
		quitButton("Exit"),
	))

	w.Canvas().Focus(urlEntry)

	w.Show()

	return nil
}

func (m *Manager) Go(url string) {
	m.supertoolURL = url
	m.ServerStatus()
	m.ProcessSupertool()
}

func (m *Manager) ServerStatus() {
	w := m.Application.NewWindow(m.title)
	urlLabel := widget.NewLabel(fmt.Sprintf("Source URL: %s", m.supertoolURL))
	statusLabel := widget.NewLabel("Server status: ")
	status := widget.NewLabel("")
	logLabel := widget.NewLabel("Server logs: ")
	logs := widget.NewMultiLineEntry()
	logs.Bind(m.log)
	status.Bind(m.WStatus)

	w.SetContent(container.NewVBox(
		urlLabel,
		statusLabel,
		status,
		logLabel,
		logs,
		widget.NewButton("Open Browser", func() {
			m.open("http://localhost:8080")
		}),
		quitButton("Exit"),
	))

	w.Show()

	if m.supertoolURL == "" {
		err := errors.New("Given URL is empty")
		dialog.ShowError(err, w)
	}

	m.StartServer()
}

func (m *Manager) ModalError(err error) {
	w := m.Application.NewWindow("Error")
	w.SetContent(container.NewVBox(
		widget.NewLabel(err.Error()),
		widget.NewButton("Oh no, but continue.", func() {
			w.Close()
		}),
		widget.NewButton("Oh no, Exit.", func() {
			os.Exit(1)
		}),
	))
	w.Show()

}

func (m *Manager) ModalMessage(title string, message string) {
	w := m.Application.NewWindow(title)
	w.SetContent(container.NewVBox(
		widget.NewLabel(message),
		widget.NewButton("Ok.", func() {
			w.Close()
		}),
	))
	w.Show()
}

func (m *Manager) Run() {
	m.Application.Run()
}

func (m *Manager) ProcessSupertool() {
	m.Info(fmt.Sprintf("Processing %s", m.supertoolURL))
	resp, err := http.Get(m.supertoolURL)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	defer resp.Body.Close()

	htmlr, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	/*
	   Modify the HTML here
	*/

	plain := html2text.HTML2Text(string(htmlr))
	plain = fmt.Sprintf("%s%s%s", openingTemplate, plain, closingTemplate)
	plain = strings.Replace(plain, "\r\n", "<br>", -1)

	plain = reFindImage.ReplaceAllString(plain, "<br><img src=\"${1}\" width=\"250\"><br>")

	plain = reFindProductId.ReplaceAllString(plain, "<br><hr><br><h2>${1}</h2>")

	plain = reFindCurrency.ReplaceAllString(plain, "<br><pre>                       ${1}</pre><br>")
	m.servedHTML = plain
	m.Info("Done, ready to serve...")
}

func (m *Manager) StartServer() {
	m.Info("Starting server...")
	http.HandleFunc("/", m.ServeTools)
	var handler http.Handler = http.DefaultServeMux
	handler = m.logRequestHandler(handler)

	go func() {
		http.ListenAndServe(":8080", handler)
	}()
	m.WStatus.Set("[started]")
	m.Info("Server started...")
}

func (m *Manager) ServeTools(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, m.servedHTML)
}

func (m *Manager) logRequestHandler(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)

		uri := r.URL.String()
		method := r.Method

		m.Info(fmt.Sprintf("[%s] %s", method, uri))
	}

	return http.HandlerFunc(fn)
}

func (m *Manager) Info(message string) {
	currentLog, err := m.log.Get()
	if err != nil {
		m.ModalError(err)
	}
	m.log.Set(fmt.Sprintf("%s\n%s", message, currentLog))
}

// open opens the specified URL in the default browser of the user.
func (m *Manager) open(url string) error {
	m.Info(fmt.Sprintf("Opening %s in the browser", url))
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	err := exec.Command(cmd, args...).Start()
	if err != nil {
		m.ModalError(err)
	}
	return nil
}
