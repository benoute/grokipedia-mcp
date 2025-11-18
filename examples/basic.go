package main

import (
	"context"
	"fmt"
	"log"

	"github.com/benoute/grokipedia/pkg/grokipedia"
)

func main() {
	ctx := context.Background()
	results, err := grokipedia.Search(ctx, "grok", grokipedia.WithLimit(10))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Total", len(results))
	for _, result := range results {
		fmt.Printf("%+v\n", result)
	}

	page, err := grokipedia.GetPage(ctx, "Grok", grokipedia.WithoutContent())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Page: %+v\n", page)
}
