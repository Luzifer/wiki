package main

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

const yamlDelimiter = `---`

var errFileNotFound = errors.New("Specified file was not found")

type storedFile struct {
	Meta    map[string]any `json:"meta"`
	Content string         `json:"content"`

	AuthorName  string `json:"-"`
	AuthorEmail string `json:"-"`
}

func loadStoredFile(filename string) (*storedFile, error) {
	if _, err := os.Stat(path.Join(cfg.DataDir, filename)); err != nil {
		return nil, errFileNotFound
	}

	content, err := os.ReadFile(path.Join(cfg.DataDir, filename)) //#nosec:G304 // Reading data file
	if err != nil {
		return nil, errors.Wrap(err, "Unable to read file")
	}

	return storedFileFromString(string(content))
}

func storedFileFromString(content string) (*storedFile, error) {
	// Look at first line and see whether this file has a metadata part
	lines := strings.Split(strings.TrimSpace(content), "\n")
	if len(lines) == 0 {
		// Empty file
		return &storedFile{}, nil
	}

	var (
		metadata     []string
		contentStart int
	)

	if lines[0] == yamlDelimiter {
		// This file has a metadata part
		for i := 1; i < len(lines); i++ {
			if lines[i] == yamlDelimiter {
				contentStart = i + 1
				break
			}

			metadata = append(metadata, lines[i])
		}
	}

	file := &storedFile{
		Content: strings.TrimSpace(strings.Join(lines[contentStart:], "\n")),
		Meta:    map[string]interface{}{},
	}

	if len(metadata) > 0 {
		if err := yaml.NewDecoder(strings.NewReader(strings.Join(metadata, "\n"))).Decode(&file.Meta); err != nil {
			return nil, errors.Wrap(err, "Unable to parse metadata part")
		}
	}

	return file, nil
}

func (s storedFile) GetMetaString(key string) string {
	if v, ok := s.Meta[key].(string); ok {
		return v
	}
	return ""
}

func (s storedFile) Save(filename string) error {
	repo, err := git.PlainOpen(cfg.DataDir)
	if err != nil {
		if err != git.ErrRepositoryNotExists {
			return errors.Wrap(err, "Unable to open repository")
		}

		if repo, err = s.initRepo(); err != nil {
			return errors.Wrap(err, "Unable to init repo")
		}
	}

	wt, err := repo.Worktree()
	if err != nil {
		return errors.Wrap(err, "Unable to get worktree")
	}

	f, err := os.Create(path.Join(cfg.DataDir, filename)) //#nosec:G304 // Creating data file
	if err != nil {
		return errors.Wrap(err, "Unable to create file")
	}
	defer func() {
		if err := f.Close(); err != nil {
			logrus.WithError(err).Error("closing data file (leaked fd)")
		}
	}()

	if len(s.Meta) > 0 {
		fmt.Fprintln(f, yamlDelimiter)
		if err := yaml.NewEncoder(f).Encode(s.Meta); err != nil {
			return errors.Wrap(err, "Unable to write metadata")
		}
		fmt.Fprintln(f, yamlDelimiter)
	}

	fmt.Fprintln(f, s.Content)

	if _, err := wt.Add(filename); err != nil {
		return errors.Wrap(err, "Unable to add file to index")
	}

	_, err = wt.Commit("Web-Update of "+filename, &git.CommitOptions{Author: s.authorSignature(), Committer: s.committerSignature()})
	return errors.Wrap(err, "Unable to commit file change")
}

func (s storedFile) authorSignature() *object.Signature {
	sig := &object.Signature{Name: "Web-User", Email: "wiki+author@luzifer.io", When: time.Now()}

	if s.AuthorName != "" {
		sig.Name = s.AuthorName
	}

	if s.AuthorEmail != "" {
		sig.Email = s.AuthorEmail
	}

	return sig
}

func (storedFile) committerSignature() *object.Signature {
	return &object.Signature{Name: "wiki " + version, Email: "wiki+committer@luzifer.io", When: time.Now()}
}

func (s storedFile) initRepo() (*git.Repository, error) {
	repo, err := git.PlainInit(cfg.DataDir, false)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to initialize repo")
	}

	if _, err := repo.Branch("master"); err == git.ErrBranchNotFound {
		if err := repo.CreateBranch(&config.Branch{Name: "master"}); err != nil {
			return nil, errors.Wrap(err, "Unable to create master branch")
		}
	}

	wt, err := repo.Worktree()
	if err != nil {
		return nil, errors.Wrap(err, "Unable to get worktree")
	}

	_, err = wt.Commit("Initial commit", &git.CommitOptions{Author: s.authorSignature(), Committer: s.committerSignature()})

	return repo, errors.Wrap(err, "Unable to create initial commit")
}
