mkdir -p generated

yarn pbjs \
  --target static-module \
  --no-encode \
  --no-decode \
  --path ../../ \
  --out ./generated/index.js \
  ./protos/HelloService.proto \

yarn pbts \
  --out ./generated/index.d.ts \
  ./generated/index.js \

### https://github.com/protobufjs/protobuf.js/issues/1222
sed -i -e "s/\[ 'Promise' \]\./Promise/g" "generated/index.d.ts"
sed -i -e "s/\[ 'object' \]\.<string, any>/{ [k: string]: any }/g" "generated/index.d.ts"
sed -i -e "s/\[ 'Array' \]\./Array/g" "generated/index.d.ts"
###