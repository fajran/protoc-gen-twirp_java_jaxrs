package main

import (
    "bytes"
    "fmt"
    "google.golang.org/protobuf/proto"
    "google.golang.org/protobuf/types/descriptorpb"
    "google.golang.org/protobuf/types/pluginpb"
    "strings"
)

type generator struct {
    Request  *pluginpb.CodeGeneratorRequest
    Response *pluginpb.CodeGeneratorResponse

    output *bytes.Buffer
    indent string
}

func newGenerator(req *pluginpb.CodeGeneratorRequest) *generator {
    return &generator{
        Request:  req,
        Response: nil,
        output:   bytes.NewBuffer(nil),
        indent:   "",
    }
}

func (g *generator) Generate() error {
    g.Response = &pluginpb.CodeGeneratorResponse{}

    g.Response.File = append(g.Response.File, g.generateResponseMetaDataProcessorInterface())
    for _, file := range g.getProtoFiles() {
        err := g.processFile(file)
        if err != nil {
            return err
        }
    }

    return nil
}

func (g *generator) processFile(file *descriptorpb.FileDescriptorProto) error {
    if file.Options.GetJavaGenericServices() {
        return fmt.Errorf("twirp_java_jaxrs cannot not work with java_generic_services option")
    }

    g.Response.File = append(g.Response.File, g.generateProvider(file))

    for _, service := range file.Service {
        out := g.generateServiceInterface(file, service)
        g.Response.File = append(g.Response.File, out)

        out = g.generateServiceClient(file, service)
        g.Response.File = append(g.Response.File, out)
    }

    return nil
}

func (g *generator) generateResponseMetaDataProcessorInterface() *pluginpb.CodeGeneratorResponse_File {

    serviceName := "IResponseMetaDataProcessor"

    g.P(`// Code generated by protoc-gen-twirp_java_jaxrs, DO NOT EDIT.`)
    g.P()
    g.P(`package plugin;`)
    g.P()
    g.P(`public interface `, serviceName, ` {`)
    g.P()
    g.P(`    void process(javax.ws.rs.core.Response response);`)
    g.P(`}`)
    g.P()

    out := &pluginpb.CodeGeneratorResponse_File{}
    out.Content = proto.String(g.output.String())
    out.Name = proto.String(fmt.Sprintf("plugin/%s.java", serviceName))
    g.Reset()
    return out
}

