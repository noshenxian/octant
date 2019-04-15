package component

import (
	"encoding/json"
	"sort"
)

// TableConfig is the contents of a Table
type TableConfig struct {
	Columns      []TableCol `json:"columns"`
	Rows         TableRows  `json:"rows"`
	EmptyContent string     `json:"emptyContent"`
}

// TableCol describes a column from a table. Accessor is the key this
// column will appear as in table rows, and must be unique within a table.
type TableCol struct {
	Name     string `json:"name"`
	Accessor string `json:"accessor"`
}

// TableRow is a row in table. Each key->value represents a particular column in the row.
type TableRow map[string]Component

// TablesRows are multiple rows.
type TableRows []TableRow

func (t TableRows) Sort(name string) error {
	sort.Slice(t, func(i, j int) bool {
		a := t[i][name]
		b := t[j][name]

		return a.String() < b.String()
	})

	return nil
}

func (t *TableRow) UnmarshalJSON(data []byte) error {
	*t = make(TableRow)

	x := map[string]TypedObject{}

	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	for k, v := range x {
		vc, err := v.ToComponent()
		if err != nil {
			return err
		}

		(*t)[k] = vc
	}

	return nil
}

// Table contains other Components
type Table struct {
	base
	Config TableConfig `json:"config"`
}

// NewTable creates a table component
func NewTable(title string, cols []TableCol) *Table {
	return &Table{
		base: newBase(typeTable, TitleFromString(title)),
		Config: TableConfig{
			Columns: cols,
		},
	}
}

// NewTableWithRows creates a table with rows.
func NewTableWithRows(title string, cols []TableCol, rows []TableRow) *Table {
	table := NewTable(title, cols)
	table.Add(rows...)
	return table
}

// NewTableCols returns a slice of table columns, each with name/accessor
// set according to the provided keys arguments.
func NewTableCols(keys ...string) []TableCol {
	if len(keys) == 0 {
		return nil
	}

	cols := make([]TableCol, len(keys))

	for i, k := range keys {
		cols[i].Name = k
		cols[i].Accessor = k
	}
	return cols
}

// IsEmpty returns true if there is one or more rows.
func (t *Table) IsEmpty() bool {
	return len(t.Config.Rows) < 1
}

// Add adds additional items to the tail of the table.
func (t *Table) Add(rows ...TableRow) {
	t.Config.Rows = append(t.Config.Rows, rows...)
}

// AddColumn adds a column to the table.
func (t *Table) AddColumn(name string) {
	t.Config.Columns = append(t.Config.Columns, TableCol{
		Name:     name,
		Accessor: name,
	})
}

type tableMarshal Table

// MarshalJSON implements json.Marshaler
func (t *Table) MarshalJSON() ([]byte, error) {
	m := tableMarshal(*t)
	m.Metadata.Type = typeTable
	return json.Marshal(&m)
}
