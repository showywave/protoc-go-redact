package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const DestFilePrefix = "redact_"

func main() {
	var input string
	flag.StringVar(&input, "input", "", "input a Go file")
	flag.Parse()

	if len(input) == 0 {
		log.Fatal("input file not found")
	}

	do(input)
}

func do(file string) {
	finfo, err := os.Stat(file)
	if err != nil {
		log.Fatal(err)
	}
	if finfo.IsDir() {
		return
	}

	if !strings.HasSuffix(strings.ToLower(finfo.Name()), ".go") ||
		strings.Contains(strings.ToLower(finfo.Name()), DestFilePrefix) {
		return
	}

	dmd := &Demand{
		Sign:     "@ban",
		FuncName: "Redact",
		Buf:      &bytes.Buffer{},
	}

	genRedact(file, dmd)

	writeFile(file, dmd)
}

func writeFile(file string, dmd *Demand) {
	dir, srcFileName := filepath.Split(file)
	fmt.Println(dir)
	destFilePath := dir + "/" + DestFilePrefix + srcFileName
	err := os.WriteFile(destFilePath, dmd.Buf.Bytes(), 0o644)
	if err != nil {
		log.Fatal(err)
	}
}

// Demand
type Demand struct {
	Sign     string
	FuncName string
	Buf      *bytes.Buffer
}

func genRedact(file string, dmd *Demand) {
	fset := token.NewFileSet()
	astFile, _ := parser.ParseFile(fset, file, nil, parser.ParseComments)

	writeLine(dmd.Buf, "package "+astFile.Name.String())
	writeLine(dmd.Buf)
	writeLine(dmd.Buf, `import "fmt"`)
	writeLine(dmd.Buf)

	for _, decl := range astFile.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		var typeSpec *ast.TypeSpec
		for _, spec := range genDecl.Specs {
			if ts, tsOK := spec.(*ast.TypeSpec); tsOK {
				typeSpec = ts
				break
			}
		}

		if typeSpec == nil {
			continue
		}

		structDecl, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			continue
		}

		kl := []string{}
		vl := []string{}
		needGen := false
		for _, field := range structDecl.Fields.List {
			comments := []*ast.Comment{}
			if field.Doc != nil {
				comments = append(comments, field.Doc.List...)
			}

			if field.Comment != nil {
				comments = append(comments, field.Comment.List...)
			}

			fieldType, ok := field.Type.(*ast.Ident)
			if !ok {
				continue
			}

			kt := "%v"
			switch fieldType.Name {
			case "string":
				kt = "%s"
			case "int", "int32":
				kt = "%d"
			default:
				kt = "%v"
			}

			pbname := GetPBNameInTag(field.Tag.Value)
			kl = append(kl, pbname+":"+kt)

			v := "x." + field.Names[0].Name
			for _, comment := range comments {
				if strings.Contains(comment.Text, dmd.Sign) {
					v = `"******"`
					needGen = true
				}
			}
			vl = append(vl, v)
		}

		if needGen {
			writeLine(dmd.Buf, "func (x *", typeSpec.Name, ") "+dmd.FuncName+" () string {")
			writeLine(dmd.Buf, `	return fmt.Sprintf("`, strings.Join(kl, " "), `", `, strings.Join(vl, ", "), ")")
			writeLine(dmd.Buf, "}")
			writeLine(dmd.Buf)
		}
	}
}

func writeLine(buf *bytes.Buffer, v ...interface{}) {
	for _, x := range v {
		fmt.Fprint(buf, x)
	}
	fmt.Fprintln(buf)
}

func GetPBNameInTag(input string) string {
	var pbname string
	var tagStr = input
	tagStr = strings.Replace(tagStr, "`", "", -1)
	tagStr = strings.Replace(tagStr, "\"", "", -1)
	tagList := strings.Split(tagStr, " ")
	for _, tag := range tagList {
		tagArr := strings.Split(tag, ":")
		values := strings.Split(tagArr[1], ",")
		if tagArr[0] == "protobuf" {
			for _, v := range values {
				if strings.Contains(v, "name") {
					pbname = strings.Split(v, "=")[1]
				}
			}
		}
	}

	return pbname
}