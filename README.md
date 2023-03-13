Twirp generator for Java/JAX-RS Client
======================================

This is a protobuf generator that creates a Java/JAX-RS client for
[Twirp](https://github.com/twitchtv/twirp).

Server support is planned #3.

**ATTENTION** This tool is still actively developed. There is no guarantee that
the generated code, interface, or even the way they are supposed to be used are
stable.

Build
-----

    go get github.com/nutshelllabs/protoc-gen-twirp_java_jaxrs


Usage
-----

    export PATH=$PATH:$GOPATH/bin
    protoc service.proto --java_out=src/main/java --twirp_java_jaxrs_out=src/main/java

