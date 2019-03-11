package cmd

import (
	"fmt"
	"log"

	"github.com/nrocco/bookmarks/storage"
	"github.com/spf13/cobra"
)

var testFetchBookmarkCmd = &cobra.Command{
	Use:   "test-bookmark",
	Short: "Test downloading a bookmark",
	RunE: func(cmd *cobra.Command, args []string) error {
		bookmark := storage.Bookmark{
			URL: args[0],
		}

		if err := bookmark.Fetch(); err != nil {
			log.Fatalf("Oh no!")
		}

		fmt.Printf("Title: %s\n", bookmark.Title)
		fmt.Printf("URL: %s\n", bookmark.URL)
		fmt.Printf("Excerpt: %s\n", bookmark.Excerpt)
		fmt.Println()
		fmt.Println(bookmark.Content)

		return nil
	},
}

var testFetchFeedCmd = &cobra.Command{
	Use:   "test-feed",
	Short: "Test downloading a feed",
	RunE: func(cmd *cobra.Command, args []string) error {
		feed := storage.Feed{
			URL: args[0],
			// Refreshed: time.Now().Add(-7 * time.Hour),
			// Etag:      os.Args[2],
		}
		feedItems := []*storage.FeedItem{}

		if err := feed.Fetch(&feedItems); err != nil {
			return err
		}

		fmt.Printf("---\n")
		fmt.Printf("Feed:\n")
		fmt.Printf("  Title: %s\n", feed.Title)
		fmt.Printf("  URL: %s\n", feed.URL)
		fmt.Printf("  Etag: %s\n", feed.Etag)
		fmt.Printf("  Refreshed: %v\n", feed.Refreshed)
		fmt.Printf("  LastAuthored: %v\n", feed.LastAuthored)
		fmt.Printf("  Items:\n")

		for _, item := range feedItems {
			fmt.Printf("  - Title: %s\n", item.Title)
			fmt.Printf("    URL: %s\n", item.URL)
			fmt.Printf("    Date: %v\n", item.Date)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(testFetchBookmarkCmd)
	rootCmd.AddCommand(testFetchFeedCmd)
}
