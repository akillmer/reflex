package reflex

import (
	"testing"
)

func TestModelSet(t *testing.T) {
	var (
		s = struct {
			Name struct {
				First string
				Last  string
			}
			Location string
		}{}
		m = NewModel(&s)
	)

	m.Set("Name.First", "Andrew")
	m.Set("Name.Last", "Killmer")
	m.Set("Location", "Aiea, HI")

	if s.Name.First != "Andrew" {
		t.Errorf("Name.First should be 'Andrew', got '%s'", s.Name.First)
	}

	if s.Name.Last != "Killmer" {
		t.Errorf("Name.Last should be 'Killmer', got '%s'", s.Name.Last)
	}

	if s.Location != "Aiea, HI" {
		t.Errorf("Location should be 'Aiea, HI', got '%s'", s.Location)
	}
}

func TestModelMap(t *testing.T) {
	type member struct {
		Name  string
		Plays string
	}

	var (
		s = struct {
			Band   string
			Member []member
		}{}
		bm = []member{
			member{
				Name:  "Jimmy Page",
				Plays: "Guitar",
			},
			member{
				Name:  "Robert Plant",
				Plays: "Harmonica",
			},
			member{
				Name:  "John Bonham",
				Plays: "Drums",
			},
			member{
				Name:  "John Paul Jones",
				Plays: "Everything else",
			},
		}
		m = NewModel(&s)
	)

	m.Set("Band", "Led Zeppelin")

	for _, v := range bm {
		var kv = make(map[string]interface{})
		kv["Member.Name"] = v.Name
		kv["Member.Plays"] = v.Plays
		m.Map(kv)
	}

	if s.Band != "Led Zeppelin" {
		t.Errorf("Band should be 'Led Zeppelin', got '%s'", s.Band)
	}

	for i := range bm {
		if bm[i].Name != s.Member[i].Name {
			t.Errorf("member #%d should have name '%s', got '%s'", i, bm[i].Name, s.Member[i].Name)
		} else if bm[i].Plays != s.Member[i].Plays {
			t.Errorf("%s should play %s, got '%s'", bm[i].Name, bm[i].Plays, s.Member[i].Plays)
		}
	}
}

func TestManySlices(t *testing.T) {
	type movie struct {
		Title string
		Year  int
	}
	type series struct {
		Name   string
		Movies []movie
	}

	var (
		s        = []series{}
		triology = []series{
			series{
				Name: "Star Wars",
				Movies: []movie{
					movie{
						Title: "A New Hope",
						Year:  1977,
					},
					movie{
						Title: "The Empire Strikes Back",
						Year:  1980,
					},
					movie{
						Title: "Return of the Jedi",
						Year:  1983,
					},
				},
			},
			series{
				Name: "Indiana Jones",
				Movies: []movie{
					movie{
						Title: "Raiders of the Lost Ark",
						Year:  1981,
					},
					movie{
						Title: "The Temple of Doom",
						Year:  1984,
					},
					movie{
						Title: "The Last Crusade",
						Year:  1989,
					},
				},
			},
		} // phew
		m = NewModel(&s)
	)

	for _, tri := range triology {
		var kv = make(map[string]interface{})
		kv["Name"] = tri.Name
		kv["Movies"] = tri.Movies
		m.Map(kv)
	}

	if len(triology) != len(s) {
		t.Fatalf("there should be %d series, got %d", len(triology), len(s))
	}

	for i, tri := range triology {
		if tri.Name != s[i].Name {
			t.Errorf("triology should be '%s', got '%s'", tri.Name, s[i].Name)
			continue
		}
		for j, mov := range tri.Movies {
			if mov.Title != s[i].Movies[j].Title {
				t.Errorf("movie title should be '%s', got '%s'", mov.Title, s[i].Movies[j].Title)
			} else if mov.Year != s[i].Movies[j].Year {
				t.Errorf("%s came out in %d, got %d", mov.Title, mov.Year, s[i].Movies[j].Year)
			}
		}
	}
}
