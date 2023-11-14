package main

import (
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/emicklei/nanny"
)

type Bike struct {
	Brand, Model, Year string
}

func main() {
	// record max 100 events
	rec := nanny.NewRecorder(nanny.WithMaxEvents(100))

	// fallback logger (cannot be the default handler)
	txt := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})

	// handler capturing debug and up
	slog.SetDefault(slog.New(nanny.NewLogHandler(rec, txt, slog.LevelDebug)))

	// serve do on /do
	http.HandleFunc("/do", do)

	// serve captured events
	http.Handle("/nanny", nanny.NewBrowser(rec))

	// serve
	slog.Info("open http://localhost:8080/do")
	slog.Info("then open http://localhost:8080/nanny")
	http.ListenAndServe(":8080", http.DefaultServeMux)
}

func do(w http.ResponseWriter, r *http.Request) {
	l := slog.Default().WithGroup("group")

	l.Debug("checking...", slog.Any("bike", Bike{Brand: "Trek", Model: "Emonda", Year: "2017"}))

	// wont see this event in the recorder
	l.Info("no attributes")

	l.Info("one attribute", slog.String("bike", "Trek"))
	io.WriteString(w, "done")
}
