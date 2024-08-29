package main

import (
	"errors"
)

type Language struct {
	Name         string
	Extensions   []string
	Boilerplate  func(string) string
	CompilerName string
	CompilerArgs []string
}

var langs = [5]Language{
	{Name: "java", Extensions: []string{".java"}, Boilerplate: javaBoilerplate, CompilerName: "javac", CompilerArgs: []string{}},
	{Name: "go", Extensions: []string{".go"}, Boilerplate: goBoileplate, CompilerName: "", CompilerArgs: []string{}},
	{Name: "c", Extensions: []string{".c"}, Boilerplate: cBoileplate, CompilerName: "", CompilerArgs: []string{}},
	{Name: "python", Extensions: []string{".py"}, Boilerplate: pythonBoileplate, CompilerName: "", CompilerArgs: []string{}},
	{Name: "c#", Extensions: []string{".cs"}, Boilerplate: csBoileplate, CompilerName: "", CompilerArgs: []string{}},
}

func GetLanguage(name string) (Language, error) {
	for _, lang := range langs {
		if name == lang.Name {
			return lang, nil
		}
	}
	return Language{}, errors.New("cannot find language")
}

func csBoileplate(s string) string {
	return ""
}

func pythonBoileplate(s string) string {
	return ""
}

func cBoileplate(s string) string {
	return ""
}

func goBoileplate(s string) string {
	return ""
}

func javaBoilerplate(projectName string) string {
	capatalizedProjectName := Capatalize(projectName)
	return "import java.util.Scanner;\n" + "\n" + "public class " + capatalizedProjectName + " {\n" + "   public static void main(String[] args) {\n" + "       Scanner sc = new Scanner(System.in);\n" + "       sc.close();\n" + "   }\n" + "}"
}
