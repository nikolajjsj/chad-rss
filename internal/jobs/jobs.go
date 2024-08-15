package jobs

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"time"

	database "chad-rss/internal/database"
	query "chad-rss/internal/database/sqlc"

	_ "github.com/joho/godotenv/autoload"
	"github.com/mmcdole/gofeed"

	"chad-rss/internal/utils"
)

type Jobs struct {
	db database.Service
}

func RunJobs() {
	NewJobs := &Jobs{
		db: database.New(),
	}

	syncTicker := time.NewTicker(1 * time.Hour)

	for {
		select {
		case <-syncTicker.C:
			log.Println("Jobs: Fetching new articles... Started")
			NewJobs.fetchNewArticles()
			log.Println("Jobs: Fetching new articles... Done")
		}
	}
}

// NOTE: Fetch all new articles for all the feeds
func (s *Jobs) fetchNewArticles() {
	feeds, err := s.db.Query().GetAllFeeds(context.Background())
	if err != nil {
		log.Println("Get all feeds error: ", err)
		return
	}

	ch := make(chan query.GetAllFeedsRow, len(feeds))
	go s.processor(ch)
	for _, feed := range feeds {
		ch <- feed
	}
}

func (s *Jobs) processor(ch chan query.GetAllFeedsRow) {
	for feed := range ch {
		fp := gofeed.NewParser()
		f, err := fp.ParseURL(feed.Url)
		if err != nil {
			log.Println("Parse feed error: ", err)
			continue
		}

		for _, item := range f.Items {
			// TODO: move this to a separate function
			if item.GUID == "" || item.Link == "" || item.Title == "" || item.PublishedParsed == nil {
				continue
			}

			media := ""
			if item.Image != nil {
				media = item.Image.URL
			}

			authors := ""
			if item.Authors != nil {
				stringSlice := make([]string, len(item.Authors))
				for i, v := range item.Authors {
					stringSlice[i] = v.Name
				}
				authors = strings.Join(stringSlice[:], ",")
			}

			nid, err := utils.GenerateNID()
			if err != nil {
				log.Println("Generate NID error: ", err)
				continue
			}

			if _, err = s.db.Query().CreateFeedArticles(context.Background(), query.CreateFeedArticlesParams{
				Nid:         nid,
				RssID:       item.GUID,
				FeedID:      feed.ID,
				Url:         item.Link,
				Title:       item.Title,
				Summary:     sql.NullString{String: item.Description, Valid: item.Description != ""},
				Content:     sql.NullString{String: item.Content, Valid: item.Content != ""},
				Authors:     sql.NullString{String: authors, Valid: authors != ""},
				Media:       sql.NullString{String: media, Valid: media != ""},
				PublishedAt: sql.NullTime{Time: *item.PublishedParsed, Valid: item.PublishedParsed != nil},
			}); err != nil {
				log.Println("Create article in DB error: ", err)
				continue
			}
		}
	}
}
