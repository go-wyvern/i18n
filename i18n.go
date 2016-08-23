package i18n

import (
	"bytes"
	"fmt"
	"strings"

	"gopkg.in/yaml.v2"
	"io/ioutil"
	"io"
	"html/template"
)

type locale struct {
	lang  string
	trans map[string]interface{}
}

type localestore struct {
	langs []string
	store map[string]locale
}

var defaultLocaleStore = new(localestore)

func localeFile(lang string) string {
	return fmt.Sprintf("locale_%s.yml", lang)
}

func ParseFile(fileName string, trans map[string]interface{}) error {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	return Parse(data, trans)
}

func Parse(data []byte, trans map[string]interface{}) error {
	return yaml.Unmarshal(data, trans)
}

func SetupLocales(path, langs string) error {
	var err error
	defaultLocaleStore.langs = strings.Split(langs, "|")
	defaultLocaleStore.store=make(map[string]locale)
	for _, lang := range defaultLocaleStore.langs {
		loc := locale{}
		loc.trans= make(map[string]interface{})
		err = ParseFile(path + "/" + localeFile(lang), loc.trans)
		if err != nil {
			return err
		}
		defaultLocaleStore.store[lang] = loc
	}
	return nil
}

func Translate(tpl string, lang string) string {
	var buff = make([]byte,0)
	b := bytes.NewBuffer(buff)
	tmpl(b, tpl, defaultLocaleStore.store[lang].trans)
	return b.String()
}

func tmpl(w io.Writer, text string, data interface{}) error {
	t := template.New("top")
	template.Must(t.Parse(text))
	if err := t.Execute(w, data); err != nil {
		return err
	}

	return nil
}