func (g *generator) generateProvider(file *descriptorpb.FileDescriptorProto) *pluginpb.CodeGeneratorResponse_File {

    multi := file.Options.GetJavaMultipleFiles()
    serviceName := "ProtoBufMessageProvider"

    if multi {
        pkg := getJavaPackage(file)
        g.P(`// Code generated by protoc-gen-twirp_java_jaxrs, DO NOT EDIT.`)
        g.P(`// source: `, file.GetName())
        g.P()
        if pkg != "" {
            g.P(`package `, pkg, `;`)
            g.P()
        }
    }

    static := ""
    if !multi {
        static = "static "
    }
    g.P(`@javax.ws.rs.ext.Provider`)
    g.P(`@javax.ws.rs.Produces({"application/protobuf", "application/json"})`)
    g.P(`@javax.ws.rs.Consumes({"application/protobuf", "application/json"})`)
    g.P(`public `, static, `class `, serviceName, ` implements javax.ws.rs.ext.MessageBodyWriter<com.google.protobuf.Message>, javax.ws.rs.ext.MessageBodyReader<com.google.protobuf.Message> {`)
    g.P()
    g.P(`    @Override`)
    g.P(`    public boolean isWriteable(Class<?> type, java.lang.reflect.Type genericType, java.lang.annotation.Annotation[] annotations, javax.ws.rs.core.MediaType mediaType) {`)
    g.P(`        return com.google.protobuf.Message.class.isAssignableFrom(type) && ("json".equals(mediaType.getSubtype()) || "protobuf".equals(mediaType.getSubtype()));`)
    g.P(`    }`)
    g.P()
    g.P(`    @Override`)
    g.P(`    public long getSize(com.google.protobuf.Message t, Class<?> type, java.lang.reflect.Type genericType, java.lang.annotation.Annotation[] annotations, javax.ws.rs.core.MediaType mediaType) {`)
    g.P(`        if (t == null) {`)
    g.P(`            return -1;`)
    g.P(`        }`)
    g.P(`        java.io.ByteArrayOutputStream out = new java.io.ByteArrayOutputStream();`)
    g.P(`        try {`)
    g.P(`            writeTo(t, type, genericType, annotations, mediaType, null, out);`)
    g.P(`        } catch (java.io.IOException e) {`)
    g.P(`            return -1;`)
    g.P(`        }`)
    g.P(`        return out.size();`)
    g.P(`    }`)
    g.P()
    g.P(`    @Override`)
    g.P(`    public void writeTo(com.google.protobuf.Message t, Class<?> type, java.lang.reflect.Type genericType, java.lang.annotation.Annotation[] annotations, javax.ws.rs.core.MediaType mediaType, javax.ws.rs.core.MultivaluedMap<String, Object> httpHeaders, java.io.OutputStream entityStream) throws java.io.IOException, javax.ws.rs.WebApplicationException {`)
    g.P(`        switch (mediaType.getSubtype()) {`)
    g.P(`            case "protobuf":`)
    g.P(`                t.writeTo(entityStream);`)
    g.P(`                break;`)
    g.P(`            case "json":`)
    g.P(`                entityStream.write(com.google.protobuf.util.JsonFormat.printer().print(t).getBytes("UTF-8"));`)
    g.P(`                break;`)
    g.P(`            default:`)
    g.P(`                throw new javax.ws.rs.WebApplicationException("MediaType not supported!");`)
    g.P(`        }`)
    g.P()
    g.P(`    }`)
    g.P()
    g.P(`    @Override`)
    g.P(`    public boolean isReadable(Class<?> type, java.lang.reflect.Type genericType, java.lang.annotation.Annotation[] annotations, javax.ws.rs.core.MediaType mediaType) {`)
    g.P(`        return com.google.protobuf.Message.class.isAssignableFrom(type) && ("json".equals(mediaType.getSubtype()) || "protobuf".equals(mediaType.getSubtype()));`)
    g.P(`    }`)
    g.P()
    g.P(`    @Override`)
    g.P(`    public com.google.protobuf.Message readFrom(Class<com.google.protobuf.Message> type, java.lang.reflect.Type genericType, java.lang.annotation.Annotation[] annotations, javax.ws.rs.core.MediaType mediaType, javax.ws.rs.core.MultivaluedMap<String, String> httpHeaders, java.io.InputStream entityStream) throws java.io.IOException, javax.ws.rs.WebApplicationException {`)
    g.P(`        try {`)
    g.P(`            switch (mediaType.getSubtype()) {`)
    g.P(`                case "protobuf":`)
    g.P(`                    java.lang.reflect.Method m = type.getMethod("parseFrom", java.io.InputStream.class);`)
    g.P(`                    return (com.google.protobuf.Message) m.invoke(null, entityStream);`)
    g.P(`                case "json":`)
    g.P(`                    com.google.protobuf.Message.Builder msg = (com.google.protobuf.Message.Builder)type.getMethod("newBuilder").invoke(null);`)
    g.P(`                    com.google.protobuf.util.JsonFormat.parser().merge(new java.io.InputStreamReader(entityStream), msg);`)
    g.P(`                    return msg.build();`)
    g.P(`                default:`)
    g.P(`                    throw new javax.ws.rs.WebApplicationException("MediaType not supported!");`)
    g.P(`            }`)
    g.P(`        } catch (Exception e) {`)
    g.P(`            throw new javax.ws.rs.WebApplicationException(e);`)
    g.P(`        }`)
    g.P(`    }`)
    g.P(`}`)
    g.P()

    out := &pluginpb.CodeGeneratorResponse_File{}
    out.Content = proto.String(g.output.String())
    if multi {
        out.Name = proto.String(getJavaServiceClientClassFileByString(file, serviceName))
    } else {
        out.Name = proto.String(getJavaOuterClassFile(file))
        out.InsertionPoint = proto.String("outer_class_scope")
    }
    g.Reset()
    return out
}

