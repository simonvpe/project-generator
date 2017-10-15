package main

import (
	"bytes"
	"log"
	"os"
	"strings"
	"text/template"
	"github.com/urfave/cli"
	"github.com/Pallinder/go-randomdata"
	"github.com/simonvpe/git"
)

const CPP int = 0

type Project struct {
	Name        string
	Language    int
	Tests       bool
	Extension   string
}

func generateName() string {
	var buffer bytes.Buffer
	buffer.WriteString(randomdata.Adjective())
	buffer.WriteString("-")
	buffer.WriteString(randomdata.FirstName(randomdata.RandomGender))
	return strings.ToLower(buffer.String())
}

func initializeRepository(proj Project) {
	log.Output(0, "Initializing git repository")
	_, err := git.Run(".", "init")
	if err != nil {
		log.Fatal(err)
	}
}

func createRootCmakeLists(proj Project) {
	log.Output(0, "Creating CMakeLists.txt")
	t := template.Must(template.ParseFiles("CMakeLists.tmpl"))
	f, err := os.OpenFile("CMakeLists.txt", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	
	err = t.Execute(f, proj)
	if err != nil {
		log.Fatal(err)
	}
}

func createTestsDirectory(proj Project) {
	_, err := os.Stat("tests")
	if proj.Tests && err != nil {
		log.Output(0, "Creating tests directory")
		err := os.Mkdir("tests", os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func createTestsCmakeLists(proj Project) {
	if proj.Tests {
		log.Output(0, "Creating tests/CMakeLists.txt")
		t := template.Must(template.ParseFiles("CMakeLists_tests.tmpl"))
		f, err := os.OpenFile("tests/CMakeLists.txt", os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		
		err = t.Execute(f, proj)		
		if err != nil {			
			log.Fatal(err)
		}

		
		if _, err = git.Run(".", "add", "tests/CMakeLists.txt"); err != nil {
			log.Fatal(err)
		}
	}
}

func addTestFramework(proj Project) {

}

func cppGenerator(ctx *cli.Context) {
	
	log.Output(0, "Generating C++ project")

	proj := Project {
		Language: CPP,
		Name: ctx.String("name"),
		Tests: ctx.Bool("tests"),
		Extension: "cpp",
	}

	initializeRepository(proj)
	createRootCmakeLists(proj)
	createTestsDirectory(proj)
	createTestsCmakeLists(proj)
	addTestFramework(proj)
}

func main() {
	app := cli.NewApp()
	app.Name = "project-generator"
	app.Usage = "generate a software project"

	app.Commands = []cli.Command {
		{
			Name: "c++",
			Aliases: []string{"cpp", "C++"},
			Usage: "generate a C++ project",
			Action: cppGenerator,
			Flags: []cli.Flag {
				cli.BoolFlag{
					Name: "tests",
					Usage: "include tests",
				},
				cli.StringFlag{
					Name: "name",
					Value: generateName(),
				},
			},
		},
	}

	app.Run(os.Args)
}
