package goobar

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"strconv"
)

// Exchange contains information about a HTTP exchange (request and response),
// providing comfort functions to get typed form or query values.
type Exchange struct {
	w  http.ResponseWriter
	r  *http.Request
	id string
}

// Request returns a pointer to the http.Request contained in the Exchange.
func (x *Exchange) Request() *http.Request {
	return x.r
}

// Writer returns the http.ResponseWriter contained in the Exchange.
func (x *Exchange) Writer() http.ResponseWriter {
	return x.w
}

// GetID returns the numeric ID passed with the request, which is the last element
// of the requested URL path (e.g. for http://host:port/item/action/1, GetID
// would return 1).
// If no numeric ID is present, false is returned for the second return value.
func (x *Exchange) GetID() (int, bool) {
	log.Printf("ID: %s -> int", x.id)
	id, err := strconv.Atoi(x.id)
	return id, err == nil
}

// MustGetID returns the numeric ID passed with the request (see GetID) or panics
// if no numeric ID is present.
// The goobar HTTP handler will recover from this panic and write a 404 error
// to the client.
func (x *Exchange) MustGetID() int {
	log.Printf("ID: %s -> int", x.id)
	id, err := strconv.Atoi(x.id)
	if err != nil {
		x.doPanic("id not present or in the wrong format")
	}
	return id
}

// GetIDString returns the ID passed with the request (see GetID) or an
// empty string if no ID is present.
func (x *Exchange) GetIDString() string {
	log.Printf("ID: %s -> string", x.id)
	return x.id
}

// MustGetIDString returns the ID passed with the request (see GetID) or panics
// if no ID is present.
// The goobar HTTP handler will recover from this panic and write a 404 error
// to the client.
func (x *Exchange) MustGetIDString() string {
	log.Printf("ID: %s -> string", x.id)
	if x.id == "" {
		x.doPanic("id not present")
	}
	return x.id
}

// GetInt returns the integer value passed by query or form with the given
// key.
// If the given key is not present or in the wrong format, false is returned
// for the second return value.
func (x *Exchange) GetInt(key string) (int, bool) {
	s := x.r.FormValue(key)
	log.Printf("%s: %s -> int", key, s)
	if s == "" {
		return 0, false
	}
	val, err := strconv.Atoi(s)
	return val, err == nil
}

// MustGetInt returns the integer value passed by query or form with the given
// key or panics if not present or the found value is in the wrong format.
func (x *Exchange) MustGetInt(key string) int {
	s := x.r.FormValue(key)
	log.Printf("%s: %s -> int", key, s)
	val, err := strconv.Atoi(s)
	if err != nil {
		x.doPanic(fmt.Sprintf("value %s not present or not an integer", key))
	}
	return val
}

// GetIntOrDefault returns the integer value passed by query or form with the
// given key or the specified default value if no such value is present.
func (x *Exchange) GetIntOrDefault(key string, defaultVal int) int {
	if val, ok := x.GetInt(key); ok {
		return val
	}
	return defaultVal
}

// GetString returns the string value passed by query or form with the
// given key or the empty string if no value with the given key is present.
func (x *Exchange) GetString(key string) string {
	s := x.r.FormValue(key)
	log.Printf("%s: %s -> string", key, s)
	return s
}

// MustGetString returns the string value passed by query or form with the
// given key or panics if no value with the given key is present.
func (x *Exchange) MustGetString(key string) string {
	s := x.r.FormValue(key)
	log.Printf("%s: %s -> string", key, s)
	if s == "" {
		x.doPanic(fmt.Sprintf("value %s not present", key))
	}
	return s
}

// GetStringOrDefault returns the string value passed by query or form with the
// given key or the specified default if no value with the given key is present.
func (x *Exchange) GetStringOrDefault(key string, defaultVal string) string {
	if val := x.GetString(key); val != "" {
		return val
	}
	return defaultVal
}

// GetFloat returns the floating-point value passed by query or form with the given
// key.
// If the given key is not present or in the wrong format, false is returned
// for the second return value.
func (x *Exchange) GetFloat(key string) (float64, bool) {
	s := x.r.FormValue(key)
	log.Printf("%s: %s -> float", key, s)
	val, err := strconv.ParseFloat(s, 64)
	return val, err == nil
}

// MustGetFloat returns the floating-point value passed by query or form with the given
// key or panics if not present or the found value is in the wrong format.
func (x *Exchange) MustGetFloat(key string) float64 {
	s := x.r.FormValue(key)
	log.Printf("%s: %s -> float", key, s)
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		x.doPanic(fmt.Sprintf("value %s not present or not a float", key))
	}
	return val
}

// GetFloatOrDefault returns the floating-point value passed by query or form with the
// given key or the specified default value if no such value is present.
func (x *Exchange) GetFloatOrDefault(key string, defaultVal float64) float64 {
	if val, ok := x.GetFloat(key); ok {
		return val
	}
	return defaultVal
}

// GetBool returns the boolean value passed by query or form with the given
// key.
// If the given key is not present or in the wrong format, false is returned
// for the second return value.
func (x *Exchange) GetBool(key string) (bool, bool) {
	s := x.r.FormValue(key)
	log.Printf("%s: %s -> bool", key, s)
	val, err := strconv.ParseBool(s)
	return val, err == nil
}

// MustGetBool returns the boolean value passed by query or form with the given
// key or panics if not present or the found value is in the wrong format.
func (x *Exchange) MustGetBool(key string) bool {
	s := x.r.FormValue(key)
	log.Printf("%s: %s -> bool", key, s)
	val, err := strconv.ParseBool(s)
	if err != nil {
		x.doPanic(fmt.Sprintf("value %s not present or not a boolean", key))
	}
	return val
}

// GetBoolOrDefault returns the boolean value passed by query or form with the
// given key or the specified default value if no such value is present.
func (x *Exchange) GetBoolOrDefault(key string, defaultVal bool) bool {
	if val, ok := x.GetBool(key); ok {
		return val
	}
	return defaultVal
}

///////////////////////////////////////////////////////////////////////////////

type actionPanic struct {
	msg string
}

func (x *Exchange) doPanic(msg string) {
	panic(actionPanic{msg})
}

func makeExchange(w http.ResponseWriter, r *http.Request) *Exchange {
	_, id := path.Split(r.URL.Path)

	return &Exchange{
		w:  w,
		r:  r,
		id: id,
	}
}
