package pg

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestEntryMgr(t *testing.T) {
	path := fmt.Sprintf("entry_test_%d", time.Now().Unix())
	persistance := NewFileEntryPersistant(path)
	defer func() {
		os.Remove(path)
	}()

	em, err := NewEntryMgr(persistance)
	if err != nil {
		t.Error(err)
	}

	entry := Entry{
		Domain:   "github.com",
		Username: "micanzhang",
		Password: Password("**********"),
	}

	// put
	em.Put(entry)

	// get
	entry1, ok := em.Get(entry.Domain, entry.Username)
	if !ok {
		t.Fatalf("me.Get(%s, %s) failed: not found", entry.Domain, entry.Username)
	}

	if entryEqual(entry, entry1) == false {
		t.Errorf("expected: %+v, got: %+v", entry, entry1)
	}

	// remove
	ok = em.Remove(entry.Domain, entry.Username)
	if !ok {
		t.Error("em.Remove  failed")
	}

	// get
	_, ok = em.Get(entry.Domain, entry.Username)
	if ok {
		t.Error("expected: false, got: true")
	}
}

func entryEqual(e1, e2 Entry) bool {
	if e1.Domain != e2.Domain {
		return false
	}

	if e1.Username != e2.Username {
		return false
	}

	if e1.Password != e2.Password {
		return false
	}

	return true
}
