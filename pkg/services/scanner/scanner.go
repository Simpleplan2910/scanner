package scanner

import (
	"bufio"
	"context"
	"os"
	"scanner/pkg/db"
	"scanner/pkg/services/git"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

//go:generate mockery --name Service
type Service interface {
	Scan(ctx context.Context, r *Repos, substr string) error
}

type serviceLocal struct {
	resultStore db.ResultStore
	scanStore   db.ScanStore
	git         git.Service
	numWorker   int
}

func NewLocalService(r db.ResultStore, git git.Service, gs db.ScanStore, nWoker int) Service {
	return serviceLocal{resultStore: r, git: git, scanStore: gs, numWorker: nWoker}
}

func (s serviceLocal) Scan(ctx context.Context, r *Repos, substr string) error {
	sctx := context.TODO()
	go s.scan(sctx, r, substr)
	return nil
}

func (s serviceLocal) scan(ctx context.Context, rep *Repos, substr string) {
	// TODO: move to search by line per worker instead of by file
	r, err := s.git.GetRepos(rep.ReposURL)
	if err != nil {
		logrus.Errorf("Error when get git repos %s", err)
		if err = s.scanStore.UpdateError(ctx, rep.ScanId); err != nil {
			logrus.Errorf("Error when update scan failed to storage with error %s", err)
		}
		return
	}
	defer r.Clean()
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
	var wg sync.WaitGroup
	wg.Add(len(files))

	for i := 0; i < s.numWorker; i++ {
		go worker(fls, results, &wg, substr)
	}

	go func() {
		for _, f := range files {
			b, err := os.ReadFile(f)
			if err != nil {
				wg.Done()
				results <- &scanResult{
					Error:    err,
					Filename: f,
				}
				continue
			}
			fls <- file{
				File:     b,
				FileName: f,
			}
		}
	}()
	go func() {
		wg.Wait()
		close(results)
	}()
	for r := range results {
		// empty result, dont add to db
		if err == nil && len(r.Lines) == 0 {
			continue
		}
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
				logrus.Errorf("Error when update scan failed to storage with error %s", err)
			}
			return
		}
	}
	logrus.Infoln("----!DONE!------")
}

// TODO: add timeout
func worker(files <-chan file, results chan<- *scanResult, wg *sync.WaitGroup, substr string) {
	for f := range files {
		lines, err := findSubstring(f.File, substr)
		results <- &scanResult{
			Lines:    lines,
			Error:    err,
			Filename: f.FileName,
		}
		wg.Done()
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
