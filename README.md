Twirp generator for Java/JAX-RS Client and Server
=================================================

Build
-----

    go get github.com/fajran/protoc-gen-twirp_java_jaxrs


Usage
-----

    export PATH=$PATH:$GOPATH/bin
    protoc service.proto --java_out=src/main/java --twirp_java_jaxrs_out=src/main/java

