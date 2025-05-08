//go:build !bench
// +build !bench

package hw10programoptimization

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetDomainStat(t *testing.T) {
	data := `{"Id":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"aliquid_qui_ea@Browsedrive.gov","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}
{"Id":2,"Name":"Jesse Vasquez","Username":"qRichardson","Email":"mLynch@broWsecat.com","Phone":"9-373-949-64-00","Password":"SiZLeNSGn","Address":"Fulton Hill 80"}
{"Id":3,"Name":"Clarence Olson","Username":"RachelAdams","Email":"RoseSmith@Browsecat.com","Phone":"988-48-97","Password":"71kuz3gA5w","Address":"Monterey Park 39"}
{"Id":4,"Name":"Gregory Reid","Username":"tButler","Email":"5Moore@Teklist.net","Phone":"520-04-16","Password":"r639qLNu","Address":"Sunfield Park 20"}
{"Id":5,"Name":"Janice Rose","Username":"KeithHart","Email":"nulla@Linktype.com","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}`

	t.Run("find 'com'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{
			"browsecat.com": 2,
			"linktype.com":  1,
		}, result)
	})

	t.Run("find 'gov'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "gov")
		require.NoError(t, err)
		require.Equal(t, DomainStat{"browsedrive.gov": 1}, result)
	})

	t.Run("find 'unknown'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "unknown")
		require.NoError(t, err)
		require.Equal(t, DomainStat{}, result)
	})
}

func TestGetDomainStatCustom(t *testing.T) {
	testCases := []struct {
		name    string
		data    string
		domain  string
		want    DomainStat
		wantErr bool
	}{
		{
			name: "find com",
			data: `{"Id":1,"Email":"arotiros@Browsedrive.gov"}
{"Id":2,"Email":"mLynch@broWsecat.com"}
{"Id":3,"Email":"RoseSmith@Browsecat.com"}`,
			domain: "com",
			want:   DomainStat{"browsecat.com": 2},
		},
		{
			name: "find gov",
			data: `{"Id":1,"Email":"arotiros@Browsedrive.gov"}
{"Id":2,"Email":"user@domain.com"}`,
			domain: "gov",
			want:   DomainStat{"browsedrive.gov": 1},
		},
		{
			name:   "unknown domain",
			data:   `{"Id":1,"Email":"foo@bar.baz"}`,
			domain: "com",
			want:   DomainStat{},
		},
		{
			name:   "empty input",
			data:   ``,
			domain: "com",
			want:   DomainStat{},
		},
		{
			name:    "invalid JSON",
			data:    `{"Id":1,"Email":"a@b.com"`,
			domain:  "com",
			wantErr: true,
		},
		{
			name:   "missing email field",
			data:   `{"Id":1,"Username":"noemail"}`,
			domain: "com",
			want:   DomainStat{},
		},
		{
			name:   "without @",
			data:   `{"Id":1,"Email":"invalid.email.com"}`,
			domain: "com",
			want:   DomainStat{},
		},
		{
			name:   "@ at end 1",
			data:   `{"Id":1,"Email":"user@"}`,
			domain: "com",
			want:   DomainStat{},
		},
		{
			name:   "@ at end 2",
			data:   `{"Id":1,"Email":"user.com@"}`,
			domain: "com",
			want:   DomainStat{},
		},
		{
			name: "case-insensitive domain",
			data: `{"Id":1,"Email":"aiwbd@Example.COM"}
{"Id":2,"Email":"oawdn@example.com"}`,
			domain: "COM",
			want:   DomainStat{"example.com": 2},
		},
		{
			name: "extra newlines",
			data: "\n{" + `"Id":1,"Email":"a@foo.biz"}` + "\n\n" +
				`{"Id":2,"Email":"b@bar.biz"}` + "\n",
			domain: "biz",
			want:   DomainStat{"foo.biz": 1, "bar.biz": 1},
		},
		{
			name: "multiple same domain",
			data: strings.Repeat(`{"Id":1,"Email":"u@dup.com"}
`, 5),
			domain: "com",
			want:   DomainStat{"dup.com": 5},
		},
		{
			name:   "only dots",
			data:   `{"Id":1,"Email":"............"}`,
			domain: "com",
			want:   DomainStat{},
		},
		{
			name:   "dot and @",
			data:   `{"Id":1,"Email":"@......@...com...@"}`,
			domain: "com",
			want:   DomainStat{},
		},
		{
			name:   ".com but invalid",
			data:   `{"Id":1,"Email":"abc@.com@"}`,
			domain: "com",
			want:   DomainStat{},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			t.Parallel()
			r := strings.NewReader(tC.data)
			got, err := GetDomainStat(r, tC.domain)
			if tC.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tC.want, got)
		})
	}
}
