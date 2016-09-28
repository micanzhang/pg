package main

import (
	"flag"
	"fmt"
	"os"

	"strings"

	"github.com/micanzhang/tools/pg"
)

// http://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
var (
	length    int
	upperCase int
	lowerCase int
	number    int
	char      int
	domain    string
	username  string
	password  string
	key       string
	mgr       *pg.EntryMgr
)

func init() {
	// TODO
	flag.IntVar(&length, "l", 16, "length of password")
	flag.IntVar(&upperCase, "U", 2, "at least width of upperCase chars")
	flag.IntVar(&lowerCase, "L", 2, "at least width of lowerCase chars")
	flag.IntVar(&number, "N", 2, "at least width of numbers")
	flag.IntVar(&char, "C", 2, "at least width of chars")
	flag.StringVar(&domain, "d", "", "domain")
	flag.StringVar(&username, "u", "", "username")
	flag.StringVar(&password, "p", "", "password")
	flag.StringVar(&key, "k", "", "key")

	path := fmt.Sprintf("%s/.pg", os.Getenv("HOME"))
	persistance := pg.NewFileEntryPersistant(path)
	var err error

	mgr, err = pg.NewEntryMgr(persistance)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error ocurred: %s", err)
	}

}

func main() {
	var (
		action string
		args   = os.Args
	)

	if len(os.Args) > 1 {
		// filter action
		action = os.Args[1]
		if strings.HasPrefix(action, "-") {
			action = ""
		} else {
			os.Args = os.Args[1:]
		}
	}

	if action == "" {
		action = "gen"
	}

	flag.Parse()

	os.Args = args

	switch action {
	case "gen":
		fmt.Println(pg.GenRandPassword(length, upperCase, lowerCase, number, char))
	case "list":
		list()
	case "new":
		newAction()
	case "update":
		updateAction()
	case "remove":
		removeAction()
	case "info":
		infoAction()
	default:
		Usage()
	}
}

func Usage() {
	fmt.Fprintf(os.Stderr, "%s is a open source command line tool for password manager.\n\n", os.Args[0])
	fmt.Fprint(os.Stderr, "Usage:\n")
	fmt.Fprintf(os.Stderr, "\t%s [command] [arguments]\n\n", os.Args[0])
	fmt.Fprint(os.Stderr, "the commands are:\n")
	fmt.Fprint(os.Stderr, `
    gen       generate password
    list      list all persistant entries
    new       new entry by domain and username
    update    update entry's password by domain and username
    remove    remove entry by domain and username
    info      get info of entry by domain and username
`)
	fmt.Fprint(os.Stderr, "\nthe arguments are:\n\n")

	flag.PrintDefaults()
	fmt.Fprint(os.Stderr, "\n")
}

func removeAction() {
	if mgr.Remove(domain, username) {
		fmt.Println("done")
	} else {
		fmt.Println("nothing to do")
	}
}

func newAction() {
	if key == "" {
		panic("key is required")
	}
	if domain == "" {
		panic("domain is required")
	}

	entry, ok := mgr.Get(domain, username)
	if ok {
		panic(fmt.Sprintf("%s has exists", username))
	}

	if password == "" {
		password = pg.GenRandPassword(length, upperCase, lowerCase, number, char)
	}

	passwd, err := pg.NewPassword(key, password)
	if err != nil {
		panic(err)
	}

	entry = pg.Entry{
		Domain:   domain,
		Username: username,
		Password: passwd,
	}

	mgr.Put(entry)
}

func updateAction() {
	if key == "" {
		panic("key is required")
	}
	if domain == "" {
		panic("domain is required")
	}

	entry, ok := mgr.Get(domain, username)
	if !ok {
		fmt.Fprintln(os.Stdout, "Not found")
		return
	}

	if password == "" {
		password = pg.GenRandPassword(length, upperCase, lowerCase, number, char)
	}

	passwd, err := pg.NewPassword(key, password)
	if err != nil {
		panic(err)
	}

	entry.Password = passwd
	mgr.Put(entry)
}

func infoAction() {
	entry, ok := mgr.Get(domain, username)
	if !ok {
		fmt.Fprintln(os.Stderr, "Not found")
		return
	}

	if key == "" {
		fmt.Fprintln(os.Stderr, "-k is required")
		return
	}

	plaintext, err := entry.Password.Plaintext(key)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid -k paramsters: %s ", key)
	}

	fmt.Println(plaintext)
}

func list() {
	if len(mgr.Entries) == 0 {
		fmt.Println("empty list.")
		return
	}

	for _, ee := range mgr.Entries {
		for _, e := range ee {
			printEntry(e)
		}
	}
}

func printEntry(entry pg.Entry) {
	//fmt.Printf("{\n\tdomain: %s,\n\tusername: %s\n}\n", entry.Domain, entry.Username)
	//fmt.Printf("domain: %s, username: %s\n", entry.Domain, entry.Username)
	fmt.Printf("%s: %s\n", entry.Domain, entry.Username)
}