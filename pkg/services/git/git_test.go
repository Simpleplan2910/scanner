package git

import "testing"

func Test_GetTextFiles(t *testing.T) {
	dir := "samples"
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
	if len(files) == 0 {
		t.Errorf("not found any file")
		return
	}
}
