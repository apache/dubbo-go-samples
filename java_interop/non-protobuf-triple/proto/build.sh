protoc -I ./ \
  --go-hessian2_out=./ --go-hessian2_opt=paths=source_relative \
  --go-triple_out=./  --go-triple_opt=paths=source_relative \
  ./greet.proto
