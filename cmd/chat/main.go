/*
 * Go Programming Blueprints - 2nd ed - by Mat Ryer
 *
 * Copyright (c) 2021 Michael D Henderson
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package main

import (
	"flag"
	"fmt"
	"github.com/mdhender/goblueprints/pkg/trace"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"
)

// templateHandler implements a handler for loading, compiling, and
// serving our template.
type templateHandler struct {
	once     sync.Once
	root     string // root of application data
	filename string
	templ    *template.Template // represents a single template
}

// ServeHTTP handles the HTTP request.
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join(t.root, "templates", t.filename)))
	})
	_ = t.templ.Execute(w, r)
}

func main() {
	cfg := struct {
		host    string
		port    string
		root    string
		tracing bool
	}{
		host: "localhost",
		port: "8080",
		root: "D:/GoLand/goblueprints/cmd/chat",
	}
	flag.StringVar(&cfg.host, "host", cfg.host, fmt.Sprintf("Host address of the application (default %q)", cfg.host))
	flag.StringVar(&cfg.port, "port", cfg.port, fmt.Sprintf("Port of the application (default %q)", cfg.port))
	flag.StringVar(&cfg.root, "root", cfg.root, fmt.Sprintf("Location of application data (default %q)", cfg.root))
	flag.BoolVar(&cfg.tracing, "tracing", cfg.tracing, fmt.Sprintf("Turn tracing on (default %v)", cfg.tracing))
	flag.Parse() // parse the flags

	r := newRoom()
	if cfg.tracing {
		r.tracer = trace.New(os.Stdout)
	}

	http.Handle("/chat", MustAuth(&templateHandler{root: cfg.root, filename: "chat.html"}))
	http.Handle("/room", r)

	// get the room going
	go r.run()

	// start the web server
	log.Println("Starting web server on", net.JoinHostPort(cfg.host, cfg.port))
	if err := http.ListenAndServe(net.JoinHostPort(cfg.host, cfg.port), nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
