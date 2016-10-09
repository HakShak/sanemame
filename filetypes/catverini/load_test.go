package catverini

import (
	"io/ioutil"
	"os"
	"testing"
)

const category = "[Category]\n"
const singleCategory = "machine1=SingleCategory\n"
const subCategory = "machine2=Primary Category / Sub Category\n"
const comment = "; Some silly comment\n"

func TestLoad_Sane(t *testing.T) {
	category := category + singleCategory
	category += subCategory
	category += comment
	err := ioutil.WriteFile("test.ini", []byte(category), 0644)
	if err != nil {
		panic(err)
	}
	defer os.Remove("test.ini")

	var machines []string
	var categories []Category

	Load("test.ini", EntryRead(func(machine string, category *Category) error {
		machines = append(machines, machine)
		categories = append(categories, *category)

		return nil
	}))

	expectedMachine := "machine1"
	if machines[0] != expectedMachine {
		t.Fatalf("Expected Machines: %s was not %s", machines[0], expectedMachine)
	}

	expectedRawCategory := "SingleCategory"
	if categories[0].Raw != expectedRawCategory {
		t.Fatalf("Expected Raw Category: %s was not %s",
			categories[0].Raw,
			expectedRawCategory)
	}
}
