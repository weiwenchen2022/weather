package weather_test

import (
	"os"
	"testing"

	"weather"

	"github.com/google/go-cmp/cmp"
)

func TestFormatURL(t *testing.T) {
	t.Parallel()

	location := "London,uk"
	token := "dummy_token"
	want := "https://api.openweathermap.org/data/2.5/weather?q=" + location + "&APPID=" + token
	got := weather.FormatURL(location, token)
	if want != got {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestParseJSON(t *testing.T) {
	t.Parallel()

	f, err := os.Open("testdata/london.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	want := weather.Conditions{
		Summary:            "Drizzle",
		TemperatureCelsius: 7.17,
	}
	got := weather.ParseJSON(f)
	if diff := cmp.Diff(want, got); diff != "" {
		t.Error(diff)
	}
}

func TestLocationFromArgs(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		input       []string
		want        string
		errExpected bool
	}{
		{[]string{"London,", "UK"}, "London,UK", false},
		{[]string{"London"}, "London", false},
		{[]string{}, "", true},
	}

	for _, tc := range testcases {
		got, err := weather.LocationFromArgs(tc.input)
		if tc.errExpected && err == nil {
			t.Error("Error expected")
		}

		if diff := cmp.Diff(tc.want, got); diff != "" {
			t.Error(diff)
		}
	}
}

func TestFormat(t *testing.T) {
	cond := weather.Conditions{
		Summary:            "Drizzle",
		TemperatureCelsius: 7.17,
	}

	want := "Drizzle 7.2ºC"
	got := cond.String()

	if diff := cmp.Diff(want, got); diff != "" {
		t.Error(diff)
	}
}
