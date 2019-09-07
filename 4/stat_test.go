package stat

import (
	"reflect"
	"testing"
)

func TestEmpty(t *testing.T) {
	result := Top10("")
	if len(result) > 0 {
		t.Errorf("Expected empty slice, not %s (len = %d)\n", result, len(result))
	}
}

func TestTopTiny(t *testing.T) {
	result := Top10("Привет, мир!")
	expected := []string{"мир", "привет"}
	assertSlicesEquals(t, expected, result)
}

func TestTopJack(t *testing.T) {
	result := Top10(`Вот дом,
Который построил Джек.
А это пшеница,
Которая в темном чулане хранится
В доме,
Который построил Джек.
А это веселая птица-синица,
Которая часто ворует пшеницу,
Которая в темном чулане хранится
В доме,
Который построил Джек.`)
	expected := []string{"джек", "которая", "который", "построил", "доме", "темном", "хранится", "чулане", "это", "веселая"}
	assertSlicesEquals(t, expected, result)
}

func assertSlicesEquals(t *testing.T, expected []string, tested []string) {
	if !reflect.DeepEqual(expected, tested) {
		t.Errorf("Expected: %s, result %s", expected, tested)
	}
}
