package git

import "testing"

func Test_gitRepos(t *testing.T) {
	dir := "."
	git := New(dir)
	r, err := git.GetRepos("https://github.com/guardrailsio/backend-engineer-challenge.git")
	if err != nil {
		t.Errorf("get repos error: %s", err)
		return
	}
	files, err := r.GetTextFiles()
	if err != nil {
		t.Errorf("get text file error: %s", err)
		return
	}
	t.Logf("list files %+v", files)
}
