package geo

import (
	"encoding/json"
	"os"
	"testing"
)

// Ensures that the library can detect if a point is in a polygon.
// Uses Brunei and the capital of Brunei as a set of test points.
func TestPointInPolygon(t *testing.T) {
	brunei, err := polygonFromFile("test/data/brunei.json")
	if err != nil {
		t.Error("brunei json file failed to parse: ", err)
	}

	point := Point{lng: 114.9480600, lat: 4.9402900}
	if !brunei.Contains(point) {
		t.Error("Expected the capital of Brunei to be in Brunei, but it wasn't.")
	}
}

// Ensures that the polygon logic can correctly identify if a polygon does not contain a point.
// Uses Brunei, Seattle, and a point directly outside of Brunei limits as test points.
func TestPointNotInPolygon(t *testing.T) {
	brunei, err := polygonFromFile("test/data/brunei.json")
	if err != nil {
		t.Error("brunei json file failed to parse: ", err)
	}

	// Seattle, WA should not be inside of Brunei
	point := NewPoint(47.45, 122.30)
	if brunei.Contains(point) {
		t.Error("Seattle, WA [47.45, 122.30] should not be inside of Brunei")
	}

	// A point just outside of the successful bounds in Brunei
	// Should not be contained in the Polygon
	precision := NewPoint(114.659596, 4.007636)
	if brunei.Contains(precision) {
		t.Error("A point just outside of Brunei should not be contained in the Polygon")
	}
}

// Ensures that a point can be contained in a complex polygon (e.g. a donut)
// This particular Polygon has a hole in it.
func TestPointInPolygonWithHole(t *testing.T) {
	nsw, err := polygonFromFile("test/data/nsw.json")
	if err != nil {
		t.Error("nsw json file failed to parse: ", err)
	}

	act, err := polygonFromFile("test/data/act.json")
	if err != nil {
		t.Error("act json file failed to parse: ", err)
	}

	// Look at two contours
	canberra := Point{lng: 149.128684300000030000, lat: -35.2819998}
	isnsw := nsw.Contains(canberra)
	isact := act.Contains(canberra)
	if !isnsw && !isact {
		t.Error("Canberra should be in NSW and also in the sub-contour ACT state")
	}

	// Using NSW as a multi-contour polygon
	nswmulti := Polygon{}
	for _, p := range nsw.Points() {
		nswmulti = nswmulti.Add(p)
	}

	for _, p := range act.Points() {
		nswmulti = nswmulti.Add(p)
	}

	isnsw = nswmulti.Contains(canberra)
	if isnsw {
		t.Error("Canberra should not be in NSW as it falls in the donut contour of the ACT")
	}

	sydney := Point{lng: 151.209, lat: -33.866}

	if !nswmulti.Contains(sydney) {
		t.Error("Sydney should be in NSW")
	}

	losangeles := Point{lng: 118.28333, lat: 34.01667}
	isnsw = nswmulti.Contains(losangeles)

	if isnsw {
		t.Error("Los Angeles should not be in NSW")
	}

}

// Ensures that jumping over the equator and the greenwich meridian
// Doesn't give us any false positives or false negatives
func TestEquatorGreenwichContains(t *testing.T) {
	point1 := NewPoint(0.0, 0.0)
	point2 := NewPoint(0.1, 0.1)
	point3 := NewPoint(0.1, -0.1)
	point4 := NewPoint(-0.1, -0.1)
	point5 := NewPoint(-0.1, 0.1)
	polygon, err := polygonFromFile("test/data/equator_greenwich.json")

	if err != nil {
		t.Errorf("error parsing polygon %v", err)
	}

	if !polygon.Contains(point1) {
		t.Errorf("Should contain middle point of earth")
	}

	if !polygon.Contains(point2) {
		t.Errorf("Should contain point %v", point2)
	}

	if !polygon.Contains(point3) {
		t.Errorf("Should contain point %v", point3)
	}

	if !polygon.Contains(point4) {
		t.Errorf("Should contain point %v", point4)
	}

	if !polygon.Contains(point5) {
		t.Errorf("Should contain point %v", point5)
	}
}

