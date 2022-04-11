package sql

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRenderIdentifier(t *testing.T) {

	r := NewRenderer(nil, NewTokenPair("L", "R").Wrap, DefaultUnwrappedIdentifiers)

	var quote = []string{
		"f o o",
		"3two_one",
		"wello$horld",
		"a/b",
	}

	for _, name := range quote {
		require.Equal(t, "L"+name+"R", r.Render(name))
	}

	var noQuote = []string{
		"FOO",
		"foo",
		"没有双引号",
		"_bar",
		"one_2_3",
	}
	for _, name := range noQuote {
		require.Equal(t, name, r.Render(name))
	}

}

func TestRenderSQLQuoteValues(t *testing.T) {

	r := NewRenderer(DefaultQuoteSanitizer, SingleQuotesWrapper(), nil)

	var testCases = map[string]string{
		"foo":            "'foo'",
		"he's 'bouta go": "'he''s ''bouta go'",
		"'moar quotes'":  "'''moar quotes'''",
		"":               "''",
	}
	for input, expected := range testCases {
		var actual = r.Render(input)
		require.Equal(t, expected, actual)
	}
}

func TestRenderComplexWrapper(t *testing.T) {

	r := NewRenderer(strings.NewReplacer("/", "_").Replace, func(text string) string {
		wrapper := SingleQuotesWrapper()
		terms := strings.Split(text, ".")
		for i, t := range terms {
			terms[i] = wrapper(t)
		}
		return strings.Join(terms, ".")
	}, nil)

	var testCases = map[string]string{
		"table1":                       "'table1'",
		"namespace.table2":             "'namespace'.'table2'",
		"bigole.namespace.table/3/yay": "'bigole'.'namespace'.'table_3_yay'",
		"":                             "''",
	}
	for input, expected := range testCases {
		var actual = r.Render(input)
		require.Equal(t, expected, actual)
	}
}

func TestRenderComment(t *testing.T) {

	testComment := fmt.Sprintf(
		"Generated by Flow for materializing collection '%s'\nto table: %s",
		"test",
		"target",
	)

	require.Equal(t, `-- Generated by Flow for materializing collection 'test'
-- to table: target
`, LineCommentRenderer().Render(testComment))

}

func TestCreateTableWithComments(t *testing.T) {

	endpoint := NewStdEndpoint(nil, nil, SQLiteSQLGenerator(), FlowTables{})

	tableRender, err := endpoint.CreateTableStatement(FlowCheckpointsTable("test"))
	require.Nil(t, err)

	require.Equal(t, `-- This table holds Flow processing checkpoints used for exactly-once processing of materializations
CREATE TABLE IF NOT EXISTS test (
	-- The name of the materialization.
	materialization TEXT NOT NULL,
	-- The inclusive lower-bound key hash covered by this checkpoint.
	key_begin INTEGER NOT NULL,
	-- The inclusive upper-bound key hash covered by this checkpoint.
	key_end INTEGER NOT NULL,
	-- This nonce is used to uniquely identify unique process assignments of a shard and prevent them from conflicting.
	fence INTEGER NOT NULL,
	-- Checkpoint of the Flow consumer shard, encoded as base64 protobuf.
	checkpoint TEXT,

	PRIMARY KEY(materialization, key_begin, key_end)
);`, tableRender)

}