func (g *generator) generateServiceClient(file *descriptorpb.FileDescriptorProto, service *descriptorpb.ServiceDescriptorProto) *pluginpb.CodeGeneratorResponse_File {
    multi := file.Options.GetJavaMultipleFiles()

    if multi {
        pkg := getJavaPackage(file)
        g.P(`// Code generated by protoc-gen-twirp_java_jaxrs, DO NOT EDIT.`)
        g.P(`// source: `, file.GetName())
        g.P()
        if pkg != "" {
            g.P(`package `, pkg, `;`)
            g.P()
        }
    }

    // TODO add comment

    serviceClass := getJavaServiceClientClassName(file, service)
    servicePath := g.GetServicePath(file, service)
    interfaceClass := g.GetType(file, getJavaServiceClassName(file, service))
    provider := g.GetType(file, "ProtoBufMessageProvider")

    static := ""
    if !multi {
        static = "static "
    }
    g.P(`public `, static, `class `, serviceClass, ` implements `, interfaceClass, ` {`)
    g.P(`  private final javax.ws.rs.client.WebTarget target;`)
    g.P(`  private plugin.IResponseMetaDataProcessor responseMetaDataProcessor;`)
    g.P()
    g.P(`  public `, serviceClass, `(javax.ws.rs.client.WebTarget target) {`)
    g.P(`    this.target = target.path("`, servicePath, `");`)
    g.P(`    this.target.register(new `, provider, `());`)
    g.P(`  }`)
    g.P()
    g.P(`  void setResponseMetaDataProcessor(plugin.IResponseMetaDataProcessor processor){ this.responseMetaDataProcessor = processor; }`)
    g.P(`  plugin.IResponseMetaDataProcessor getResponseMetaDataProcessor(){ return this.responseMetaDataProcessor; }`)
    g.P()
    g.P(`  private <R> R call(String path, com.google.protobuf.Message req, javax.ws.rs.core.MultivaluedMap<String, Object> headers, Class<R> responseClass) {`)
    g.P(`    headers.add("Accept", "application/protobuf");`)
    g.P(`    javax.ws.rs.core.Response response = target.path(path)`)
    g.P(`        .request()`)
    g.P(`        .headers(headers)`)
    g.P(`        .post(javax.ws.rs.client.Entity.entity(req, "application/protobuf"));`)
    g.P(`    if (response.getStatusInfo().getFamily() == javax.ws.rs.core.Response.Status.Family.SUCCESSFUL) {`)
    g.P(`      R r = response.readEntity(responseClass);`)
    g.P(`      response.close();`)
    g.P(`      if(responseMetaDataProcessor != null){ responseMetaDataProcessor.process(response); }`)
    g.P(`      return r;`)
    g.P(`    } else {`)
    g.P(`      throw new javax.ws.rs.WebApplicationException(response);`)
    g.P(`    }`)
    g.P(`  }`)
    g.P()
    g.P(`  private <R> R call(String path, com.google.protobuf.Message req, Class<R> responseClass) {`)
    g.P(`    javax.ws.rs.core.Response response = target.path(path)`)
    g.P(`        .request("application/protobuf")`)
    g.P(`        .post(javax.ws.rs.client.Entity.entity(req, "application/protobuf"));`)
    g.P(`    if (response.getStatusInfo().getFamily() == javax.ws.rs.core.Response.Status.Family.SUCCESSFUL) {`)
    g.P(`      R r = response.readEntity(responseClass);`)
    g.P(`      response.close();`)
    g.P(`      if(responseMetaDataProcessor != null){ responseMetaDataProcessor.process(response); }`)
    g.P(`      return r;`)
    g.P(`    } else {`)
    g.P(`      throw new javax.ws.rs.WebApplicationException(response);`)
    g.P(`    }`)
    g.P(`  }`)
    g.P()
    g.P(`  private <R> java.util.concurrent.Future<R> call(String path, com.google.protobuf.Message request, Class<R> responseClass, int retries) {`)
    g.P(`    if(retries <= 0){throw new IllegalArgumentException("Retries count can't be less than or equal to 0");}`)
    g.P(`    final javax.ws.rs.client.AsyncInvoker invoker = this.target.path(path).request("application/protobuf").async();`)
    g.P(`    final javax.ws.rs.client.Entity<com.google.protobuf.Message> entity = javax.ws.rs.client.Entity.entity(request, "application/protobuf");`)
    g.P(`    java.util.concurrent.CompletableFuture<R> future = new java.util.concurrent.CompletableFuture<>();`)
    g.P(`    invoker.post(entity, new javax.ws.rs.client.InvocationCallback<javax.ws.rs.core.Response>() {`)
    g.P(`      private final java.util.concurrent.atomic.AtomicInteger count = new java.util.concurrent.atomic.AtomicInteger(0);`)
    g.P()
    g.P(`      @Override`)
    g.P(`      public void completed(javax.ws.rs.core.Response response) {`)
    g.P(`        if (response.getStatusInfo().getFamily() == javax.ws.rs.core.Response.Status.Family.SUCCESSFUL) {`)
    g.P(`          future.complete(response.readEntity(responseClass));`)
    g.P(`          response.close();`)
    g.P(`          if(responseMetaDataProcessor != null){ responseMetaDataProcessor.process(response); }`)
    g.P(`        } else {`)
    g.P(`          failed(new javax.ws.rs.WebApplicationException(response));`)
    g.P(`        }`)
    g.P(`      }`)
    g.P()
    g.P(`      @Override`)
    g.P(`      public void failed(Throwable throwable) {`)
    g.P(`        if (count.getAndIncrement() < retries) {`)
    g.P(`          try {`)
    g.P(`            Thread.sleep(500);`)
    g.P(`          } catch (InterruptedException e) {`)
    g.P(`            future.completeExceptionally(e);`)
    g.P(`          }`)
    g.P(`          invoker.post(entity, this);`)
    g.P(`        } else {`)
    g.P(`          future.completeExceptionally(throwable);`)
    g.P(`        }`)
    g.P(`      }`)
    g.P(`    });`)
    g.P(`    return future;`)
    g.P(`  }`)
    g.P()
    

    for _, method := range service.Method {
        inputType := g.GetType(file, method.GetInputType())
        outputType := g.GetType(file, method.GetOutputType())
        methodName := lowerCamelCase(method.GetName())
        methodPath := camelCase(method.GetName())

        g.P()
        g.P(`// Performing a request with headers`)
        g.P(`  public `, outputType, ` `, methodName, `(`, inputType, ` request, javax.ws.rs.core.MultivaluedMap<String, Object> headers) {`)
        g.P(`    return call("/`, methodPath, `", request, headers, `, outputType, `.class);`)
        g.P(`  }`)
        g.P()
        // add comment
        g.P(`  @Override`)
        g.P(`  public `, outputType, ` `, methodName, `(`, inputType, ` request) {`)
        g.P(`    return call("/`, methodPath, `", request, `, outputType, `.class);`)
        g.P(`  }`)
        g.P()
        // add comment
        g.P(`  @Override`)
        g.P(`  public java.util.concurrent.Future<`, outputType, `> `, methodName, `(`, inputType, ` request, int retries) {`)
        g.P(`    return call("/`, methodPath, `", request, `, outputType, `.class, retries);`)
        g.P(`  }`)
    }

    g.P(`}`)
    g.P()

    out := &pluginpb.CodeGeneratorResponse_File{}
    out.Content = proto.String(g.output.String())
    if multi {
        out.Name = proto.String(getJavaServiceClientClassFile(file, service))
    } else {
        out.Name = proto.String(getJavaOuterClassFile(file))
        out.InsertionPoint = proto.String("outer_class_scope")
    }
    g.Reset()

    return out
}

