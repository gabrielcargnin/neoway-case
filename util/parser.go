package util

import (
	"bufio"
	"github.com/shopspring/decimal"
	"io"
	"neoway-case/errors"
	"neoway-case/schema"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func Parse(r io.Reader) ([]schema.Consumption, error) {
	const errorMessage errors.Message = "Could not parse file"
	const op errors.Op = "util.Parse"
	s := bufio.NewScanner(r)
	s.Split(bufio.ScanLines)

	consumptions, err := parseRows(s)
	if err != nil {
		return nil, errors.E(op, err, errorMessage)
	}
	return consumptions, err
}

// Iterate over the scanner pointer to scan every line from the file, parsing one by one till EOF.
func parseRows(s *bufio.Scanner) ([]schema.Consumption, error) {
	const errorMessage = "Incorrect record at line "
	const op errors.Op = "util.Parse.parseRows"
	var consumptions []schema.Consumption
	s.Scan()
	count := 2
	for s.Scan() {
		var c schema.Consumption
		rowText := s.Text()
		columns := strings.Fields(rowText)
		if len(columns) != reflect.TypeOf(schema.Consumption{}).NumField() {
			return nil, errors.E(op, errors.Message(errorMessage+strconv.Itoa(count)+". row has more columns than expected"))
		}
		c, err := parseColumns(columns)
		if err != nil {
			return nil, errors.E(op, err, errors.Message(errorMessage+strconv.Itoa(count)))
		}
		consumptions = append(consumptions, c)
		count++
	}
	return consumptions, nil
}

// Parse every column to its type defined in Consumption struct and return it afterwards
func parseColumns(columns []string) (schema.Consumption, error) {
	const errorMessage = "Invalid value for column "
	const op errors.Op = "util.Parse.parseRows.parseColumns"
	const kind errors.Kind = "Parser Error"
	c := schema.Consumption{}
	c.CPF = parseCPF(columns[0])
	private, err := parseInt(columns[1])
	if err != nil {
		return c, errors.E(op, err, errors.Message(errorMessage+"private"), http.StatusBadRequest, kind)
	}
	c.Private = private
	incompleto, err := parseInt(columns[2])
	if err != nil {
		return c, errors.E(op, err, errors.Message(errorMessage+"incompleto"), http.StatusBadRequest, kind)
	}
	c.Incompleto = incompleto
	data, err := parseDataUltimaCompra(columns[3])
	if err != nil {
		return c, errors.E(op, err, errors.Message(errorMessage+"data ultima compra"), http.StatusBadRequest, kind)
	}
	c.DataUltimaCompra = data
	ticketMedio, err := parseTicket(columns[4])
	if err != nil {
		return c, errors.E(op, err, errors.Message(errorMessage+"ticket medio"), http.StatusBadRequest, kind)
	}
	c.TicketMedio = ticketMedio
	ticketUltimaCompra, err := parseTicket(columns[5])
	if err != nil {
		return c, errors.E(op, err, errors.Message(errorMessage+"ticket ultima compra"), http.StatusBadRequest, kind)
	}
	c.TicketUltimaCompra = ticketUltimaCompra
	c.LojaFrequente = parseCNPJ(columns[6])
	c.LojaUltimaCompra = parseCNPJ(columns[7])
	return c, nil
}

func parseInt(field string) (int, error) {
	return strconv.Atoi(field)
}

//Parse a string in the yyyy-MM-dd format, returns error if value can't be converted.
func parseDataUltimaCompra(field string) (*time.Time, error) {
	if nullStringAsNil(field) {
		return nil, nil
	}
	parse, err := time.Parse("2006-01-02", field)
	return &parse, err
}

// Parse a string into time.Time type. If field has commas, will be replaced to dots.
func parseTicket(field string) (*decimal.Decimal, error) {
	if nullStringAsNil(field) {
		return nil, nil
	}
	s := strings.ReplaceAll(field, ",", ".")
	d, err := decimal.NewFromString(s)
	return &d, err
}

func nullStringAsNil(field string) bool {
	re := regexp.MustCompile("NULL")
	return re.MatchString(field)
}

// Receive a string and returns a pointer to one string without dots, slashes and hyphens
func parseCPF(field string) string {
	re := regexp.MustCompile("[.-]")
	return re.ReplaceAllString(field, "")
}

// Receive a string and returns a pointer to one string without dots, slashes and hyphens. If field value equals NULL, a nil pointer is sent back
func parseCNPJ(field string) *string {
	if nullStringAsNil(field) {
		return nil
	}
	re := regexp.MustCompile("[./-]")
	s := re.ReplaceAllString(field, "")
	return &s
}
