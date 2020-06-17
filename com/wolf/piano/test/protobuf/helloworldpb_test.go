package protobuf

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"log"
	"testing"
)

func TestProtobufBase(t *testing.T) {

	data_encode := &Helloworld{
		Id:  proto.Int32(11),
		Str: proto.String("hello world!"),
		Opt: proto.Int32(17),
	}

	data, err := proto.Marshal(data_encode)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}

	data_decode := &Helloworld{}
	err = proto.Unmarshal(data, data_decode)
	if err != nil {
		log.Fatal("unmarshaling error: ", err)
	}

	if data_encode.GetId() != data_decode.GetId() {
		log.Fatalf("data mismatch %q != %q", data_encode.GetId(), data_decode.GetId())
	}
	fmt.Println("ID:", data_decode.GetId())
	fmt.Println("Str:", data_decode.GetStr())
	fmt.Println("Opt:", data_decode.GetOpt())
}