func (g *generator) generateServiceInterface(file *descriptorpb.FileDescriptorProto, service *descriptorpb.ServiceDescriptorProto) *pluginpb.CodeGeneratorResponse_File {
    // TODO add comment

    serviceClass := getJavaServiceClassName(file, service)
    servicePath := g.GetServicePath(file, service)
    multi := file.Options.GetJavaMultipleFiles()

    if multi {
        pkg := getJavaPackage(file)
        g.P(`// Code generated by protoc-gen-twirp_java_jaxrs, DO NOT EDIT.`)
        g.P(`// source: `, file.GetName())
        g.P()
        if pkg != "" {
            g.P(`package `, pkg, `;`)
            g.P()
        }
    }

    g.P(`@javax.ws.rs.Path( "/`, servicePath, `" )`)
    g.P(`public interface `, serviceClass, ` {`)

    for _, method := range service.Method {
        inputType := g.GetType(file, method.GetInputType())
        outputType := g.GetType(file, method.GetOutputType())
        methodName := lowerCamelCase(method.GetName())

        // add comment
        g.P(`  `, `@javax.ws.rs.POST`)
        g.P(`  `, `@javax.ws.rs.Path( "/`, strings.Title(methodName), `" )`)
        g.P(`  `, `@javax.ws.rs.Consumes({"application/protobuf", "application/json"})`)
        g.P(`  `, `@javax.ws.rs.Produces({"application/protobuf", "application/json"})`)
        g.P(`  `, outputType, ` `, methodName, `(`, inputType, ` request);`)
        g.P()
        // add comment
        g.P(`  `, `@javax.ws.rs.POST`)
        g.P(`  `, `@javax.ws.rs.Path( "/`, strings.Title(methodName), `" )`)
        g.P(`  `, `@javax.ws.rs.Consumes({"application/protobuf", "application/json"})`)
        g.P(`  `, `@javax.ws.rs.Produces({"application/protobuf", "application/json"})`)
        g.P(`  java.util.concurrent.Future<`, outputType, `> `, methodName, `(`, inputType, ` request, int retries);`)
    }

    g.P(`}`)
    g.P()

    out := &pluginpb.CodeGeneratorResponse_File{}
    out.Content = proto.String(g.output.String())

    if multi {
        out.Name = proto.String(getJavaServiceClassFile(file, service))
    } else {
        out.Name = proto.String(getJavaOuterClassFile(file))
        out.InsertionPoint = proto.String("outer_class_scope")
    }
    g.Reset()

    return out
}

