package schema

import (
	"context"
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"errors"
	geni "github.com/noamcattan/geni/ent"
	"github.com/noamcattan/geni/ent/hook"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"log"
	"time"
)

// Expense holds the schema definition for the Expense entity.
type Expense struct {
	ent.Schema
}

// Fields of the Expense.
func (Expense) Fields() []ent.Field {
	return []ent.Field{
		field.Float("amount"),
		field.String("description"),
	}
}

// Mixin of the Expense.
func (Expense) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.CreateTime{},
	}
}

// Edges of the Expense.
func (Expense) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("category", Category.Type).Ref("expenses"),
	}
}

func (Expense) Hooks() []ent.Hook {
	return []ent.Hook{
		hook.On(
			UpdateSheetsHook,
			ent.OpCreate,
		),
	}
}

func UpdateSheetsHook(next ent.Mutator) ent.Mutator {
	return hook.ExpenseFunc(func(ctx context.Context, mutation *geni.ExpenseMutation) (geni.Value, error) {
		newVal, nextErr := next.Mutate(ctx, mutation)
		if nextErr != nil {
			return newVal, nextErr
		}

		e, ok := newVal.(*geni.Expense)
		if !ok {
			log.Printf("error: %s", errors.New("bad value"))
			return newVal, nextErr
		}

		ca, err := e.QueryCategory().First(ctx)
		if err != nil {
			log.Printf("error: %s", err.Error())
			return newVal, nextErr
		}

		ac, err := e.QueryCategory().QueryAccount().First(ctx)
		if err != nil {
			log.Printf("error: %s", err.Error())
			return newVal, nextErr
		}

		srv, err := sheets.NewService(ctx, option.WithCredentialsJSON(ac.SheetsCredentials))
		if err != nil {
			log.Printf("error: %s", err.Error())
			return newVal, nextErr
		}

		spreadsheet, err := srv.Spreadsheets.Get(ac.SpreadsheetID).Do()
		if err != nil {
			log.Printf("error: %s", err.Error())
			return newVal, nextErr
		}

		var sheet *sheets.Sheet

		now := time.Now()
		sheetTitle := now.Format("January 2006")

		for _, s := range spreadsheet.Sheets {
			if s.Properties.Title == sheetTitle {
				sheet = s
				break
			}
		}

		if sheet == nil {
			// create sheet
			req := sheets.Request{
				AddSheet: &sheets.AddSheetRequest{
					Properties: &sheets.SheetProperties{
						Title: sheetTitle,
					},
				},
			}

			rbb := &sheets.BatchUpdateSpreadsheetRequest{
				Requests: []*sheets.Request{&req},
			}
			_, err = srv.Spreadsheets.BatchUpdate(ac.SpreadsheetID, rbb).Context(context.Background()).Do()
			if err != nil {
				log.Printf("error: %s", err.Error())
				return newVal, nextErr
			}
			row := &sheets.ValueRange{
				Values: [][]interface{}{{"time", "category", "amount", "description"}},
			}
			_, err = srv.Spreadsheets.Values.
				Append(ac.SpreadsheetID, sheetTitle, row).
				ValueInputOption("USER_ENTERED").
				InsertDataOption("INSERT_ROWS").
				Context(ctx).Do()
			if err != nil {
				log.Printf("error: %s", err.Error())
				return newVal, nextErr
			}

		}

		row := &sheets.ValueRange{
			Values: [][]interface{}{{
				e.CreateTime.Format("02/01/2006 15:04:05"),
				ca.Name,
				e.Amount,
				e.Description}},
		}

		_, err = srv.Spreadsheets.Values.
			Append(ac.SpreadsheetID, sheetTitle, row).
			ValueInputOption("USER_ENTERED").
			InsertDataOption("INSERT_ROWS").
			Context(ctx).Do()
		if err != nil {
			log.Printf("error: %s", err.Error())
			return newVal, nextErr
		}

		return newVal, nextErr
	})
}