// A test struct used to encapsulate and
// Unmarshal JSON into.
type testPoints struct {
	Points []Point
}

// Opens a JSON file and unmarshals the data into a Polygon
func polygonFromFile(filename string) (Polygon, error) {
	p := Polygon{}
	file, err := os.Open(filename)
	if err != nil {
		return Polygon{}, err
	}

	points := new(testPoints)
	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(&points); err != nil {
		return Polygon{}, err
	}

	for _, point := range points.Points {
		p = p.Add(point)
	}

	return p, nil
}

func TestPolygonNotClosed(t *testing.T) {
	points := []Point{
		NewPoint(0, 0), NewPoint(0, 1), //! NewPoint(1, 1), NewPoint(1, 0), if stricter rules for a closed polygon apply
	}
	poly := NewPolygon(points)
	if poly.IsClosed() {
		t.Error("Nope! The polygon is not closed!")
	}
}

func TestPolygonClosed(t *testing.T) {
	points := []Point{
		NewPoint(0, 0), NewPoint(0, 1), NewPoint(1, 1), NewPoint(1, 0), NewPoint(0, 0),
	}
	poly := NewPolygon(points)
	if !poly.IsClosed() {
		t.Error("Nope! The polygon is closed!")
	}
}

type testPoint struct {
	P        Point
	Expected bool
}

// ntp = newTestPoint
func ntp(lat float64, lng float64, expectedOutput bool) (t testPoint) {
	t.P = NewPoint(lat, lng)
	t.Expected = expectedOutput
	return t
}
func TestPolygonSimple(t *testing.T) {
	// A closed polygon - however the IsClosed() does not work properly.
	polygonPoints := []Point{
		NewPoint(0, 0), NewPoint(0, 1), NewPoint(1, 1), NewPoint(1, 0), NewPoint(0, 0),
	}

	testers := []testPoint{
		ntp(0, 0, true),
		ntp(0.5, 0.5, true),
		ntp(9, 9, false),
		ntp(0.5, 0, true),
		ntp(0.2, 0.4, true),
		//ntp(1, 0.9, true),
		ntp(1, 1.1, false),
		//ntp(0.9, 1, true),
		ntp(1, 1, true),
		ntp(0.00000000000000000000009, 0.00000000000000000000000000000001, true),
	}

	polygon := NewPolygon(polygonPoints)

	for _, q := range testers {
		res := polygon.Contains(q.P)
		if res != q.Expected {
			t.Log("x: ", q.P.lat, "\ty: ", q.P.lng, "\tExpected: ", q.Expected, "\tGot: ", res)
			t.Fail()
		}
	}

}

