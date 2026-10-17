package main

import "testing"

func TestAggregate(t *testing.T) {
	tracts := []Tract{
		{"Vermont", "Addison", 100}, {"Vermont", "Addison", 50}, {"Vermont", "Windsor", 7},
		{"Idaho", "Ada", 1000},
	}
	got := Aggregate(tracts)
	if got["Vermont"]["Addison"] != (CountyStats{Tracts: 2, Pop: 150}) {
		t.Errorf("Addison = %+v", got["Vermont"]["Addison"])
	}
	if got["Vermont"]["Windsor"].Pop != 7 || got["Idaho"]["Ada"].Tracts != 1 {
		t.Errorf("stats = %+v", got)
	}
}

func TestReadTractsFixture(t *testing.T) {
	tracts, err := ReadTracts("census.xlsx")
	if err != nil {
		t.Fatal(err)
	}
	if len(tracts) < 20 {
		t.Fatalf("only %d tracts read", len(tracts))
	}
	stats := Aggregate(tracts)
	if len(stats) != 3 || len(stats["Vermont"]) != 4 {
		t.Fatalf("states/counties = %d/%d", len(stats), len(stats["Vermont"]))
	}
}
