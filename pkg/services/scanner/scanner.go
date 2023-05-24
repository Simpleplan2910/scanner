package scanner

import (
	"bufio"
	"context"
	"os"
	"scanner/pkg/db"
	"scanner/pkg/services/git"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type Service interface {
	Scan(ctx context.Context, r *Repos, substr string) error
}

type serviceLocal struct {
	resultStore db.ResultStore
	scanStore   db.ScanStore
	git         git.Service
	numWorker   int
}

func NewLocalService(r db.ResultStore, git git.Service) Service {
	return serviceLocal{resultStore: r, git: git}
}

func (s serviceLocal) Scan(ctx context.Context, r *Repos, substr string) error {
	sctx := context.TODO()
	go s.scan(sctx, r, substr)
	return nil
}

func (s serviceLocal) scan(ctx context.Context, rep *Repos, substr string) {
	// get repo to local disk
	// bc this local, just read 1 file and send to other worker
	// buffer 1 or more files for the next read
	// TODO: move to search by line per worker instead of by file
	r, err := s.git.GetRepos(rep.ReposURL)
	if err != nil {
		logrus.Errorf("Error when get git repos %s", err)
		if err = s.scanStore.UpdateError(ctx, rep.ScanId); err != nil {
			logrus.Errorf("Error when update scan failed to storage with error %s", err)
		}
		return
	}
	files, err := r.GetTextFiles()
	if err != nil {
		logrus.Errorf("Error when get text files in git repos %s", err)
		if err = s.scanStore.UpdateError(ctx, rep.ScanId); err != nil {
			logrus.Errorf("Error when update scan failed to storage with error %s", err)
		}
		return
	}

	fls := make(chan file)
	results := make(chan *scanResult)

	for i := 0; i < s.numWorker; i++ {
		go worker(fls, results, substr)
	}

	for _, f := range files {
		b, err := os.ReadFile(f)
		if err != nil {
			return
		}
		fls <- file{
			File:     b,
			FileName: f,
		}
	}

	for r := range results {
		er := ""
		if r.Error != nil {
			er = r.Error.Error()
		}
		res := &db.Result{
			ScanID:    rep.ScanId,
			Lines:     r.Lines,
			Filename:  r.Filename,
			Error:     er,
			CreatedAt: time.Now(),
		}
		_, err = s.resultStore.Add(ctx, res)
		if err != nil {
			logrus.Errorf("Failed to add restult to storage with error %s", err)
			if err = s.scanStore.UpdateError(ctx, rep.ScanId); err != nil {
				logrus.Errorf("Error when update scan failed tto storage with error %s", err)
			}
			return
		}

	}
}

// TODO: add timeout
func worker(files <-chan file, results chan<- *scanResult, substr string) {
	for f := range files {
		lines, err := findSubstring(f.File, substr)
		results <- &scanResult{
			Lines:    lines,
			Error:    err,
			Filename: f.FileName,
		}
	}
}

func findSubstring(file []byte, substr string) (lines []db.Line, err error) {
	// maxsize token 64 * 1024
	bs := bufio.NewScanner(strings.NewReader(string(file)))
	var lineNum uint32 = 0
	for bs.Scan() {
		line := db.Line{
			LineNum: lineNum,
		}
		ind := 0
		// get next line from file
		text := bs.Text()
		substrInd := strings.Index(text[ind:], substr)
		for substrInd != -1 {
			line.Indexes = append(line.Indexes, substrInd)
			ind = substrInd
			substrInd = strings.Index(text[ind+len(substr):], substr)
			// not found
			if substrInd == -1 {
				break
			}
			substrInd = substrInd + ind + len(substr)
		}
		if len(line.Indexes) != 0 {
			lines = append(lines, line)
		}
		lineNum += 1
	}
	if err := bs.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}
