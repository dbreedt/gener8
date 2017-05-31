package main

import "testing"

func TestTimeConsuming(t *testing.T) {

	kws := "a,b,c"

	keywords, err := parseKws(kws)

	if err != nil {
		t.Error(err)
	}

	if len(*keywords) != 3 {
		t.Error("Found ", len(*keywords), " should be 3")
	}

	kws = "a,b,\"d,e\""

	keywords, err = parseKws(kws)

	if err != nil {
		t.Error(err)
	}

	if len(*keywords) != 3 {
		t.Error("Found ", len(*keywords), " should be 3")
	}

}
