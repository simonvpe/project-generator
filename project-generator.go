package main

import (
	"bytes"
	"log"
	"os"
	"path"
	"strings"
	"github.com/urfave/cli"
	"github.com/Pallinder/go-randomdata"
	"github.com/simonvpe/git"
	"github.com/simonvpe/cmake"
	"github.com/kelseyhightower/envconfig"
)

type Environment struct {
	Share string `default:"/usr/share/myapp"`
}

func generateName() string {
	var buffer bytes.Buffer
	buffer.WriteString(randomdata.Adjective())
	buffer.WriteString("-")
	buffer.WriteString(randomdata.FirstName(randomdata.RandomGender))
	return strings.ToLower(buffer.String())
}

func cppGenerator(ctx *cli.Context) {
	var env Environment
	envconfig.Process("myapp", &env)
	
	log.Output(0, "Generating C++ project")

	cmakeCtx := cmake.CMakeContext {
		ProjectName: ctx.String("name"),
		MinimumVersion: "3.8",
		TestSuite: ctx.String("tests"),
		Language: "cpp",
	}
	
	log.Output(0, "Initializing git repository")
	if _, err := git.Run(".", "init"); err != nil {
		log.Fatalf("Failed to initialize repository %q\n", err)
	}

	log.Output(0, "Generating CMakeLists.txt")
	if err := cmake.Generate("CMakeLists.txt",
		path.Join(env.Share, "CMakeLists.tmpl"),
		&cmakeCtx); err != nil {
		log.Fatalf("Failed to generate CMakeLists.txt %q", err)
	}

	if cmakeCtx.TestSuite == "catch" {

		if _, err := os.Stat("test"); os.IsNotExist(err) {
			log.Output(0, "Creating test directory")
			if err := os.Mkdir("test", os.ModePerm); err != nil {
				log.Fatalf("Failed to create test directory %q", err)
			}
		}
		
		log.Output(0, "Generating test/CMakeLists.txt")
		if err := cmake.Generate(path.Join("test", "CMakeLists.txt"),
			path.Join(env.Share, "test", "CMakeLists.tmpl"),
			&cmakeCtx); err != nil {
			log.Fatalf("Failed to generate CMakeLists.txt %q", err)
		}

		if _, err := git.Run("test",
			"submodule",
			"add",
			"https://github.com/philsquared/Catch.git"); err != nil {
			log.Fatalf("Failed to add submodule %q", err)
		}

		if _, err := git.Run("test",
			"submodule",
			"update",
			"--init"); err != nil {
			log.Fatalf("Failed to init and update submodules %q", err)
		}
	}
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
				cli.StringFlag{
					Name: "tests",
					Value: "catch",
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
