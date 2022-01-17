package queuejob

import (
	"bufio"
	"context"
	"io"
	"strings"
	"sync"
	"test_guard/pkg/db"
	"test_guard/pkg/services/git"
	"time"
)

type QueueJob interface {
	AddJob(ctx context.Context, j *Job) error
	Start()
}

func New() QueueJob {
	return &queueJob{}
}

// we have some choices with the implementation
// do multiple repos at the same time across a pool of workers or smaller services
// or do one repos at a time but multiple files process across multiple cores
// I chose the later bc we run on single machine
type queueJob struct {
	nWorker     int
	git         git.Service
	queue       []repos
	done        chan struct{}
	resultStore db.ResultStore
	sync.Mutex
}

func NewQueueJob(git git.Service, resultStore db.ResultStore, nWorker int) QueueJob {
	return &queueJob{
		nWorker:     nWorker,
		git:         git,
		resultStore: resultStore,
	}
}

func (q *queueJob) AddJob(ctx context.Context, j *Job) error {
	r, err := q.git.GetRepos(j.ReposURL)
	if err != nil {
		return err
	}
	repos := repos{
		Job:   j,
		Repos: r,
	}
	q.Lock()
	result := &db.Result{
		ReposId:        repos.Job.ReposId,
		Status:         db.Queued,
		RepositoryUrl:  j.ReposURL,
		RepositoryName: j.ReposName,
		QueuedAt:       time.Now(),
	}
	//TODO: use transaction instead of in app lock
	id, err := q.resultStore.Add(ctx, result)
	if err != nil {
		return err
	}
	repos.ResultId = id
	q.queue = append(q.queue, repos)
	q.Unlock()

	return nil
}

func (q *queueJob) Start() {
	q.queue = make([]repos, 0)
	go q.waitJob()
}

func (q *queueJob) Stop() {
	close(q.done)
}

func (q *queueJob) scanFile(bChan <-chan singleFile, resultChan chan<- scanResult) {
	for b := range bChan {
		select {
		case resultChan <- q.scan(b):
		case <-q.done:
			return
		}
	}
}

func (q *queueJob) scan(f singleFile) (r scanResult) {
	br := bufio.NewReader(f.Reader)
	r = scanResult{
		Filename: f.FileName,
	}
	nLine := 1
	for {
		// TODO: limit length
		line, err := br.ReadString('\n')
		if err != nil {
			// end of file, it return the last line too
			if err == io.EOF {
				if isContainSecretKey(line) {
					r.IsContainVulnerable = true
					r.Line = append(r.Line, nLine)
					return r
				}
				return r
			}
			r.Error = err
			return r
		}
		if isContainSecretKey(line) {
			r.IsContainVulnerable = true
			r.Line = append(r.Line, nLine)
		}
		nLine += 1
	}
}

func isContainSecretKey(str string) bool {
	// space before the keyword
	if strings.Contains(str, " public_key") || strings.Contains(str, " private_key") {
		return true
	}
	// line start with keyword
	if len(str) > len("public_key") && str[:len("public_key")] == "public_key" {
		return true
	}
	// line start with keyword
	if len(str) > len("private_key") && str[:len("private_key")] == "private_key" {
		return true
	}
	return false
}

func (q *queueJob) waitJob() {
	for {
		select {
		case <-q.done:
			return
		default:
			var repos *repos
			q.Lock()
			if len(q.queue) > 0 {
				repos = &q.queue[0]
				q.queue = q.queue[1:]
			}
			q.Unlock()
			if repos != nil {
				ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
				err := q.resultStore.UpdateScanningAt(ctx, repos.ResultId, time.Now())
				if err != nil {
					// update failed bc of db, reschedule it or sth
					continue
				}
				err = q.resultStore.UpdateStatus(ctx, repos.ResultId, db.InProgress)
				if err != nil {
					// update failed bc of db, reschedule it or sth
					continue

				}

				results, err := q.doJob(repos)
				if err != nil {
					ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
					err = q.resultStore.UpdateStatus(ctx, repos.ResultId, db.Failure)
					if err != nil {
						// update failed bc of db, reschedule it or sth
						continue

					}
					continue
				}
				ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
				err = q.resultStore.UpdateStatus(ctx, repos.ResultId, db.Success)
				if err != nil {
					// update failed bc of db, reschedule it or sth
					continue
				}
				err = q.resultStore.UpdateStatus(ctx, repos.ResultId, db.Success)
				if err != nil {
					// update failed bc of db, reschedule it or sth
					continue
				}
				resultJson, err := toJsonb(results)
				if err != nil {
					// reschedule it or sth
					continue
				}
				err = q.resultStore.UpdateFinding(ctx, repos.ResultId, resultJson)
				if err != nil {
					// reschedule it or sth
					continue
				}
			}
			time.Sleep(300 * time.Millisecond)
		}

	}
}

func toJsonb(r []scanResult) (s string, err error) {
	return "test", nil
}

func (q *queueJob) doJob(r *repos) (results []scanResult, err error) {
	bChan := make(chan singleFile)
	resultChan := make(chan scanResult)
	results = []scanResult{}
	files, err := r.Repos.GetTextFiles()
	if err != nil {
		return nil, err
	}

	for i := 0; i < q.nWorker; i++ {
		go q.scanFile(bChan, resultChan)
	}

	go func() {
		for _, v := range files {
			b, err := r.Repos.ReadFile(v)
			if err != nil {
				resultChan <- scanResult{
					Error:    err,
					Filename: v,
				}
				continue
			}
			select {
			case <-q.done:
				return
			case bChan <- singleFile{Reader: b, FileName: v}:
			}
		}
		close(resultChan)
	}()
	for result := range resultChan {
		results = append(results, result)
	}

	close(bChan)
	return results, nil
}
