package goobar

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"image"
	"image/png"
	"io"
	"log"
	"mime"
	"net/http"
	"path"
)

// Responder is the type returned by all goobar exchange handlers.
type Responder interface {
	// Respond writes an HTTP response to the specified io.Writer.
	Respond(writer io.Writer) error
	// ContentType returns the MIME type that describes the format
	// of the data returned by Respond.
	ContentType() string
}

// JSON returns a Responder that sends an object in JSON format to
// the client.
func JSON(v interface{}) Responder {
	return &jsonResponder{v}
}

type jsonResponder struct {
	value interface{}
}

func (r jsonResponder) Respond(writer io.Writer) error {
	return json.NewEncoder(writer).Encode(r.value)
}

func (r jsonResponder) ContentType() string {
	return "application/json; charset=utf-8"
}

// PlainText returns a Responder that writes a plain UTF-8 encoded
// text response to the client.
func PlainText(text string) Responder {
	return plainTextResponder(text)
}

type plainTextResponder string

func (r plainTextResponder) Respond(writer io.Writer) error {
	_, err := io.WriteString(writer, string(r))
	return err
}

func (r plainTextResponder) ContentType() string {
	return "text/plain; charset=utf-8"
}

// XML returns a Responder that sends an object in XML format to
// the client. See the xml package for more info on the encoding.
func XML(v interface{}) Responder {
	return &xmlResponder{v}
}

type xmlResponder struct {
	value interface{}
}

func (r xmlResponder) Respond(writer io.Writer) error {
	return xml.NewEncoder(writer).Encode(r.value)
}

func (r xmlResponder) ContentType() string {
	return "text/xml; charset=utf-8"
}

// Nop returns a responder that does not write any data back to the client.
// The handler that returns Nop() is responsible for writing to the client
// using Exchange.Writer().
func Nop() Responder {
	return &nop
}

var nop = nopResponder{}

type nopResponder struct{}

func (r nopResponder) Respond(writer io.Writer) error {
	return nil
}

func (r nopResponder) ContentType() string {
	return ""
}

// ImagePNG returns a Responder that writes a PNG encoded image to the client.
func ImagePNG(img image.Image) Responder {
	return &pngResponder{img}
}

type pngResponder struct {
	img image.Image
}

func (r pngResponder) Respond(writer io.Writer) error {
	return png.Encode(writer, r.img)
}

func (r pngResponder) ContentType() string {
	return "image/png"
}

// Binary returns a Responder that writes binary data to the client.
func Binary(reader io.Reader) Responder {
	return &binaryResponder{reader}
}

type binaryResponder struct {
	reader io.Reader
}

func (r binaryResponder) Respond(writer io.Writer) error {
	_, err := io.Copy(writer, r.reader)
	return err
}

func (r binaryResponder) ContentType() string {
	return "application/octet-stream"
}

// Error returns a Responder that writes an HTTP error to the client.
func Error(statusCode int, message string) Responder {
	return &errorResponder{statusCode, message}
}

type errorResponder struct {
	statusCode int
	message    string
}

func (r errorResponder) Respond(writer io.Writer) error {
	if w, ok := writer.(http.ResponseWriter); ok {
		http.Error(w, r.message, r.statusCode)
		return nil
	}
	_, err := fmt.Fprintf(writer, "Error %d: %s", r.statusCode, r.message)
	return err
}

func (r errorResponder) ContentType() string {
	return "text/plain; charset=utf-8"
}

// View returns a Responder that renders a view from a Go template.
// The template is identified by the specified file path and the given
// model is pipelined to the template renderer.
func View(path string, model interface{}) Responder {
	return &viewResponder{path, model}
}

type viewResponder struct {
	path  string
	model interface{}
}

func (r viewResponder) Respond(writer io.Writer) error {
	templ, name, err := getTemplate(r.path)
	if err != nil {
		log.Print(err)
		return err
	}
	return templ.ExecuteTemplate(writer, name, r.model)
}

func (r viewResponder) ContentType() string {
	return mime.TypeByExtension(path.Ext(r.path))
}
