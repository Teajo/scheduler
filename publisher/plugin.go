package publisher

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"plugin"
)

var pluginExt = ".so"

func loadPublisherPlugins(dir string) map[string]Publisher {
	fmt.Println(fmt.Sprintf("load plugins in %s", dir))

	m := make(map[string]Publisher)
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if !f.IsDir() && filepath.Ext(f.Name()) == pluginExt {
			fmt.Println("lookup plugin", f.Name())

			var name = f.Name()[0 : len(f.Name())-len(pluginExt)]

			p, err := plugin.Open(filepath.Join("plugins", f.Name()))
			if err != nil {
				fmt.Println(fmt.Sprintf("%s unable to cast New function, plugin %s will not be added", err.Error(), name))
				continue
			}

			new, err := p.Lookup("New")
			if err != nil {
				fmt.Println(fmt.Sprintf("%s unable to cast New function, plugin %s will not be added", err.Error(), name))
				continue
			}

			newPublisher, ok := new.(func() Publisher)
			if !ok {
				fmt.Println(fmt.Sprintf("unable to cast New function, plugin %s will not be added", name))
				continue
			}

			publisher := newPublisher()
			m[name] = publisher
			fmt.Println(fmt.Sprintf("publisher plugin %s added", name))
		}
	}

	return m
}
