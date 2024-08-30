package githubreleases

import (
	"context"
	"fmt"
	"io"

	"github.com/google/go-github/v63/github"
)

type FileDownloader struct {
	client *github.Client

	owner  string
	repo   string
	gitref string
}

func NewFileDownloader(client *github.Client, owner, repo, gitRef string) *FileDownloader {
	return &FileDownloader{
		client: client,
		owner:  owner,
		repo:   repo,
		gitref: gitRef,
	}
}

func (d *FileDownloader) DownloadFile(ctx context.Context, file string) ([]byte, error) {
	_, res, err := d.client.Repositories.DownloadContents(ctx, d.owner, d.repo, file, &github.RepositoryContentGetOptions{
		Ref: d.gitref,
	})
	if err != nil {
		return nil, fmt.Errorf("downloading github file (%s/%s/%s @ %s): %w", d.owner, d.repo, file, d.gitref, err)
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}
	return data, nil
}
