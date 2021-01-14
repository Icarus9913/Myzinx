package main

import (
	"Myzinx/protobufDemo/pb"
	"fmt"
	"github.com/golang/protobuf/proto"
)

func main()  {
	//定义一个Person结构对象
	person := &pb.Person{
		Name: "wk",
		Age: 16,
		Emails: []string{"110@qq.com","120.gmail.com","119@163.com"},
		Phones: []*pb.PhoneNumber{
			&pb.PhoneNumber{
				Number: "111111",
				Type: pb.PhoneType_MOBILE,
			},
			&pb.PhoneNumber{
				Number: "222222",
				Type: pb.PhoneType_HOME,
			},
			&pb.PhoneNumber{
				Number: "333333",
				Type: pb.PhoneType_WORK,
			},
		},

	}

	//编码
	//将person对象，就是将protobuf的message进行序列化，得到一个二进制文件
	data, err := proto.Marshal(person)
	//data就是我们要进行网络传输的数据，对端需要按照Message Person格式进行解析
	if nil!=err{
		fmt.Println("marshal err:",err)
	}

	//解码
	newPerson := &pb.Person{}
	err = proto.Unmarshal(data, newPerson)
	if nil!=err{
		fmt.Println("unmarshal err:",err)
	}
	fmt.Println("原数据:",person)
	fmt.Println("解码之后的数据:",newPerson)
}