func (g *generator) Reset() {
    g.indent = ""
    g.output.Reset()
}

func (g *generator) In() {
    g.indent += "  "
}

func (g *generator) Out() {
    g.indent = g.indent[2:]
}

func (g *generator) P(str ...string) {
    for _, v := range str {
        g.output.WriteString(v)
    }
    g.output.WriteByte('\n')
}

func (g *generator) getProtoFiles() []*descriptorpb.FileDescriptorProto {
    files := make([]*descriptorpb.FileDescriptorProto, 0)
    for _, fname := range g.Request.GetFileToGenerate() {
        for _, _proto := range g.Request.GetProtoFile() {
            if _proto.GetName() == fname {
                files = append(files, _proto)
            }
        }
    }
    return files
}

func (g *generator) GetServicePath(file *descriptorpb.FileDescriptorProto, service *descriptorpb.ServiceDescriptorProto) string {
    name := camelCase(service.GetName())
    pkg := file.GetPackage()
    if pkg != "" {
        name = pkg + "." + name
    }
    return name
}

func (g *generator) GetType(file *descriptorpb.FileDescriptorProto, name string) string {

    multi := file.Options.GetJavaMultipleFiles()
    if name[0:1] == "." {
        name = name[1:]
    }

    path, pkg, class, dot := "", "", "", strings.LastIndex(name, ".")
    if dot > -1 {
        slice := strings.Split(name, ".")
        pkg, class = slice[0], slice[1]
    } else {
        pkg, class = "", name
    }

    if containsType(class, file){
        path = getJavaPackage(file)
    } else {
        for _, dep := range g.Request.GetProtoFile() {
            if dep.GetPackage() == pkg && containsType(class, dep){
                path = getJavaPackage(dep)
                break
            }
        }
    }

    if multi {
        return fmt.Sprintf("%s.%s", path, class)
    } else {
        return fmt.Sprintf("%s.%s.%s", path, getJavaOuterClassName(file), class)
    }
}
