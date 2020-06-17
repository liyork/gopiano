path=$(dirname $0)
#echo $path
#echo $(pwd)
path=${path/\./$(pwd)}
echo $path

# 编译Protobuf协议

protoc --version

protoc --go_out=$path/../ -I=$path $path/helloworld.proto

# 最后正确的编写以及编译方式：
# 工程用mod，以工程的根为基准进行相互import，添加go_package让生成最后的go文件以这个进行被依赖
# 当protoc编译时，在根开始，然后--proto_path=.即可，生成目录--go_out=paths=source_relative:.添加source_relative就是为了在相同
# 目录生成，也可另行指定