func TestPolygon_Contains(t *testing.T) {
	testPoint1 := NewPoint(32.349556, -86.257392)
	testPoint2 := NewPoint(33.245495, -87.523993)
	testPoint3 := NewPoint(29.361532, -100.882520)
	testPoint4 := NewPoint(29.361530, -100.882169)
	testPoint5 := NewPoint(29.360826, -100.883000)
	testPoint6 := NewPoint(29.360826, -100.882000)
	testPoint7 := NewPoint(34.709741, -86.669305)
	type fields struct {
		points []Point
	}
	type args struct {
		point Point
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{"Simple incorrect test", fields{[]Point{{0.0, 0.0}, {2.0, 0.0}, {2.0, 2.0}, {0.0, 2.0}, {0.0, 0.0}}}, args{NewPoint(3.0, 1.0)}, false},
		{"testPoint1 is not in block 011010028001000", fields{[]Point{
			{32.352226, -86.263596},
			{32.352358, -86.263566},
			{32.352367, -86.263383},
			{32.352370, -86.262824},
			{32.352382, -86.262643},
			{32.349556, -86.260067},
			{32.343598, -86.253016},
			{32.343214, -86.253470},
			{32.343803, -86.254162},
			{32.344282, -86.254726},
			{32.345042, -86.255614},
			{32.345059, -86.255637},
			{32.347299, -86.258256},
			{32.347734, -86.258773},
			{32.347880, -86.258952},
			{32.347924, -86.259024},
			{32.347954, -86.259090},
			{32.347982, -86.259172},
			{32.347999, -86.259258},
			{32.348007, -86.259358},
			{32.348002, -86.259772},
			{32.347988, -86.260983},
			{32.348053, -86.260985},
			{32.349114, -86.261014},
			{32.349226, -86.261027},
			{32.349375, -86.261056},
			{32.349508, -86.261101},
			{32.349615, -86.261151},
			{32.349775, -86.261253},
			{32.349840, -86.261303},
			{32.349916, -86.261365},
			{32.349998, -86.261451},
			{32.350949, -86.262569},
			{32.351430, -86.263142},
			{32.351514, -86.263218},
			{32.351630, -86.263345},
			{32.351719, -86.263428},
			{32.351806, -86.263491},
			{32.351920, -86.263546},
			{32.352044, -86.263581},
			{32.352167, -86.263596},
			{32.352226, -86.263596}}}, args{testPoint1}, false},
		{"testPoint1 is not in block 011010033014007", fields{[]Point{
			{32.352389, -86.262545},
			{32.352403, -86.261315},
			{32.352412, -86.260666},
			{32.352396, -86.260458},
			{32.352399, -86.260006},
			{32.352410, -86.258055},
			{32.352186, -86.258056},
			{32.351362, -86.258042},
			{32.350316, -86.258025},
			{32.349272, -86.258009},
			{32.349114, -86.258006},
			{32.349030, -86.257998},
			{32.348951, -86.257976},
			{32.348871, -86.257936},
			{32.348841, -86.257914},
			{32.348768, -86.257845},
			{32.348415, -86.257428},
			{32.348357, -86.257338},
			{32.348308, -86.257241},
			{32.348264, -86.257126},
			{32.348233, -86.257017},
			{32.348216, -86.256917},
			{32.348205, -86.256815},
			{32.348207, -86.256523},
			{32.348022, -86.256510},
			{32.347886, -86.256477},
			{32.347805, -86.256445},
			{32.347746, -86.256416},
			{32.347683, -86.256377},
			{32.347628, -86.256335},
			{32.347546, -86.256261},
			{32.347482, -86.256192},
			{32.346998, -86.255624},
			{32.346804, -86.255394},
			{32.346044, -86.254499},
			{32.345329, -86.253659},
			{32.344530, -86.252721},
			{32.344264, -86.252406},
			{32.344194, -86.252312},
			{32.343598, -86.253016},
			{32.349556, -86.260067},
			{32.352382, -86.262643},
			{32.352389, -86.262545},
		},
		}, args{testPoint1}, false},
		{"testPoint2 is not in block 011250104062021", fields{
			[]Point{
				{33.244945, -87.525365},
				{33.245156, -87.525299},
				{33.245345, -87.525172},
				{33.245573, -87.524987},
				{33.245550, -87.524926},
				{33.245559, -87.524827},
				{33.245607, -87.524772},
				{33.245495, -87.524693},
				{33.243568, -87.521510},
				{33.243527, -87.521451},
				{33.243480, -87.521369},
				{33.242871, -87.521996},
				{33.242792, -87.522083},
				{33.242941, -87.522309},
				{33.243155, -87.522655},
				{33.243339, -87.522998},
				{33.243479, -87.523201},
				{33.243808, -87.523755},
				{33.243912, -87.523956},
				{33.244075, -87.524203},
				{33.244237, -87.524484},
				{33.244301, -87.524642},
				{33.244343, -87.524813},
				{33.244362, -87.524989},
				{33.244366, -87.525167},
				{33.244360, -87.525379},
				{33.244945, -87.525365},
			},
		}, args{testPoint2}, false},
		{"testPoint3 is not in Texas block 484659506023021", fields{
			[]Point{
				{29.360755, -100.883863},
				{29.361532, -100.882980},
				{29.360826, -100.882169},
				{29.360046, -100.883046},
				{29.360755, -100.883863},
			},
		}, args{testPoint3}, false},
		{"Previous flipped", fields{
			[]Point{
				{-100.883863, 29.360755},
				{-100.882980, 29.361532},
				{-100.882169, 29.360826},
				{-100.883046, 29.360046},
				{-100.883863, 29.360755},
			},
		}, args{NewPoint(testPoint3.lng, testPoint3.lat)}, false},
		{"testPoint4 is not in Texas block 484659506023021", fields{
			[]Point{
				{29.360755, -100.883863},
				{29.361532, -100.882980},
				{29.360826, -100.882169},
				{29.360046, -100.883046},
				{29.360755, -100.883863},
			},
		}, args{testPoint4}, false},
		{"testPoint6 is not in Texas block 484659506023019", fields{
			[]Point{
				{29.360755, -100.883863},
				{29.361532, -100.882980},
				{29.360826, -100.882169},
				{29.360046, -100.883046},
				{29.360755, -100.883863},
			},
		}, args{testPoint6}, false},
		{"testPoint7 is not in Alabama block 010890014042002", fields{
			[]Point{
				{34.710450, -86.671007},
				{34.711249, -86.671000},
				{34.712443, -86.670998},
				{34.712665, -86.670993},
				{34.713262, -86.670980},
				{34.713240, -86.670280},
				{34.713241, -86.670149},
				{34.713246, -86.669575},
				{34.713250, -86.669169},
				{34.713248, -86.669037},
				{34.713245, -86.668810},
				{34.713012, -86.668803},
				{34.712973, -86.669617},
				{34.712467, -86.670034},
				{34.712238, -86.670065},
				{34.711941, -86.670066},
				{34.711681, -86.670106},
				{34.711705, -86.670533},
				{34.711626, -86.670564},
				{34.710773, -86.670595},
				{34.710531, -86.670595},
				{34.710294, -86.670595},
				{34.710090, -86.670480},
				{34.709506, -86.670167},
				{34.709419, -86.670102},
				{34.709366, -86.670024},
				{34.709363, -86.669908},
				{34.709517, -86.669317},
				{34.709670, -86.668895},
				{34.709755, -86.668714},
				{34.709860, -86.668685},
				{34.710118, -86.668672},
				{34.710565, -86.668673},
				{34.711325, -86.668657},
				{34.712334, -86.668666},
				{34.712476, -86.668685},
				{34.712540, -86.668735},
				{34.712544, -86.668882},
				{34.712706, -86.668713},
				{34.712704, -86.668586},
				{34.712686, -86.668351},
				{34.712908, -86.668281},
				{34.713237, -86.668282},
				{34.713224, -86.666978},
				{34.713219, -86.666216},
				{34.713215, -86.665502},
				{34.712941, -86.665459},
				{34.712682, -86.665461},
				{34.712353, -86.665464},
				{34.712126, -86.665497},
				{34.712003, -86.665539},
				{34.711647, -86.665433},
				{34.711515, -86.665416},
				{34.711118, -86.665438},
				{34.710567, -86.665297},
				{34.710420, -86.665260},
				{34.710392, -86.665360},
				{34.710125, -86.666345},
				{34.709319, -86.669305},
				{34.709062, -86.670252},
				{34.708956, -86.670643},
				{34.708907, -86.670823},
				{34.708870, -86.670960},
				{34.708857, -86.671006},
				{34.710037, -86.671022},
				{34.710450, -86.671007},
			},
		}, args{testPoint7}, false},
		{"Simple correct test", fields{[]Point{{0.0, 0.0}, {2.0, 0.0}, {2.0, 2.0}, {0.0, 2.0}, {0.0, 0.0}}}, args{NewPoint(1.0, 1.0)}, true},
		{"testPoint1 is in this block", fields{
			[]Point{
				{32.350316, -86.258025},
				{32.350354, -86.254757},
				{32.349308, -86.254739},
				{32.349299, -86.255634},
				{32.349272, -86.258009},
				{32.350316, -86.258025},
			},
		}, args{testPoint1}, true},
		{"testPoint2 is in block 011250104062016", fields{
			[]Point{
				{33.246460, -87.525798},
				{33.246557, -87.525795},
				{33.246628, -87.525781},
				{33.246705, -87.525747},
				{33.246793, -87.525690},
				{33.246879, -87.525624},
				{33.246936, -87.525549},
				{33.247078, -87.525304},
				{33.247283, -87.524910},
				{33.247344, -87.524787},
				{33.247819, -87.523839},
				{33.248018, -87.523486},
				{33.248070, -87.523361},
				{33.248078, -87.523309},
				{33.248090, -87.523251},
				{33.248104, -87.523184},
				{33.248118, -87.523110},
				{33.248129, -87.523029},
				{33.248137, -87.522856},
				{33.248136, -87.522765},
				{33.248134, -87.522675},
				{33.248132, -87.522584},
				{33.248129, -87.522493},
				{33.248124, -87.522402},
				{33.248122, -87.522311},
				{33.248120, -87.522220},
				{33.248110, -87.522041},
				{33.248103, -87.521952},
				{33.248107, -87.521790},
				{33.248105, -87.521709},
				{33.248102, -87.521627},
				{33.248098, -87.521540},
				{33.248095, -87.521464},
				{33.248092, -87.521379},
				{33.248092, -87.521295},
				{33.248095, -87.521212},
				{33.248089, -87.521030},
				{33.248086, -87.520939},
				{33.248071, -87.520675},
				{33.248066, -87.520592},
				{33.247547, -87.520675},
				{33.247059, -87.520740},
				{33.246504, -87.520815},
				{33.246375, -87.520843},
				{33.246254, -87.520932},
				{33.246152, -87.521027},
				{33.246194, -87.521211},
				{33.246240, -87.521626},
				{33.246280, -87.522054},
				{33.246313, -87.522489},
				{33.246254, -87.522652},
				{33.246102, -87.522841},
				{33.245921, -87.523037},
				{33.245807, -87.523035},
				{33.245734, -87.523016},
				{33.245669, -87.522977},
				{33.245568, -87.522860},
				{33.245354, -87.522511},
				{33.244509, -87.521099},
				{33.244217, -87.520632},
				{33.244034, -87.520809},
				{33.243557, -87.521295},
				{33.243670, -87.521486},
				{33.243889, -87.521839},
				{33.244666, -87.523120},
				{33.245375, -87.524301},
				{33.245581, -87.524628},
				{33.245607, -87.524772},
				{33.245559, -87.524827},
				{33.245550, -87.524926},
				{33.245573, -87.524987},
				{33.245617, -87.525034},
				{33.245677, -87.525041},
				{33.245748, -87.525011},
				{33.245832, -87.525072},
				{33.245883, -87.525192},
				{33.245908, -87.525251},
				{33.245891, -87.525357},
				{33.246024, -87.525564},
				{33.246115, -87.525667},
				{33.246203, -87.525724},
				{33.246275, -87.525769},
				{33.246357, -87.525795},
				{33.246460, -87.525798},
			},
		}, args{testPoint2}, true},
		{"testPoint3 is in Texas block 484659506023019", fields{
			[]Point{
				{29.361635, -100.882863},
				{29.362469, -100.881915},
				{29.361767, -100.881112},
				{29.360930, -100.882053},
				{29.361635, -100.882863},
			},
		}, args{testPoint3}, true},
		{"testPoint5 is in Texas block 484659506023019", fields{
			[]Point{
				{29.360755, -100.883863},
				{29.361532, -100.882980},
				{29.360826, -100.882169},
				{29.360046, -100.883046},
				{29.360755, -100.883863},
			},
		}, args{testPoint5}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Polygon{
				points: tt.fields.points,
			}
			if got := p.Contains(tt.args.point); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}
