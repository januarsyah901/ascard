package articles

import "github.com/januarsyah901/ascard/api/internal/httpx"

// Article is the JSON shape of an "Article" row, keys in the same order as the Prisma output.
type Article struct {
	ID        string     `json:"id"`
	Title     string     `json:"title"`
	Content   string     `json:"content"`
	Category  *string    `json:"category"`
	Published bool       `json:"published"`
	CreatedAt httpx.Time `json:"createdAt"`
	UpdatedAt httpx.Time `json:"updatedAt"`
}

// Input is a validated article body for create and update.
type Input struct {
	Title     string
	Content   string
	Category  *string
	Published bool
}
