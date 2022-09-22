package gedcom

import (
	"testing"
)

func TestSplitPersonalName(t *testing.T) {
	testCases := []struct {
		name string
		want ParsedName
	}{
		{
			name: "First Name Only",
			want: ParsedName{
				Full:     "First Name Only",
				Given:    "First Name Only",
				Nickname: "",
			},
		},
		{
			name: "First Name /Last Name/",
			want: ParsedName{
				Full:     "First Name Last Name",
				Given:    "First Name",
				Surname:  "Last Name",
				Nickname: "",
			},
		},
		{
			name: "First/ Last / ",
			want: ParsedName{
				Full:     "First/ Last",
				Given:    "First/ Last",
				Surname:  "",
				Nickname: "",
			},
		},
		{
			name: " First /Last ",
			want: ParsedName{
				Full:     "First Last",
				Given:    "First",
				Surname:  "Last",
				Nickname: "",
			},
		},
		{
			name: "First /Last/ II ",
			want: ParsedName{
				Full:     "First Last II",
				Given:    "First",
				Surname:  "Last",
				Suffix:   "II",
				Nickname: "",
			},
		},
		{
			name: "/Last/ Karl II",
			want: ParsedName{
				Full:     "Last Karl II",
				Surname:  "Last",
				Suffix:   "Karl II",
				Nickname: "",
			},
		},
		{
			name: "Жанна /Иванова (Д'Арк)/",
			want: ParsedName{
				Full:     "Жанна Иванова (Д'Арк)",
				Given:    "Жанна",
				Surname:  "Иванова (Д'Арк)",
				Nickname: "",
			},
		},
		{
			name: "First/Alt /Last/ II ",
			want: ParsedName{
				Full:     "First/Alt Last II",
				Given:    "First/Alt",
				Surname:  "Last",
				Suffix:   "II",
				Nickname: "",
			},
		},
		{
			name: "/Last/",
			want: ParsedName{
				Full:     "Last",
				Given:    "",
				Surname:  "Last",
				Suffix:   "",
				Nickname: "",
			},
		},
		{
			name: " /Last/",
			want: ParsedName{
				Full:     "Last",
				Given:    "",
				Surname:  "Last",
				Suffix:   "",
				Nickname: "",
			},
		},
		{
			name: "/Last/ Jr",
			want: ParsedName{
				Full:     "Last Jr",
				Given:    "",
				Surname:  "Last",
				Suffix:   "Jr",
				Nickname: "",
			},
		},
		{
			name: `Lydia /Fetters/Fletcher/`,
			want: ParsedName{
				Full:     "Lydia Fetters/Fletcher",
				Given:    "Lydia",
				Surname:  "Fetters/Fletcher",
				Suffix:   "",
				Nickname: "",
			},
		},
		{
			name: `John "Jack" /Bryan/`,
			want: ParsedName{
				Full:     "John Bryan",
				Given:    "John",
				Surname:  "Bryan",
				Suffix:   "",
				Nickname: "Jack",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := SplitPersonalName(tc.name)
			if actual != tc.want {
				t.Logf("original input: %q", tc.name)
				if actual.Full != tc.want.Full {
					t.Errorf("got Full=%q, wanted %q", actual.Full, tc.want.Full)
				}
				if actual.Given != tc.want.Given {
					t.Errorf("got Given=%q, wanted %q", actual.Given, tc.want.Given)
				}
				if actual.Surname != tc.want.Surname {
					t.Errorf("got Surname=%q, wanted %q", actual.Surname, tc.want.Surname)
				}
				if actual.Suffix != tc.want.Suffix {
					t.Errorf("got Suffix=%q, wanted %q", actual.Suffix, tc.want.Suffix)
				}
				if actual.Nickname != tc.want.Nickname {
					t.Errorf("got Nickname=%q, wanted %q", actual.Nickname, tc.want.Nickname)
				}
			}
		})
	}
}
