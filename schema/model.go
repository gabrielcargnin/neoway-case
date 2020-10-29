package schema

import (
	"github.com/shopspring/decimal"
	"time"
)

// Consumption represents a row coming from the file.
type Consumption struct {
	CPF                string `validate:"cpf,required"`
	Private            int    `validate:"eq=0|eq=1"`
	Incompleto         int    `validate:"eq=0|eq=1"`
	DataUltimaCompra   *time.Time
	TicketMedio        *decimal.Decimal
	TicketUltimaCompra *decimal.Decimal
	LojaFrequente      *string `validate:"cnpj"`
	LojaUltimaCompra   *string `validate:"cnpj"`
}
