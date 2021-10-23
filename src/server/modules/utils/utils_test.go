package utils

import (
	"testing"
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
			name: "Adding one child to SitemapPath",
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.addChild(tt.args.parentID, tt.args.child); got != tt.want {
				t.Errorf("Sitemap.addChild() = %v, want %v", got, tt.want)
			}
		})
	}
}
