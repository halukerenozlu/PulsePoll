package routes

import (
	"reflect"
	"testing"
)

func TestSurveyModelIDHasUUIDDefaultTag(t *testing.T) {
	field, ok := reflect.TypeOf(surveyModel{}).FieldByName("ID")
	if !ok {
		t.Fatalf("surveyModel.ID field not found")
	}
	got := field.Tag.Get("gorm")
	want := "column:id;type:uuid;default:gen_random_uuid();primaryKey"
	if got != want {
		t.Fatalf("surveyModel.ID gorm tag = %q, want %q", got, want)
	}
}

func TestOptionDBIDHasUUIDDefaultTag(t *testing.T) {
	field, ok := reflect.TypeOf(optionDB{}).FieldByName("ID")
	if !ok {
		t.Fatalf("optionDB.ID field not found")
	}
	got := field.Tag.Get("gorm")
	want := "column:id;type:uuid;default:gen_random_uuid();primaryKey"
	if got != want {
		t.Fatalf("optionDB.ID gorm tag = %q, want %q", got, want)
	}
}
