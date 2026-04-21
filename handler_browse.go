package main

import (
	"context"
	"flag"
	"fmt"
	"sort"
	"strings"

	"github.com/jasonsoprovich/rss-aggregator/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	fs := flag.NewFlagSet("browse", flag.ContinueOnError)
	limit := fs.Int("limit", 2, "number of posts per page")
	page := fs.Int("page", 1, "page number (1-based)")
	sortBy := fs.String("sort", "date", "sort field: date, title, feed")
	order := fs.String("order", "desc", "sort order: asc, desc")
	feedFilter := fs.String("feed", "", "filter by feed name (case-insensitive substring match)")

	if err := fs.Parse(cmd.Args); err != nil {
		return fmt.Errorf("invalid arguments: %w", err)
	}

	switch *sortBy {
	case "date", "title", "feed":
	default:
		return fmt.Errorf("invalid sort field %q: use date, title, or feed", *sortBy)
	}
	switch *order {
	case "asc", "desc":
	default:
		return fmt.Errorf("invalid order %q: use asc or desc", *order)
	}

	posts, err := s.db.GetAllPostsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("couldn't get posts for user: %w", err)
	}

	if *feedFilter != "" {
		lower := strings.ToLower(*feedFilter)
		filtered := posts[:0]
		for _, p := range posts {
			if strings.Contains(strings.ToLower(p.FeedName), lower) {
				filtered = append(filtered, p)
			}
		}
		posts = filtered
	}

	sort.SliceStable(posts, func(i, j int) bool {
		var less bool
		switch *sortBy {
		case "title":
			less = strings.ToLower(posts[i].Title) < strings.ToLower(posts[j].Title)
		case "feed":
			less = strings.ToLower(posts[i].FeedName) < strings.ToLower(posts[j].FeedName)
		default: // date
			ti := posts[i].PublishedAt.Time
			tj := posts[j].PublishedAt.Time
			less = ti.Before(tj)
		}
		if *order == "desc" {
			return !less
		}
		return less
	})

	if *page < 1 {
		return fmt.Errorf("page must be 1 or greater")
	}
	if *limit < 1 {
		return fmt.Errorf("limit must be 1 or greater")
	}

	total := len(posts)
	totalPages := (total + *limit - 1) / *limit
	if totalPages == 0 {
		totalPages = 1
	}

	offset := (*page - 1) * *limit
	if offset >= total {
		return fmt.Errorf("page %d out of range (total pages: %d)", *page, totalPages)
	}
	end := offset + *limit
	if end > total {
		end = total
	}
	posts = posts[offset:end]

	fmt.Printf("Page %d of %d (%d total posts) for user %s:\n", *page, totalPages, total, user.Name)
	for _, post := range posts {
		fmt.Printf("%s from %s\n", post.PublishedAt.Time.Format("Mon Jan 2"), post.FeedName)
		fmt.Printf("--- %s ---\n", post.Title)
		fmt.Printf("    %v\n", post.Description.String)
		fmt.Printf("Link: %s\n", post.Url)
		fmt.Println("=====================================")
	}
	if *page < totalPages {
		fmt.Printf("(use --page=%d to see the next page)\n", *page+1)
	}

	return nil
}
