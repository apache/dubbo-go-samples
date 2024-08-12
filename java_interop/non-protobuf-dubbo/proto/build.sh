protoc -I ./ \
  --go-hessian2_out=./ --go-hessian2_opt=paths=source_relative \
  --go-dubbo_out=./  --go-dubbo_opt=paths=source_relative \
  ./greet.proto
