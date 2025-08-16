package main

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sleklere/gator/internal/database"
	"github.com/sleklere/gator/internal/feed"
)

func scrapeFeeds(s *state) error {
	fmt.Println("scraping feeds")
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}

	err = s.db.MarkFeedFetched(context.Background(), nextFeed.ID)
	if err != nil {
		return err
	}

	feedRes, err := feed.FetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		return err
	}


	for _, item := range feedRes.Channel.Item {
		publishedAt, err := parseRSSDate(item.PubDate)
		if err != nil {
			fmt.Printf("error in item with title '%s': %v\n", item.Title, err)
			fmt.Printf("skipping post save for item with title '%s'\n", item.Title)
			continue
		}
		postParams := database.CreatePostParams{
			ID: uuid.New(),
			Title: item.Title,
			Url: item.Link,
			Description: toNullString(item.Description),
			FeedID: nextFeed.ID,
			PublishedAt: publishedAt,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		_, err = s.db.CreatePost(context.Background(), postParams)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint \"posts_url_key\"") {
				fmt.Println("duplicate url err")
				continue
			}

			fmt.Printf("error in item with title '%s': %v", item.Title, err)
			fmt.Printf("skipping post save for item with title '%s'", item.Title)
			continue
		}
	}

	fmt.Printf("finished scraping feed '%s'\n", nextFeed.Name)
	return nil
}

func parseRSSDate(dateStr string) (time.Time, error) {
	// try common RSS date formats
	formats := []string{
		time.RFC1123Z, // "Mon, 02 Jan 2006 15:04:05 -0700"
		time.RFC1123,  // "Mon, 02 Jan 2006 15:04:05 MST"
		time.RFC3339,  // "2006-01-02T15:04:05Z07:00"
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
}

func toNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}
