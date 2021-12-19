package utils

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
)

func TestSitemapPath_addChild(t *testing.T) {
	type args struct {
		parentID int64
		child    *SitemapPath
	}

	tests := []struct {
		name string
		p    *SitemapPath
		args args
		want bool
	}{
		{
			name: "Adding one child to SitemapPath child",
			p: &SitemapPath{
				ID:       1,
				Title:    "Some title",
				Path:     "/some/path",
				Priority: 1,
				Children: []SitemapPath{
					{
						ID:       2,
						ParentID: 1,
						Title:    "First child title",
						Path:     "/child1/path",
						Priority: 0.8,
					},
				},
			},
			args: args{
				parentID: 2,
				child: &SitemapPath{
					ID:       3,
					ParentID: 2,
					Title:    "Second child title",
					Path:     "/child2/path",
					Priority: 0.8,
				},
			},
			want: true,
		},
		{
			name: "Adding one child to SitemapPath child's child",
			p: &SitemapPath{
				ID:       1,
				Title:    "Some title",
				Path:     "/some/path",
				Priority: 1,
				Children: []SitemapPath{
					{
						ID:       2,
						ParentID: 1,
						Title:    "First child title",
						Path:     "/child1/path",
						Priority: 0.8,
						Children: []SitemapPath{
							{
								ID:       3,
								ParentID: 2,
								Title:    "Second child title",
								Path:     "/child2/path",
								Priority: 0.6,
							},
						},
					},
				},
			},
			args: args{
				parentID: 3,
				child: &SitemapPath{
					ID:       4,
					ParentID: 3,
					Title:    "Third child title",
					Path:     "/child3/path",
					Priority: 0.4,
				},
			},
			want: true,
		},
		{
			name: "Error on adding one child to SitemapPath child",
			p: &SitemapPath{
				ID:       1,
				Title:    "Some title",
				Path:     "/some/path",
				Priority: 1,
				Children: []SitemapPath{
					{
						ID:       2,
						ParentID: 1,
						Title:    "First child title",
						Path:     "/child1/path",
						Priority: 0.8,
					},
				},
			},
			args: args{
				parentID: 3,
				child: &SitemapPath{
					ID:       3,
					ParentID: 3,
					Title:    "Second child title",
					Path:     "/child2/path",
					Priority: 0.8,
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.addChild(tt.args.parentID, tt.args.child); got != tt.want {
				t.Errorf("SitemapPath.addChild() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSitemap_addChild(t *testing.T) {
	type args struct {
		parentID int64
		child    *SitemapPath
	}

	tests := []struct {
		name string
		p    *Sitemap
		args args
		want bool
	}{
		{
			name: "Adding one child to Sitemap top-level",
			p: &Sitemap{
				{
					ID:       1,
					Title:    "Some title 1",
					Path:     "/some/path1",
					Priority: 1,
				},
				{
					ID:       2,
					Title:    "Some title 2",
					Path:     "/some/path2",
					Priority: 1,
				},
			},
			args: args{
				parentID: 2,
				child: &SitemapPath{
					ID:       3,
					ParentID: 2,
					Title:    "Second child title",
					Path:     "/child/path",
					Priority: 0.8,
				},
			},
			want: true,
		},
		{
			name: "Adding one child to Sitemap child",
			p: &Sitemap{
				{
					ID:       1,
					Title:    "Some title 1",
					Path:     "/some/path1",
					Priority: 1,
					Children: []SitemapPath{
						{
							ID:       3,
							ParentID: 1,
							Title:    "First child title",
							Path:     "/child1/path",
							Priority: 0.8,
						},
					},
				},
				{
					ID:       2,
					Title:    "Some title 2",
					Path:     "/some/path2",
					Priority: 1,
				},
			},
			args: args{
				parentID: 3,
				child: &SitemapPath{
					ID:       4,
					ParentID: 3,
					Title:    "Second child title",
					Path:     "/child2/path",
					Priority: 0.8,
				},
			},
			want: true,
		},
		{
			name: "Error on adding one child to Sitemap top-level",
			p: &Sitemap{
				{
					ID:       1,
					Title:    "Some title 1",
					Path:     "/some/path1",
					Priority: 1,
				},
				{
					ID:       2,
					Title:    "Some title 2",
					Path:     "/some/path2",
					Priority: 1,
				},
			},
			args: args{
				parentID: 4,
				child: &SitemapPath{
					ID:       3,
					ParentID: 4,
					Title:    "Second child title",
					Path:     "/child/path",
					Priority: 0.8,
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.addChild(tt.args.parentID, tt.args.child); got != tt.want {
				t.Errorf("Sitemap.addChild() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSitemap_removePath(t *testing.T) {
	type args struct {
		index int
	}

	tests := []struct {
		name string
		p    *Sitemap
		args args
		want *Sitemap
	}{
		{
			name: "Remove one path from Sitemap",
			p: &Sitemap{
				{
					ID:       1,
					Title:    "Some title 1",
					Path:     "/some/path1",
					Priority: 1,
				},
				{
					ID:       2,
					Title:    "Some title 2",
					Path:     "/some/path2",
					Priority: 1,
				},
			},
			args: args{index: 1},
			want: &Sitemap{
				{
					ID:       1,
					Title:    "Some title 1",
					Path:     "/some/path1",
					Priority: 1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.removePath(tt.args.index); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Sitemap.removePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSitemap_Nest(t *testing.T) {
	tests := []struct {
		name string
		p    *Sitemap
		want *Sitemap
	}{
		{
			name: "Nest the Sitemap",
			p: &Sitemap{
				{
					ID:       1,
					Title:    "Some title 1",
					Path:     "/some/path1",
					Priority: 1,
				},
				{
					ID:       2,
					ParentID: 1,
					Title:    "Some title 2",
					Path:     "/some/path2",
					Priority: 1,
				},
			},
			want: &Sitemap{
				{
					ID:       1,
					Title:    "Some title 1",
					Path:     "/some/path1",
					Priority: 1,
					Children: []SitemapPath{
						{
							ID:       2,
							ParentID: 1,
							Title:    "Some title 2",
							Path:     "/some/path2",
							Priority: 1,
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.Nest(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Sitemap.Nest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_childLookup(t *testing.T) {
	type args struct {
		item *SitemapPath
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "HTMLize the SitemapPath",
			args: args{
				item: &SitemapPath{
					ID:       1,
					Title:    "Some title 1",
					Path:     "/some/path1",
					Priority: 1,
					Children: []SitemapPath{
						{
							ID:       2,
							ParentID: 1,
							Title:    "Some title 2",
							Path:     "/some/path2",
							Priority: 0.9,
						},
					},
				},
			},
			want: "\n<li><a href=\"/some/path1\" class=\"spa\">Some title 1</a><ul>" +
				"\n<li><a href=\"/some/path2\" class=\"spa\">Some title 2</a></li></ul></li>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := childLookup(tt.args.item); got != tt.want {
				t.Errorf("childLookup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSitemap_ToHTMLString(t *testing.T) {
	tests := []struct {
		name string
		p    *Sitemap
		want string
	}{
		{
			name: "HTMLize the Sitemap",
			p: &Sitemap{
				{
					ID:       1,
					Title:    "Some title 1",
					Path:     "/some/path1",
					Priority: 1,
					Children: []SitemapPath{
						{
							ID:       3,
							ParentID: 1,
							Title:    "Some title 3",
							Path:     "/some/path3",
							Priority: 0.9,
							Children: []SitemapPath{
								{
									ID:       4,
									ParentID: 3,
									Title:    "Some title 4",
									Path:     "/some/path4",
									Priority: 0.8,
								},
							},
						},
					},
				},
				{
					ID:       2,
					Title:    "Some title 2",
					Path:     "/some/path2",
					Priority: 1,
				},
			},
			want: string("<ul>" +
				"\n<li><a href=\"/some/path1\" class=\"spa\">Some title 1</a><ul>" +
				"\n<li><a href=\"/some/path3\" class=\"spa\">Some title 3</a><ul>" +
				"\n<li><a href=\"/some/path4\" class=\"spa\">Some title 4</a></li></ul></li></ul></li>" +
				"\n<li><a href=\"/some/path2\" class=\"spa\">Some title 2</a></li></ul>"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.ToHTMLString(); *got != tt.want {
				t.Errorf("Sitemap.ToHTMLString() = %v, want %v", *got, tt.want)
			}
		})
	}
}

func TestURLParams_ToString(t *testing.T) {
	type args struct {
		escaped bool
	}

	tests := []struct {
		name string
		p    *URLParams
		args args
		want string
	}{
		{
			name: "Stringify empty URLParams",
			p:    &URLParams{},
			args: args{escaped: false},
			want: "",
		},
		{
			name: "Stringify URLParams without escaping",
			p: &URLParams{
				"param1": 1,
				"param2": "some param",
			},
			args: args{escaped: false},
			want: "?param1=1&param2=some param",
		},
		{
			name: "Stringify URLParams with escaping",
			p: &URLParams{
				"param1": 10,
				"param2": "some param @",
			},
			args: args{escaped: true},
			want: "?param1=10&param2=some+param+%40",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.ToString(tt.args.escaped); got != tt.want {
				t.Errorf("URLParams.ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAuthLinks(t *testing.T) {
	tests := []struct {
		name string
		want fiber.Map
	}{
		{
			name: "Test GetAuthLinks",
			want: fiber.Map{
				"vk": "https://oauth.vk.com/authorize" +
					"?client_id=&redirect_uri=https%3A%2F%2F%2Fauth%2F&response_type=code&scope=email&state=vk",
				"fb": "https://www.facebook.com/v11.0/dialog/oauth" +
					"?client_id=&redirect_uri=https%3A%2F%2F%2Fauth%2F&response_type=code&scope=email&state=facebook",
				"ya": "https://oauth.yandex.ru/authorize" +
					"?client_id=&redirect_uri=https%3A%2F%2F%2Fauth%2F&response_type=code&state=yandex",
				"gl": "https://accounts.google.com/o/oauth2/v2/auth" +
					"?access_type=online&client_id=&include_granted_scopes=false&redirect_uri=https%3A%2F%2F%2Fauth%2F&response_type=code" +
					"&scope=https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fuserinfo.profile+https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fuserinfo.email&state=google",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetAuthLinks(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAuthLinks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_IsEmptyStruct(t *testing.T) {
	type args struct {
		object interface{}
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Test nil to true",
			args: args{
				object: nil,
			},
			want: true,
		},
		{
			name: "Test empty bool to true",
			args: args{
				object: false,
			},
			want: true,
		},
		{
			name: "Test empty string to true",
			args: args{
				object: "",
			},
			want: true,
		},
		{
			name: "Test pointer to empty string to true",
			args: args{
				object: new(string),
			},
			want: true,
		},
		{
			name: "Test non-empty string to false",
			args: args{
				object: "Some",
			},
			want: false,
		},
		{
			name: "Test pointer to non-empty slice of int to false",
			args: args{
				object: &[1]int{1},
			},
			want: false,
		},
		{
			name: "Test empty Sitemap struct to true",
			args: args{
				object: Sitemap{},
			},
			want: true,
		},
		{
			name: "Test pointer to empty Sitemap struct to true",
			args: args{
				object: &Sitemap{},
			},
			want: true,
		},
		{
			name: "Test non-empty Sitemap struct to false",
			args: args{
				object: Sitemap{
					{
						ID:       1,
						Title:    "Some title",
						Path:     "/some/path",
						Priority: 1,
					},
				},
			},
			want: false,
		},
		{
			name: "Test pointer to non-empty Sitemap struct to false",
			args: args{
				object: &Sitemap{
					{
						ID:       1,
						Title:    "Some title",
						Path:     "/some/path",
						Priority: 1,
					},
				},
			},
			want: false,
		},
		{
			name: "Test empty SitemapPath struct to true",
			args: args{
				object: SitemapPath{},
			},
			want: true,
		},
		{
			name: "Test pointer to empty SitemapPath struct to true",
			args: args{
				object: &SitemapPath{},
			},
			want: true,
		},
		{
			name: "Test non-empty SitemapPath struct to false",
			args: args{
				object: SitemapPath{
					ID:       1,
					Title:    "Some title",
					Path:     "/some/path",
					Priority: 1,
				},
			},
			want: false,
		},
		{
			name: "Test pointer to non-empty SitemapPath struct to false",
			args: args{
				object: &SitemapPath{
					ID:       1,
					Title:    "Some title",
					Path:     "/some/path",
					Priority: 1,
				},
			},
			want: false,
		},
		{
			name: "Test empty URLParams map to true",
			args: args{
				object: URLParams{},
			},
			want: true,
		},
		{
			name: "Test pointer to empty URLParams map to true",
			args: args{
				object: &URLParams{},
			},
			want: true,
		},
		{
			name: "Test non-empty URLParams map to false",
			args: args{
				object: URLParams{
					"param1": 0,
				},
			},
			want: false,
		},
		{
			name: "Test pointer to non-empty URLParams map to false",
			args: args{
				object: &URLParams{
					"param1": 0,
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsEmptyStruct(tt.args.object); got != tt.want {
				t.Errorf("IsEmptyStruct() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIterateStruct(t *testing.T) {
	type AX []interface{}

	type struct1 struct {
		Id  *int
		Str *string
	}

	type struct2 struct {
		Id  *int
		Str *string
		Rts *int
	}

	type struct3 struct {
		Id       int
		StrFirst string
		Rts      *int
	}

	id := 5
	str := "asd"

	type args struct {
		s interface{}
	}

	tests := []struct {
		name    string
		args    args
		wantSet string
	}{
		{
			name: "Test for struct{*int, *string}",
			args: args{
				s: struct1{
					Id:  &id,
					Str: &str,
				},
			},
			wantSet: "[[= field_id 5] [= field_str asd]]",
		},
		{
			name: "Test for pointer to struct{*int, *string}",
			args: args{
				s: &struct1{
					Id:  &id,
					Str: &str,
				},
			},
			wantSet: "[[= field_id 5] [= field_str asd]]",
		},
		{
			name: "Test for pointer to struct{*int, *string, nil}",
			args: args{
				s: &struct2{
					Id:  &id,
					Str: &str,
				},
			},
			wantSet: "[[= field_id 5] [= field_str asd]]",
		},
		{
			name: "Test for pointer to struct{*int, *string, nil} with big field names",
			args: args{
				s: &struct3{
					Id:       id,
					StrFirst: str,
				},
			},
			wantSet: "[[= field_id 5] [= field_str_first asd]]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set := make([]interface{}, 0)
			re := regexp.MustCompile(`[A-Z][a-z0-9]*`)

			IterateStruct(tt.args.s, func(field string, value interface{}) {
				f := strings.Join(re.FindAllString(field, -1), "_")
				set = append(set, AX{"=", "field_" + strings.ToLower(f), value})
			})

			if fmt.Sprint(set) != tt.wantSet {
				t.Errorf("IterateStruct() = %v, wantSet %v", fmt.Sprint(set), tt.wantSet)
			}
		})
	}
}

func TestContains(t *testing.T) {
	type args struct {
		val   uint64
		slice []uint64
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "[1, 2, 3] contains 1",
			args: args{
				val:   1,
				slice: []uint64{1, 2, 3},
			},
			want: true,
		},
		{
			name: "[1, 2, 3] not contains 4",
			args: args{
				val:   4,
				slice: []uint64{1, 2, 3},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Contains(tt.args.val, tt.args.slice...); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimeMSK_ToTime(t *testing.T) {
	tNow := time.Now()

	loc, err := time.LoadLocation("UTC")
	if err != nil {
		panic(err)
	}

	type args struct {
		mock []time.Time
	}

	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{
			name: "Get mockable time.Now",
			args: args{
				mock: []time.Time{tNow},
			},
			want: tNow.In(loc).Add(time.Hour * time.Duration(3)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TimeMSK_ToTime(tt.args.mock...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TimeMSK_ToTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimeMSK_ToString(t *testing.T) {
	tNow := time.Now()

	loc, err := time.LoadLocation("UTC")
	if err != nil {
		panic(err)
	}

	type args struct {
		mock []time.Time
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Get string of mockable TimeMSK",
			args: args{
				mock: []time.Time{tNow},
			},
			want: tNow.In(loc).Add(time.Hour * time.Duration(3)).Format("2006-01-02 15:04:05"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TimeMSK_ToString(tt.args.mock...); got != tt.want {
				t.Errorf("TimeMSK_ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimeMSK_ToLocaleString(t *testing.T) {
	tNow := time.Now()

	loc, err := time.LoadLocation("UTC")
	if err != nil {
		panic(err)
	}

	type args struct {
		mock []time.Time
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Get locale string of mockable TimeMSK",
			args: args{
				mock: []time.Time{tNow},
			},
			want: tNow.In(loc).Add(time.Hour * time.Duration(3)).Format("02.01.2006 15:04:05"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TimeMSK_ToLocaleString(tt.args.mock...); got != tt.want {
				t.Errorf("TimeMSK_ToLocaleString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDateMSK_ToLocaleString(t *testing.T) {
	tNow := time.Now()

	loc, err := time.LoadLocation("UTC")
	if err != nil {
		panic(err)
	}

	type args struct {
		mock []time.Time
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Get locale string of mockable DateMSK",
			args: args{
				mock: []time.Time{tNow},
			},
			want: tNow.In(loc).Add(time.Hour * time.Duration(3)).Format("02.01.2006"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DateMSK_ToLocaleString(tt.args.mock...); got != tt.want {
				t.Errorf("DateMSK_ToLocaleString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDateMSK_ToLocaleSepString(t *testing.T) {
	tNow := time.Now()

	loc, err := time.LoadLocation("UTC")
	if err != nil {
		panic(err)
	}

	type args struct {
		mock []time.Time
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Get locale sep string of mockable DateMSK",
			args: args{
				mock: []time.Time{tNow},
			},
			want: tNow.In(loc).Add(time.Hour * time.Duration(3)).Format("02-01-2006"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DateMSK_ToLocaleSepString(tt.args.mock...); got != tt.want {
				t.Errorf("DateMSK_ToLocaleSepString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHashSHA512(t *testing.T) {
	type args struct {
		str string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test HashSHA512",
			args: args{
				str: "hash",
			},
			want: "30163935c002fc4e1200906c3d30a9c4956b4af9f6dcaef1eb4b1fcb8fba69e7a7acdc491ea5b1f2864ea8c01b01580ef09defc3b11b3f183cb21d236f7f1a6b",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HashSHA512(tt.args.str); got != tt.want {
				t.Errorf("HashSHA512() = %v, want %v", got, tt.want)
			}
		})
	}
}
