// Code generated by ent, DO NOT EDIT.

package runtime

import (
	"time"

	"github.com/noamcattan/geni/ent/expense"
	"github.com/noamcattan/geni/ent/schema"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	expenseMixin := schema.Expense{}.Mixin()
	expenseHooks := schema.Expense{}.Hooks()
	expense.Hooks[0] = expenseHooks[0]
	expenseMixinFields0 := expenseMixin[0].Fields()
	_ = expenseMixinFields0
	expenseFields := schema.Expense{}.Fields()
	_ = expenseFields
	// expenseDescCreateTime is the schema descriptor for create_time field.
	expenseDescCreateTime := expenseMixinFields0[0].Descriptor()
	// expense.DefaultCreateTime holds the default value on creation for the create_time field.
	expense.DefaultCreateTime = expenseDescCreateTime.Default.(func() time.Time)
}

const (
	Version = "v0.11.4"                                         // Version of ent codegen.
	Sum     = "h1:grwVY0fp31BZ6oEo3YrXenAuv8VJmEw7F/Bi6WqeH3Q=" // Sum of ent codegen.
)
