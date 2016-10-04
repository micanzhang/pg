package pg

type EntryMgr struct {
	persistance EntryPersistant
	Entries     map[string]map[string]Entry
}

func NewEntryMgr(persistance EntryPersistant) (*EntryMgr, error) {
	entries, err := persistance.Restore()
	if err != nil {
		return nil, err
	}
	return &EntryMgr{
		persistance: persistance,
		Entries:     sliceToMap(entries),
	}, nil
}

func (em *EntryMgr) Get(domain string, username string) (entry Entry, ok bool) {
	ee, ok := em.Entries[domain]
	if !ok {
		return
	}

	entry, ok = ee[username]
	return
}

func (em *EntryMgr) Put(entry Entry) {
	defer em.Flush()

	ee, ok := em.Entries[entry.Domain]
	if !ok {
		ee = make(map[string]Entry)
	}

	ee[entry.Username] = entry
	em.Entries[entry.Domain] = ee
}

func (em *EntryMgr) Remove(domain string, username string) bool {
	ee, ok := em.Entries[domain]
	if !ok {
		return false
	}

	if _, ok := ee[username]; !ok {
		return false
	}

	defer em.Flush()

	delete(ee, username)
	em.Entries[domain] = ee

	return true
}

func (em *EntryMgr) Flush() {
	if err := em.persistance.Save(mapToSlice(em.Entries)); err != nil {
		panic(err)
	}
}

func mapToSlice(entriesMap map[string]map[string]Entry) []Entry {
	var entries = make([]Entry, 0)
	for _, ee := range entriesMap {
		for _, e := range ee {
			entries = append(entries, e)
		}
	}

	return entries
}

func sliceToMap(entries []Entry) map[string]map[string]Entry {
	res := make(map[string]map[string]Entry)
	for _, e := range entries {
		ee, ok := res[e.Domain]
		if !ok {
			ee = make(map[string]Entry)
		}
		ee[e.Username] = e
		res[e.Domain] = ee
	}

	return res
}
