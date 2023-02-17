package helpers

import "testing"

type testStructStub struct {
	name string
	age  int
}

func TestMapFunction1(t *testing.T) {
	slice := []string{"a", "b", "c"}

	result := Map(slice, func(s string) string {
		return s + "1"
	})

	if len(result) != len(slice) {
		t.Errorf("wrong result length")
	}

	for i := 0; i < len(result); i++ {
		if result[i] != slice[i]+"1" {
			t.Errorf("wrong result")
		}
	}
}

func TestMapFunction2(t *testing.T) {
	slice := []string{"abdf", "adfb", "dfdfc"}

	result := Map(slice, func(s string) testStructStub {
		return testStructStub{name: s + "g", age: 1 + int(s[0])}
	})

	if len(result) != len(slice) {
		t.Errorf("wrong result length")
	}

	expected := []testStructStub{testStructStub{
		name: "abdfg", age: 98},
		testStructStub{name: "adfbg", age: 98},
		testStructStub{name: "dfdfcg", age: 101}}

	for i := 0; i < len(result); i++ {
		if result[i].name != expected[i].name || result[i].age != expected[i].age {
			t.Errorf("wrong result")
		}
	}
}

func TestMapFunction3(t *testing.T) {

	slice := []string{}

	result := Map(slice, func(s string) testStructStub {
		return testStructStub{name: s + "g", age: 1 + int(s[0])}
	})

	if len(result) != 0 {
		t.Errorf("wrong result length")
	}
}

func TestStringFilter1(t *testing.T) {
	slice := []string{"a", "b", "c"}

	result := Filter(slice, func(s string) bool {
		return s == "a"
	})

	if len(result) != 1 {
		t.Errorf("wrong result length")
	}

	if result[0] != "a" {
		t.Errorf("wrong result")
	}
}

func TestStringFilter2(t *testing.T) {
	slice := []string{"a", "b", "c"}

	result := Filter(slice, func(s string) bool {
		return s != "d"
	})

	for i := 0; i < len(result); i++ {
		if result[i] != slice[i] {
			t.Errorf("wrong result")
		}
	}
}

func TestValidateEmailAddresses1(t *testing.T) {
	emails := []string{"test@gmail.com", "test@hello.sg", "h123@outlook.com"}

	err := ValidateEmailAddresses(emails)

	if err != nil {
		t.Errorf("All emails should be valid")
	}
}

func TestValidateEmailAddresses2(t *testing.T) {
	emails := []string{"hexi@@gmail.com", "test@hello.sg", "h123@outlook.com"}

	err := ValidateEmailAddresses(emails)

	if err == nil {
		t.Errorf("The first email should be invalid")
	}
}

func TestFindValidEmailsInText1(t *testing.T) {
	text := "test test yo @12-3@u.nus.edu lorem ipsum @dragon44@gmeow.com niceboy@example.com @hansomexample.com"

	emails := FindValidEmailsInText(text)

	if len(emails) != 2 {
		t.Errorf("wrong number of emails")
	}

}

func TestFindValidEmailsInText2(t *testing.T) {
	text := "hello there"

	emails := FindValidEmailsInText(text)

	if len(emails) != 0 {
		t.Errorf("wrong number of emails")
	}

}

func TestValidateEmailFormat(t *testing.T) {

	if ValidateEmailFormat("bulls-eye@haello.co") != nil {
		t.Errorf("email should be valid")
	}
}

func TestValidateEmailFormat2(t *testing.T) {

	if ValidateEmailFormat("123@u.nus.edu") != nil {
		t.Errorf("email should be valid")
	}
}
