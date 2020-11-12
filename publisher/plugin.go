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
				fmt.Println(err.Error())
			}

			new, err := p.Lookup("New")
			if err != nil {
				fmt.Println(err.Error())
			}

			fmt.Println(fmt.Sprintf("added publisher plugin %s", name))
			publisher := new.(func() Publisher)()
			m[name] = publisher
		}
	}

	return m
}
