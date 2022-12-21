// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"entgo.io/ent/dialect/sql"
	"github.com/noamcattan/geni/ent/account"
)

// Account is the model entity for the Account schema.
type Account struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// SheetsCredentials holds the value of the "sheets_credentials" field.
	SheetsCredentials []byte `json:"sheets_credentials,omitempty"`
	// SpreadsheetID holds the value of the "spreadsheet_id" field.
	SpreadsheetID string `json:"spreadsheet_id,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the AccountQuery when eager-loading is set.
	Edges AccountEdges `json:"edges"`
}

// AccountEdges holds the relations/edges for other nodes in the graph.
type AccountEdges struct {
	// Member holds the value of the member edge.
	Member []*User `json:"member,omitempty"`
	// Categories holds the value of the categories edge.
	Categories []*Category `json:"categories,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [2]bool
}

// MemberOrErr returns the Member value or an error if the edge
// was not loaded in eager-loading.
func (e AccountEdges) MemberOrErr() ([]*User, error) {
	if e.loadedTypes[0] {
		return e.Member, nil
	}
	return nil, &NotLoadedError{edge: "member"}
}

// CategoriesOrErr returns the Categories value or an error if the edge
// was not loaded in eager-loading.
func (e AccountEdges) CategoriesOrErr() ([]*Category, error) {
	if e.loadedTypes[1] {
		return e.Categories, nil
	}
	return nil, &NotLoadedError{edge: "categories"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Account) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case account.FieldSheetsCredentials:
			values[i] = new([]byte)
		case account.FieldID:
			values[i] = new(sql.NullInt64)
		case account.FieldName, account.FieldSpreadsheetID:
			values[i] = new(sql.NullString)
		default:
			return nil, fmt.Errorf("unexpected column %q for type Account", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Account fields.
func (a *Account) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case account.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			a.ID = int(value.Int64)
		case account.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				a.Name = value.String
			}
		case account.FieldSheetsCredentials:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field sheets_credentials", values[i])
			} else if value != nil {
				a.SheetsCredentials = *value
			}
		case account.FieldSpreadsheetID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field spreadsheet_id", values[i])
			} else if value.Valid {
				a.SpreadsheetID = value.String
			}
		}
	}
	return nil
}

// QueryMember queries the "member" edge of the Account entity.
func (a *Account) QueryMember() *UserQuery {
	return (&AccountClient{config: a.config}).QueryMember(a)
}

// QueryCategories queries the "categories" edge of the Account entity.
func (a *Account) QueryCategories() *CategoryQuery {
	return (&AccountClient{config: a.config}).QueryCategories(a)
}

// Update returns a builder for updating this Account.
// Note that you need to call Account.Unwrap() before calling this method if this Account
// was returned from a transaction, and the transaction was committed or rolled back.
func (a *Account) Update() *AccountUpdateOne {
	return (&AccountClient{config: a.config}).UpdateOne(a)
}

// Unwrap unwraps the Account entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (a *Account) Unwrap() *Account {
	_tx, ok := a.config.driver.(*txDriver)
	if !ok {
		panic("ent: Account is not a transactional entity")
	}
	a.config.driver = _tx.drv
	return a
}

// String implements the fmt.Stringer.
func (a *Account) String() string {
	var builder strings.Builder
	builder.WriteString("Account(")
	builder.WriteString(fmt.Sprintf("id=%v, ", a.ID))
	builder.WriteString("name=")
	builder.WriteString(a.Name)
	builder.WriteString(", ")
	builder.WriteString("sheets_credentials=")
	builder.WriteString(fmt.Sprintf("%v", a.SheetsCredentials))
	builder.WriteString(", ")
	builder.WriteString("spreadsheet_id=")
	builder.WriteString(a.SpreadsheetID)
	builder.WriteByte(')')
	return builder.String()
}

// Accounts is a parsable slice of Account.
type Accounts []*Account

func (a Accounts) config(cfg config) {
	for _i := range a {
		a[_i].config = cfg
	}
}
