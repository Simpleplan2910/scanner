package git

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
)

type Service interface {
	GetRepos(reposURl string) (r Repos, err error)
}

type gitRepos struct {
	TmpDir string
}

func New(dir string) Service {
	return &gitRepos{TmpDir: dir}
}

func (g *gitRepos) GetRepos(reposURl string) (r Repos, err error) {
	l := strings.Split(reposURl, "/")
	name := l[len(l)-1]
	dir := filepath.Join(g.TmpDir, name)
	_, err = os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			_, err = git.PlainClone(dir, false, &git.CloneOptions{
				URL: reposURl,
			})
			if err != nil {
				return nil, err
			}
			return &repos{
				dir: dir,
			}, nil
		}
		return nil, err
	}
	return &repos{
		dir: dir,
	}, nil
}
