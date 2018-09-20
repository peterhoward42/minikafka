To compile toykafka.proto into toykafka.pb.go use the embedded *go generate* 
directive in generate.go.

I.e.

    cd .
    go generate

Note generate.go exists ONLY as a place for this go generate directive 
comment to live in.