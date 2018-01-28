package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
)

const (
	javaOuterClassSuffix = "OuterClass"
)

func getProtoName(file *descriptor.FileDescriptorProto) string {
	name := file.GetName()
	ext := filepath.Ext(name)
	if ext == ".proto" || ext == ".protodevel" {
		name = name[0 : len(name)-len(ext)]
	}
	return name
}

func getJavaOuterClassName(file *descriptor.FileDescriptorProto) string {
	name := file.Options.GetJavaOuterClassname()
	if name != "" {
		return name
	}

	name = camelCase(getProtoName(file))
	outer := name + javaOuterClassSuffix
	for _, desc := range file.GetMessageType() {
		if strings.Title(desc.GetName()) == name {
			return outer
		}
	}

	for _, desc := range file.GetService() {
		if strings.Title(desc.GetName()) == name {
			return outer
		}
	}

	for _, desc := range file.GetEnumType() {
		if strings.Title(desc.GetName()) == name {
			return outer
		}
	}

	return name
}

func getJavaOuterClassFile(file *descriptor.FileDescriptorProto) string {
	className := getJavaOuterClassName(file)
	pkg := getJavaPackage(file)
	if pkg == "" {
		return fmt.Sprintf("%s.java", className)
	} else {
		dir := strings.Replace(pkg, ".", "/", -1)
		return fmt.Sprintf("%s/%s.java", dir, className)
	}
}

func getJavaServiceClassName(file *descriptor.FileDescriptorProto, service *descriptor.ServiceDescriptorProto) string {
	serviceName := camelCase(service.GetName())
	return fmt.Sprintf("%s", serviceName)
}

func getJavaServiceClassFile(file *descriptor.FileDescriptorProto, service *descriptor.ServiceDescriptorProto) string {
	serviceClass := getJavaServiceClassName(file, service)
	pkg := getJavaPackage(file)
	if pkg == "" {
		return fmt.Sprintf("%s.java", serviceClass)
	} else {
		dir := strings.Replace(pkg, ".", "/", -1)
		return fmt.Sprintf("%s/%s.java", dir, serviceClass)
	}
}

func getJavaServiceClientClassName(file *descriptor.FileDescriptorProto, service *descriptor.ServiceDescriptorProto) string {
	serviceName := camelCase(service.GetName())
	return fmt.Sprintf("%sClient", serviceName)
}

func getJavaServiceClientClassFile(file *descriptor.FileDescriptorProto, service *descriptor.ServiceDescriptorProto) string {
	serviceClass := getJavaServiceClientClassName(file, service)
	pkg := getJavaPackage(file)
	if pkg == "" {
		return fmt.Sprintf("%s.java", serviceClass)
	} else {
		dir := strings.Replace(pkg, ".", "/", -1)
		return fmt.Sprintf("%s/%s.java", dir, serviceClass)
	}
}

func getJavaPackage(file *descriptor.FileDescriptorProto) string {
	pkg := file.Options.GetJavaPackage()
	if pkg != "" {
		return pkg
	}
	return file.GetPackage()
}

func getJavaType(file *descriptor.FileDescriptorProto, name string) string {
	pkg := getJavaPackage(file)
	multi := file.Options.GetJavaMultipleFiles()

	p := strings.LastIndex(name, ".")
	typeName := name[p+1:]

	if pkg == "" {
		return fmt.Sprintf("%s", typeName)
	} else if multi {
		return fmt.Sprintf("%s.%s", pkg, typeName)
	} else {
		outerClass := getJavaOuterClassName(file)
		return fmt.Sprintf("%s.%s.%s", pkg, outerClass, typeName)
	}
}

func camelCase(str string) string {
	parts := strings.Split(str, "_")
	for i, part := range parts {
		runes := []rune(part)
		runes[0] = unicode.ToUpper(runes[0])
		parts[i] = string(runes)
	}
	return strings.Join(parts, "")
}

func lowerCamelCase(str string) string {
	cc := camelCase(str)
	runes := []rune(cc)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}
