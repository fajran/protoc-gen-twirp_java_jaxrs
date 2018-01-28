Twirp generator for Java/JAX-RS Client and Server
=================================================

**ATTENTION** This tool is still actively developed. There is no guarantee that
the generated code, interface, or even the way they are supposed to be used are
stable.

Build
-----

    go get github.com/fajran/protoc-gen-twirp_java_jaxrs


Usage
-----

    export PATH=$PATH:$GOPATH/bin
    protoc service.proto --java_out=src/main/java --twirp_java_jaxrs_out=src/main/java

