package isso

import (
	"wrong.wang/x/go-isso/isso/model"
)

// Storage handles all operations related to the database.
type Storage interface {
	ThreadStorage
	CommentStorage
	PreferenceStorage
}

// ThreadStorage handles all operations related to Thread and the database.
type ThreadStorage interface {
	ContainsThread(uri string) bool
	GetThreadByURI(uri string) (model.Thread, error)
	GetThreadByID(id int) (model.Thread, error)
	NewThread(uri string, title string, Host string) (model.Thread, error)
}

// CommentStorage handles all operations related to Comment and the database.
type CommentStorage interface {
	IsApprovedAuthor(email string) bool
	NewComment(model.AcceptComment) (model.ReplyComment, error)
}

// PreferenceStorage handles all operations related to Preference and the database.
type PreferenceStorage interface {
	GetPreference(key string) (string, error)
	SetPreference(key string, value string) error
}

// type validThreads struct {
// 	sync.RWMutex
// 	data map[string]model.Thread
// 	ts   ThreadStorage
// }

// func NewValidThreads(ts ThreadStorage) (*validThreads, error) {
// 	var data map[string]model.Thread
// 	mt, err := ts.GetAllThreds()
// 	if err != nil {
// 		return nil, err
// 	}
// 	for _, m := range mt {
// 		data[m.URI] = m
// 	}
// 	return &validThreads{
// 		sync.RWMutex{},
// 		data,
// 		ts,
// 	}, nil
// }

// func (vt *validThreads) Contains(uri string) bool {
// 	vt.RLock()
// 	defer vt.RUnlock()
// 	_, ok := vt.data[uri]
// 	return ok
// }

// func (vt *validThreads) Get(uri string) (model.Thread, error) {
// 	vt.RLock()
// 	defer vt.RUnlock()
// 	var mt model.Thread
// 	if mt, ok := vt.data[uri]; !ok {
// 		return mt, errors.New("thread not found")
// 	}
// 	return mt, nil
// }

// // Add new uri to validThreads
// func (vt *validThreads) Add(uri string, title string, Host string) (model.Thread, error) {
// 	if title == "" {
// 		url := path.Join(Host, uri)
// 		resp, err := http.Get(url)
// 		if err != nil {
// 			return model.Thread{}, fmt.Errorf("failed to load page %s (%w)", url, err)
// 		}
// 		defer resp.Body.Close()
// 		if resp.StatusCode != 200 {
// 			return model.Thread{}, fmt.Errorf("can't load page %s, code %d", url, resp.StatusCode)
// 		}
// 		title, uri, err = extract.TitleAndThreadURI(resp.Body, "Untitled", uri)
// 		if err != nil {
// 			return model.Thread{}, err
// 		}
// 	}
// 	vt.Lock()
// 	defer vt.Unlock()
// 	t, err := vt.ts.NewThread(uri, title)
// 	if err != nil {
// 		return model.Thread{}, fmt.Errorf("can't save thread into database :%w", err)
// 	}
// 	vt.data[uri] = t
// 	return t, nil
// }
