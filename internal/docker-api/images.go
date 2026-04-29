package dockerapi

import (
	"context"
	"strings"

	"github.com/moby/moby/client"
)

type Image struct {
	ID         string
	Repository string
	Tag        string
	Size       int64
	Created    int64
}

func FetchImages() []Image {
	var results []Image
	cli := GetClient()
	if cli == nil {
		return results
	}

	ctx := context.Background()
	images, err := cli.ImageList(ctx, client.ImageListOptions{All: true})
	if err != nil {
		return results
	}

	for _, img := range images.Items {
		repo := "<none>"
		tag := "<none>"

		if len(img.RepoTags) > 0 {
			// A tag is usually in the format repository:tag
			// So we split from the last colon.
			lastColonIndex := strings.LastIndex(img.RepoTags[0], ":")
			if lastColonIndex != -1 && lastColonIndex != 0 {
				repo = img.RepoTags[0][:lastColonIndex]
				tag = img.RepoTags[0][lastColonIndex+1:]
			} else {
				repo = img.RepoTags[0]
			}
		} else if len(img.RepoDigests) > 0 {
			// Some images without tags but with repo digests
			repo = strings.Split(img.RepoDigests[0], "@")[0]
		}

		// Clean up ID (remove sha256:)
		id := strings.TrimPrefix(img.ID, "sha256:")
		if len(id) > 12 {
			id = id[:12]
		}

		results = append(results, Image{
			ID:         id,
			Repository: repo,
			Tag:        tag,
			Size:       img.Size,
			Created:    img.Created,
		})
	}

	return results
}
