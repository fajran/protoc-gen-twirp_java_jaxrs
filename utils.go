package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"unicode"

	"google.golang.org/protobuf/types/descriptorpb"
)

const (
	javaOuterClassSuffix = "OuterClass"
)

func getProtoName(file *descriptorpb.FileDescriptorProto) string {
	name := file.GetName()
	ext := filepath.Ext(name)
	if ext == ".proto" || ext == ".protodevel" {
		name = name[0 : len(name)-len(ext)]
	}
	return name
}

func getJavaOuterClassName(file *descriptorpb.FileDescriptorProto) string {
	name := file.Options.GetJavaOuterClassname()
	if name != "" {
		return name
	}

	name = camelCase(getProtoName(file))
	outer := name + javaOuterClassSuffix
	for _, desc := range file.MessageType {
		if strings.Title(desc.GetName()) == name {
			return outer
		}
	}

	for _, desc := range file.Service {
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

func getJavaOuterClassFile(file *descriptorpb.FileDescriptorProto) string {
	className := getJavaOuterClassName(file)
	pkg := getJavaPackage(file)
	if pkg == "" {
		return fmt.Sprintf("%s.java", className)
	} else {
		dir := strings.Replace(pkg, ".", "/", -1)
		return fmt.Sprintf("%s/%s.java", dir, className)
	}
}

func containsType(name string, file *descriptorpb.FileDescriptorProto) bool {
	for _, t := range file.GetEnumType() {
		if t.GetName() == name {
			return true
		}
	}
	for _, t := range file.GetMessageType() {
		if t.GetName() == name {
			return true
		}
	}
	for _, t := range file.GetService() {
		if t.GetName() == name {
			return true
		}
	}
	return false
}

func getJavaServiceClassName(file *descriptorpb.FileDescriptorProto, service *descriptorpb.ServiceDescriptorProto) string {
	serviceName := camelCase(service.GetName())
	return fmt.Sprintf("%s", serviceName)
}

func getJavaServiceClassFile(file *descriptorpb.FileDescriptorProto, service *descriptorpb.ServiceDescriptorProto) string {
	serviceClass := getJavaServiceClassName(file, service)
	pkg := getJavaPackage(file)
	if pkg == "" {
		return fmt.Sprintf("%s.java", serviceClass)
	} else {
		dir := strings.Replace(pkg, ".", "/", -1)
		return fmt.Sprintf("%s/%s.java", dir, serviceClass)
	}
}

func getJavaServiceClientClassName(file *descriptorpb.FileDescriptorProto, service *descriptorpb.ServiceDescriptorProto) string {
	serviceName := camelCase(service.GetName())
	name := fmt.Sprintf("%sClient", serviceName)
	if containsType(name, file) {
		name = fmt.Sprintf("%sTwirpClient", serviceName)
	}
	return name
}

func getJavaServiceClientClassFile(file *descriptorpb.FileDescriptorProto, service *descriptorpb.ServiceDescriptorProto) string {
	serviceClass := getJavaServiceClientClassName(file, service)
	return getJavaServiceClientClassFileByString(file, serviceClass)
}

func getJavaServiceClientClassFileByString(file *descriptorpb.FileDescriptorProto, serviceClass string) string {
	pkg := getJavaPackage(file)
	if pkg == "" {
		return fmt.Sprintf("%s.java", serviceClass)
	} else {
		dir := strings.Replace(pkg, ".", "/", -1)
		return fmt.Sprintf("%s/%s.java", dir, serviceClass)
	}
}

func getJavaPackage(file *descriptorpb.FileDescriptorProto) string {
	pkg := file.Options.GetJavaPackage()
	if pkg != "" {
		return pkg
	}
	return file.GetPackage()
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
