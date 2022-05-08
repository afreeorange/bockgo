package main

import (
	"fmt"
	"os"
	"path/filepath"

	// TODO: Implement this yourself
	cp "github.com/otiai10/copy"
)

// TODO: This can be smarter...
func copyTemplateAssets(outputFolder string) {
	fmt.Print("Creating template assets...")

	// Copy all the css, js, etc
	for _, a := range [3]string{"css", "img", "js"} {
		d, err := templatesContent.ReadDir("template/" + a)
		if err != nil {
			fmt.Print("Could not read " + a + "...skipping...")
			break
		}

		os.MkdirAll(outputFolder+"/"+a, os.ModePerm)

		for _, de := range d {
			f, _ := templatesContent.ReadFile("template/" + a + "/" + de.Name())
			os.WriteFile(outputFolder+"/"+a+"/"+de.Name(), f, os.ModePerm)
		}
	}

	// Then copy anything at the root level of the template folder except the
	// actual template HTML files!
	d, _ := templatesContent.ReadDir("template")
	for _, de := range d {
		if filepath.Ext(de.Name()) != ".html" {
			f, _ := templatesContent.ReadFile("template/" + de.Name())
			os.WriteFile(outputFolder+"/"+de.Name(), f, os.ModePerm)
		}
	}

	fmt.Println("done.")
}

func copyAssets(articleRoot string, outputFolder string) {
	fmt.Print("Copying assets... ")

	err := cp.Copy(articleRoot+"/__assets", outputFolder+"/assets")
	if err != nil {
		fmt.Print("Oops, could not copy assets: ", err)
	}

	fmt.Println("done.")
}
