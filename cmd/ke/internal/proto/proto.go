package proto

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/mod/modfile"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// proto 表示一个 API proto 文件的生成描述
type proto struct {
	Name        string // demo
	Dir         string // api/demo/v1
	File        string // demo.proto
	Package     string // demo.v1
	Service     string // Demo
	GoPackage   string
	JavaPackage string
}

// newProto 根据 API 名称构建 proto 对象
func newProto(name, version, protoName string) *proto {
	name = strings.ToLower(name)

	dir := filepath.Join("api", name, version)

	return &proto{
		Name:        name,
		Dir:         dir,
		File:        protoName + ".proto",
		Package:     "api." + name + "." + version,
		Service:     toUpperCamelCase(protoName),
		GoPackage:   goPackage(dir),
		JavaPackage: protoName,
	}
}

// FilePath 返回最终生成的 proto 文件路径
func (p *proto) FilePath() string {
	return filepath.Join(p.Dir, p.File)
}

// Generate 生成 proto 文件
func (p *proto) Generate() error {
	body, err := p.execute()
	if err != nil {
		return err
	}

	// 创建目录
	if err := os.MkdirAll(p.Dir, 0o755); err != nil {
		return err
	}

	path := p.FilePath()

	// 防止覆盖已有文件
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("file already exists：%s", path)
	}

	if err := os.WriteFile(path, body, 0o644); err != nil {
		return err
	}

	return nil
}

// goPackage 生成 go_package 选项
func goPackage(dir string) string {
	mod := findModuleName()
	if mod == "" {
		return dir
	}
	base := filepath.Base(dir)
	return fmt.Sprintf("%s/%s;%s", mod, dir, base)
}

// findModuleName 向上查找 go.mod 并解析 module 名称
func findModuleName() string {
	dir, _ := os.Getwd()

	for {
		modPath := filepath.Join(dir, "go.mod")
		if data, err := os.ReadFile(modPath); err == nil {
			return modfile.ModulePath(data)
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return ""
}

// toUpperCamelCase 将 demo_name 转为 DemoName
func toUpperCamelCase(s string) string {
	s = strings.ReplaceAll(s, "_", " ")
	s = cases.Title(language.Und, cases.NoLower).String(s)
	return strings.ReplaceAll(s, " ", "")
}
