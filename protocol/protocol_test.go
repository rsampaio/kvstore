package protocol

import (
	"errors"
	"reflect"
	"testing"
)

type TestCase struct {
	Name         string
	Text         string
	Parsed       *Protocol
	ParsingError error
}

func TestProtocolParse(t *testing.T) {
	for _, tt := range []TestCase{
		{
			Name: "TestSetSuccess",
			Text: "SET foo 3",
			Parsed: &Protocol{
				Command:       "SET",
				Args:          []string{"foo", "3"},
				ReceivesValue: true,
			},
		},
		{
			Name:         "TestSetError",
			Text:         "SET foo",
			ParsingError: errors.New("set invalid arguments"),
		},
		{
			Name:         "TestSetInvalidSize",
			Text:         "SET foo a",
			ParsingError: errors.New("set invalid size"),
		},
		{
			Name: "TestGetSucess",
			Text: "GET foo",
			Parsed: &Protocol{
				Command: "GET",
				Args:    []string{"foo"},
			},
		},
		{
			Name:         "TestGetInvalidArgument",
			Text:         "GET",
			ParsingError: errors.New("get invalid arguments"),
		},
		{
			Name: "TestDeleteSucess",
			Text: "DELETE foo",
			Parsed: &Protocol{
				Command: "DELETE",
				Args:    []string{"foo"},
			},
		},
		{
			Name:         "TestDeleteInvalidArguments",
			Text:         "DELETE",
			ParsingError: errors.New("delete invalid arguments"),
		},
		{
			Name:         "TestInvalidCommand",
			Text:         "NONE",
			ParsingError: errors.New("invalid command"),
		},
	} {
		t.Run(tt.Name, func(t *testing.T) {
			p := &Protocol{}
			err := p.Parse(tt.Text)
			if err != nil {
				if err.Error() != tt.ParsingError.Error() {
					t.Errorf("got %v, expected %v", err, tt.ParsingError)
				}
			}

			if tt.Parsed != nil {
				if !reflect.DeepEqual(p, tt.Parsed) {
					t.Errorf("got %#v, expected %#v", p, tt.Parsed)
				}
			}
		})
	}
}
