package importer

import "testing"

func TestGetPrj(t *testing.T) {
	prj := `GEOGCS["GCS_WGS_1984",DATUM["D_WGS_1984",SPHEROID["WGS_1984",6378137,298.257223563]],PRIMEM["Greenwich",0],UNIT["Degree",0.017453292519943295]]`
	code, err := getPrj(prj)
	if err != nil {
		t.Fatal(err)
	}

	if code != 4326 {
		t.Fatalf("failed to get correct code, wanted %d, got %d", 4326, code)
	}
}
