package file

import (
	"github.com/digital-monster-1997/digicore/pkg/dlog"
	//"github.com/digital-monster-1997/digicore/pkg/dlog"
	"github.com/digital-monster-1997/digicore/pkg/utils/dfile"
	"github.com/digital-monster-1997/digicore/pkg/utils/dgo"
	"github.com/fsnotify/fsnotify"
	"io/ioutil"
	"log"
	"path/filepath"
)

type fileDataSourceProvider struct {
	path 			string
	dir 			string
	enableWatch 	bool
	changed 		chan struct{}
}

func NewDataSource(path string, watch bool) *fileDataSourceProvider {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		//xlog.Panic("new datasource", xlog.Any("err", err))
		panic("new datasource")
	}
	dir := dfile.CheckAndGetParentDir(absolutePath)
	ds := &fileDataSourceProvider{path: absolutePath, dir: dir, enableWatch: watch}
	if watch {
		ds.changed = make(chan struct{}, 1)
		go dgo.RecoverGo(ds.watch, nil)
	}
	return ds
}

// ReadConfig ...
func (fp *fileDataSourceProvider) ReadConfig() (content []byte, err error) {
	return ioutil.ReadFile(fp.path)
}

// Close ...
func (fp *fileDataSourceProvider) Close() error {
	close(fp.changed)
	return nil
}

// IsConfigChanged ...
func (fp *fileDataSourceProvider) IsConfigChanged() <-chan struct{} {
	return fp.changed
}


// Watch file and automate update.
func (fp *fileDataSourceProvider) watch() {
	// fsnotify 監看檔案變動
	w, err := fsnotify.NewWatcher()
	if err != nil {
		//log.Printf("new file watcher")
		dlog.Fatal("new file watcher", dlog.FieldMod("file datasource"), dlog.Any("err", err))
	}

	defer w.Close()
	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-w.Events:
				//log.Printf("read watch even, file datasource,%s, %s", filepath.Clean(event.Name),filepath.Clean(fp.path))
				dlog.Debug("read watch event",
					dlog.FieldMod("file datasource"),
					dlog.String("event", filepath.Clean(event.Name)),
					dlog.String("path", filepath.Clean(fp.path)),
				)
				const writeOrCreateMask = fsnotify.Write | fsnotify.Create
				if event.Op&writeOrCreateMask != 0 && filepath.Clean(event.Name) == filepath.Clean(fp.path) {
					log.Println("modified file: ", event.Name)
					select {
					case fp.changed <- struct{}{}:
					default:
					}
				}
			case err := <-w.Errors:
				//log.Printf("read watch error: %s",err)
				dlog.Error("read watch error", dlog.FieldMod("file datasource"), dlog.Any("err", err))
			}
		}
	}()

	err = w.Add(fp.dir)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
