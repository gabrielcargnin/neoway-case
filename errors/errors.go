package errors

import (
	"log"
	"net/http"
	"strings"
)

type Error struct {
	Op      Op
	Err     error
	Message Message
	Code    int
	Kind    Kind
}

type ResponseError struct {
	ShortMessage    string `json:"shortMessage"`
	DetailedMessage string `json:"detailedMessage"`
	Type            string `json:"type"`
}

func (e Error) Error() string {
	panic("Whoops, something went wrong")
}

type Kind string

type Op string

type Message string

func E(args ...interface{}) error {
	e := &Error{}
	for _, arg := range args {
		switch arg := arg.(type) {
		case Op:
			e.Op = arg
		case error:
			e.Err = arg
		case Message:
			e.Message = arg
		case int:
			e.Code = arg
		case Kind:
			e.Kind = arg
		default:
			panic("Incorrect call to E")
		}

	}
	return e
}

func LogCleanStackTrace(err error) {
	log.Printf("Error in ")
	ops := ops(err.(*Error))
	for i := len(ops) - 1; i >= 0; i-- {
		log.Println(ops[i])
	}
}

func ops(err *Error) []Op {
	res := []Op{err.Op}
	subErr, ok := err.Err.(*Error)
	if !ok {
		return res
	}
	res = append(res, ops(subErr)...)
	return res
}

func message(err *Error) []Message {
	res := []Message{err.Message}
	subErr, ok := err.Err.(*Error)
	if !ok {
		return res
	}
	res = append(res, message(subErr)...)
	return res
}

func shortMessage(err *Error) string {
	messages := message(err)
	i := len(messages)
	return string(messages[i-1])
}

func detailedMessage(err *Error) string {
	messages := message(err)
	var res []string
	for _, m := range messages {
		res = append(res, string(m))
	}
	return strings.Join(res, ". ")
}

func Code(err error) int {
	e, ok := err.(*Error)
	if !ok {
		return http.StatusInternalServerError
	}
	if e.Code != 0 {
		return e.Code
	}
	return Code(e.Err)
}

func kind(err error) string {
	e, ok := err.(*Error)
	if !ok {
		return "Unknown"
	}
	if e.Kind != "" {
		return string(e.Kind)
	}
	return kind(e.Err)
}

func GetResponseErr(err error) ResponseError {
	return ResponseError{
		DetailedMessage: detailedMessage(err.(*Error)),
		ShortMessage:    shortMessage(err.(*Error)),
		Type:            kind(err),
	}
}
