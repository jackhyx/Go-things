package go源码

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
)
/*
1.5、Go-Micro特性
Registry：主要负责服务注册和发现功能。我们之前学习过的consul，就可以和此处的Registry结合起来，实现服务的发现功能。

Selector：selector主要的作用是实现服务的负载均衡功能。当某个客户端发起请求时，将首先查询服务注册表，返回当前系统中可用的服务列表，然后从中选择其中一个节点进行查询，保证节点可用。

Broker：Broker是go-micro框架中事件发布和订阅的接口，主要是用消息队列的方式实现信息的接收和发布，用于处理系统间的异步功能。

Codec：go-micro中数据传输过程中的编码和解码接口。go-micro中有多重编码方式，默认的实现方式是protobuf，除此之外，还有json等格式。

Transport：go-micro框架中的通信接口，有很多的实现方案可以选择，默认使用的是http形式的通信方式，除此以外，还有grpc等通信方式。

Client和Server：分别是go-micro中的客户端接口和服务端接口。client负责调用，server负责等待请求处理。
 */


/* 服务的定义:在micro框架中，服务用接口来进行定义，服务被定义为Service，完整的接口定义如下
在该接口中，定义了一个服务实例具体要包含的方法
分别是：Init、Options、Client、Server、Run、String等6个方法。 */

type Service interface {
	Init(...Option)
	Options() Options
	Client() client.Client
	Server() server.Server
	Run() error
	String() string
}
// 初始化服务实例: micro框架，除了提供Service的定义外，提供创建服务实例的方法供开发者调用：

service := micro.NewService()

// 如上是最简单一种创建service实例的方式。NewService可以接受一个Options类型的可选项参数。NewService的定义如下：

func NewService(opts ...Option) Service {
	return newService(opts...)
}

/* Options可选项配置
关于Options可配置选项，有很多可以选择的设置项。
micro框架包中包含了options.go文件，定义了详细的可选项配置的内容。
最基本常见的配置项有：服务名称，服务的版本，服务的地址，服务： */

//服务名称
func Name(n string) Option {
	return func(o *Options) {
		o.Server.Init(server.Name(n))
	}
}

//服务版本
func Version(v string) Option {
	return func(o *Options) {
		o.Server.Init(server.Version(v))
	}
}

//服务部署地址
func Address(addr string) Option {
	return func(o *Options) {
		o.Server.Init(server.Address(addr))
	}
}

//元数据项设置
func Metadata(md map[string]string) Option {
	return func(o *Options) {
		o.Server.Init(server.Metadata(md))
	}
}

// 完整的实例化对象代码如下所示：

func main() {
	//创建一个新的服务对象实例
	service := micro.NewService(
		micro.Name("helloservice"),
		micro.Version("v1.0.0"),
	)
}

// 开发者可以直接调用micro.Name为服务设置名称，设置版本号等信息。
// 在对应的函数内部，调用了server.Server.Init函数对配置项进行初始化。

/* 定义服务接口,实现服务业务逻辑
使用protobuf定义服务接口并自动生成go语言文件,需要经过以下几个步骤
 */

// 定义.proto文件:使用proto3语法定义数据结构体和服务方法。具体定义内容如下
syntax = 'proto3';
package message;

//学生数据体
message Student {
string name = 1; //姓名
string classes = 2; //班级
int32 grade = 3; //分数
}

//请求数据体定义
message StudentRequest {
string name = 1;
}

//学生服务
service StudentService {
//查询学生信息服务
rpc GetStudent (StudentRequest) returns (Student);
}

/*
编译.proto文件:在原来学习gRPC框架时，我们是将.proto文件按照grpc插件的标准来进行编译。
而现在，我们学习的是go-micro，因此我们可以按照micro插件来进行编译。
micro框架中的protobuf插件，我们需要单独安装;
 */

// 编码实现服务功能 :在项目目录下，实现StudentService定义的rpc GetStudent功能。
// 新建studentManager.go文件，具体实现如下：

//学生服务管理实现
type StudentManager struct {
}

//获取学生信息的服务接口实现
func GetStudent(ctx context.Context, request *message.StudentRequest, response *message.Student) error {

	studentMap := map[string]message.Student{
		"davie":  message.Student{Name: "davie", Classes: "软件工程专业", Grade: 80},
		"steven": message.Student{Name: "steven", Classes: "计算机科学与技术", Grade: 90},
		"tony":   message.Student{Name: "tony", Classes: "计算机网络工程", Grade: 85},
		"jack":   message.Student{Name: "jack", Classes: "工商管理", Grade: 96},
	}

	if request.Name == "" {
		return errors.New(" 请求参数错误,请重新请求。")
	}

	student := studentMap[request.Name]

	if student.Name != "" {
		response = &student
	}
	return errors.New(" 未查询当相关学生信息 ")
}
// 运行服务: 我们用micro框架来实现服务的运行。完整的运行服务的代码如下：

func main() {

	//创建一个新的服务对象实例
	service := micro.NewService(
		micro.Name("student_service"),
		micro.Version("v1.0.0"),
	)

	//服务初始化
	service.Init()

	//注册
message.RegisterStudentServiceHandler(service.Server(), new(StudentManager))

	//运行
	err := service.Run()
	if err != nil {
		log.Fatal(err)
	}
}
// 客户端调用: 客户端可以构造请求对象，并访问对应的服务方法。具体方法实现如下：

func main() {

	service := micro.NewService(
		micro.Name("student.client"),
	)
	service.Init()

	studentService := message.NewStudentServiceClient("student_service", service.Client())

	res, err := studentService.GetStudent(context.TODO(), &message.StudentRequest{Name: "davie"})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res.Name)
	fmt.Println(res.Classes)
	fmt.Println(res.Grade)
}
// 注册服务到consul
// 默认注册到mdns
// 在我们运行服务端的程序时，我们可以看到Registry [mdns] Registering node:xxx这个日志,该日志显示go-micro框架将我们的服务使用默认的配置注册到了mdns中。
// mdns是可简单翻译为mdns，是go-micro的默认配置选项。


/*
事件发布: 只有当用户服务中的某个功能执行时，才会触发相应的事件，并将对应的用户数据等消息发送到消息队列组件中，这个过程我们称之为事件发布。

事件订阅:与事件发布对应的是事件订阅。
我们增加消息队列组件的目的是实现模块程序的解耦，原来是程序调用端主动进行程序调用;
现在需要由另外一方模块的程序到消息队列组件中主动获取需要相关数据并进行相关功能调用。这个主动获取的过程称之为订阅。

基于消息发布/订阅的消息系统有很多种框架的实现，常见的有：Kafka、RabbitMQ、ActiveMQ、Kestrel、NSQ等。

MQTT全称是Message Queuing Telemetry Transport，翻译为消息队列遥测传输协议，是一种基于发布/订阅模式的"轻量级"的通讯协议，该协议基于TCP/IP协议，由IBM在1999年发布;
MQTT的最大优点在于，可以用极少的代码和有限的宽带,为连接远程设备提供提供实时可靠的消息服务。

 */

/*
 编程实现
接下来进行订阅和发布机制的编程的实现。
消息组件初始化
如果要想使用消息组件完成消息的发布和订阅，首先应该让消息组件正常工作。因此，需要先对消息组件进行初始化。我们可以在服务创建时，对消息组件进行初始化，并进行可选项配置,设置使用mqtt作为消息组件。代码实现如下：
*/
...
server := micro.NewService(
		micro.Name("go.micro.srv"),
		micro.Version("latest"),
		micro.Broker(mqtt.NewBroker()),
)
...
// 可以使用micro.Broker来指定特定的消息组件，并通过mqtt.NewBroker初始化一个mqtt实例对象,作为broker参数。

// 消息订阅:因为是时间驱动机制，消息的发送方随时可能发布相关事件。因此需要消息的接收方先进行订阅操作，避免遗漏消息。
// go-micro框架中可以通过broker.Subscribe实现消息订阅。编程代码如下所示：

...
pubSub := service.Server().Options().Broker
_, err := pubSub.Subscribe("go.micro.srv.message", func(event broker.Event) error {
	var req *message.StudentRequest
	if err := json.Unmarshal(event.Message().Body, &req); err != nil {
		return err
	}
	fmt.Println(" 接收到信息：", req)
	//去执行其他操作

	return nil
})
...
// 消息发布:完成了消息的订阅，我们再来实现消息的发布。在客户端实现消息的发布。
// 在go-micro框架中，可以使用broker.Publish来进行消息的发布,具体的代码如下：

...

brok := service.Server().Options().Broker
if err := brok.Connect(); err != nil {
log.Fatal(" broker connection failed, error : ", err.Error())
}

student := &message.Student{Name: "davie", Classes: "软件工程专业", Grade: 80, Phone: "12345678901"}
msgBody, err := json.Marshal(student)
if err != nil {
log.Fatal(err.Error())
}
msg := &broker.Message{
Header: map[string]string{
"name": student.Name,
},
Body: msgBody,
}

err = brok.Publish("go.micro.srv.message", msg)
if err != nil {
log.Fatal(" 消息发布失败：%s\n", err.Error())
} else {
log.Print("消息发布成功")
}

...

// 弊端:在服务端通过fmt.println日志，可以输出event.Message().Body)数据，其格式为：

{"name":"davie","classes":"软件工程专业","grade":80,"phone":"12345678901"}


/* 我们可以看到在服务实例之间传输的数据格式是json格式。
根据之前学习proto知识可以知道，在进行消息通信时，采用JSON格式进行数据传输，其效率比较低。
因此，这意味着，当我们在使用第三方消息组件进行消息发布/订阅时，会失去对protobuf的使用。
这对追求高消息的开发者而言，是需要解决和改进的问题。
因为使用protobuf可以直接在多个服务之间使用二进制流数据进行传输，要比json格式高效的多。
 */
// Micro负载均衡组件--Selector：
// 所谓负载均衡，英文为Load Balance，其意思是将负载进行平衡、分摊到多个操作单元上进行执行。

/*
负载均衡器主要处理四种请求，分别是：HTTP、HTTPS、TCP、UDP。

负载均衡算法
负载均衡器的作用既然是负责接收请求，并实现请求的分发，因此需要按照一定的规则进行转发处理。
负载均衡器可以按照不同的规则实现请求的转发，其遵循的转发规则称之为负载均衡算法。常用的负载均衡算法有以下几个：

Round Robin（轮训算法）：所谓轮训算法，其含义很简单，就是按照一定的顺序进行依次排队分发。
当有请求队列需要转发时，为第一个请求选择可用服务列表中的第一个服务器，为下一个请求选择服务列表中的第二个服务器。
按照此规则依次向下进行选择分发，直到选择到服务器列表的最后一个。当第一次列表转发完毕后，重新选择第一个服务器进行分发，此为轮训。

Least Connections（最小连接）：因为分布式系统中有多台服务器程序在运行，每台服务器在某一个时刻处理的连接请求数量是不一样的。
因此，当有新的请求需要转发时，按照最小连接数原则，负载均衡器会有限选择当前连接数最小的服务器，以此来作为转发的规则。

Source（源）：还有一种常见的方式是将请求的IP进行hash计算，根据结算结果来匹配要转发的服务器，然后进行转发。
这种方式可以一定程度上保证特定用户能够连接到相同的服务器。



 */
/*
Micro的Selector
Selector的英文是选择器的意思，在Micro中实现了Selector组件，运行在客户端实现负载均衡功能。当客户端需要调用服务端方法时，客户端会根据其内部的selector组件中指定的负载均衡策略选择服务注册中中的一个服务实例。Go-micro中的Selector是基于Register模块构建的，提供负载均衡策略，同时还提供过滤、缓存和黑名单等功能。

Selector定义
首先，让我们来看一下Selector的定义：
*/
type Selector interface {
	Init(opts ...Option) error
	Options() Options
	// Select returns a function which should return the next node
	Select(service string, opts ...SelectOption) (Next, error)
	// Mark sets the success/error against a node 标记服务节点的状态
	Mark(service string, node *registry.Node, err error)
	// Reset returns state back to zero for a service
	Reset(service string)
	// Close renders the selector unusable
	Close() error
	// Name of the selector
	String() string
}

// 如上是go-micro框架中的Selector的定义，Selector接口定义中包含Init、Options、Mark、Reset、Close、String方法。
// 其中Select是核心方法，可以实现自定义的负载均衡策略，Mark方法用于标记服务节点的状态,String方法返回自定义负载均衡器的名称。

// DefaultSelector : 在selector包下，除Selector接口定义外，还包含DefaultSelector的定义，作为go-micro默认的负载均衡器而被使用。
// DefaultSelector是通过NewSelector函数创建生成的。NewSelector函数实现如下:

func NewSelector(opts ...Option) Selector {
	sopts := Options{
		Strategy: Random,
	}

	for _, opt := range opts {
		opt(&sopts)
	}

	if sopts.Registry == nil {
		sopts.Registry = registry.DefaultRegistry
	}

	s := &registrySelector{
		so: sopts,
	}
	s.rc = s.newCache()

	return s
}
// 在NewSelector中，实例化了registrySelector对象并进行了返回,在实例化的过程中，配置了Selector的Options选项，默认的配置是Random。我们进一步查看会发现Random是一个func，定义如下：

func Random(services []*registry.Service) Next {
	var nodes []*registry.Node

	for _, service := range services {
		nodes = append(nodes, service.Nodes...)
	}

	return func() (*registry.Node, error) {
		if len(nodes) == 0 {
			return nil, ErrNoneAvailable
		}

		i := rand.Int() % len(nodes)
		return nodes[i], nil
	}
}
// 该算法是go-micro中默认的负载均衡器，会随机选择一个服务节点进行分发；除了Random算法外，还可以看到RoundRobin算法，如下所示：

func RoundRobin(services []*registry.Service) Next {
	var nodes []*registry.Node

	for _, service := range services {
		nodes = append(nodes, service.Nodes...)
	}

	var i = rand.Int()
	var mtx sync.Mutex

	return func() (*registry.Node, error) {
		if len(nodes) == 0 {
			return nil, ErrNoneAvailable
		}

		mtx.Lock()
		node := nodes[i%len(nodes)]
		i++
		mtx.Unlock()
		return node, nil
	}
}
//registrySelector : registrySelector是selector包下default.go文件中的结构体定义，具体定义如下:

type registrySelector struct {
	so Options
	rc cache.Cache
}
// 缓存Cache:目前已经有了负载均衡器，我们可以看到在Selector的定义中，还包含一个cache.Cache结构体类型，这是什么作用呢？
// 有了Selector以后，我们每次请求负载均衡器都要去Register组件中查询一次，这样无形之中就增加了成本，降低了效率，没有办法达到高可用。
// 为了解决以上这种问题，在设计Selector的时候设计一个缓存，Selector将自己查询到的服务列表数据缓存到本地Cache中。
// 当需要处理转发时，先到缓存中查找，如果能找到即分发；如果缓存当中没有，会执行请求服务发现注册组件，然后缓存到本地。 具体的实现机制如下所示：

type Cache interface {
	// embed the registry interface
	registry.Registry
	// stop the cache watcher
	Stop()
}

func (c *cache) watch(w registry.Watcher) error {
	// used to stop the watch
	stop := make(chan bool)

	// manage this loop
	go func() {
		defer w.Stop()

		select {
		// wait for exit
		case <-c.exit:
			return
		// we've been stopped
		case <-stop:
			return
		}
	}()

	for {
		res, err := w.Next()
		if err != nil {
			close(stop)
			return err
		}
		c.update(res)
	}
}
// 通过watch实现缓存的更新、创建、移除等操作。

// 黑名单 : 在了解完了缓存后，我们再看看Selector中其他的方法。在Selector接口的定义中，还可以看到有Mark和Reset的声明。具体声明如下：

// Mark sets the success/error against a node
Mark(service string, node *registry.Node, err error)
// Reset returns state back to zero for a service
Reset(service string)

/* Mark方法可以用于标记服务注册和发现组件中的某一个节点的状态，这是因为在某些情况下，负载均衡器跟踪请求的执行情况。
如果请求被转发到某天服务节点上，多次执行失败，就意味着该节点状态不正常，此时可以通过Mark方法设置节点变成黑名单;
以过滤掉掉状态不正常的节点。

 */

/*
Go-Micro API网关
Micro框架中有API网关的功能。
API网关的作用是为微服务做代理，负责将微服务的RPC方法代理成支持HTTP协议的web请求，同时将用户端使用的URL进行暴露。
 */

// 服务定义和编译
// 定义学生消息体proto文件：

syntax = 'proto3';

package proto;

message Student {
string id = 1;
string name = 2;
int32 grade = 3;
string classes = 4;
}

message Request {
string name = 1;
}

service StudentService {
rpc GetStudent (Request) returns (Student);
}
// 在proto文件中定义了Student、Request消息体和rpc服务。使用micro api网关功能，编译proto文件，需要生成micro文件。编译生成该文件需要使用到一个新的protoc-gen-micro库，安装protoc-gen-micro库命令如下：

go get github.com/micro/protoc-gen-micro
// 再次编译proto文件，需要指定两个参数，分别是：go_out和micro_out，详细命令如下：

protoc --go_out=. --micro_out=. student.proto
// 上述命令执行成功后，会自动生成两个go语言文件：student.pb.go和student.micro.go。

// micro.go文件中生成的内容包含服务的实例化，和相应的服务方法的底层实现。

// 服务端实现: 我们都知道正常的Web服务，是通过路由处理http的请求的。
// 在此处也是一样的，我们可以通过路由处理来解析HTTP请求的接口，service对象中包含路由处理方法。详细代码如下所示：

...
type StudentServiceImpl struct {
}

//服务实现
func (ss *StudentServiceImpl) GetStudent(ctx context.Context, request *proto.Request, resp *proto.Student) error {

	//tom
	studentMap := map[string]proto.Student{
		"davie":  proto.Student{Name: "davie", Classes: "软件工程专业", Grade: 80},
		"steven": proto.Student{Name: "steven", Classes: "计算机科学与技术", Grade: 90},
		"tony":   proto.Student{Name: "tony", Classes: "计算机网络工程", Grade: 85},
		"jack":   proto.Student{Name: "jack", Classes: "工商管理", Grade: 96},
	}

	if request.Name == "" {
		return errors.New(" 请求参数错误,请重新请求。")
	}

	//获取对应的student
	student := studentMap[request.Name]
	if student.Name != "" {
		fmt.Println(student.Name, student.Classes, student.Grade)
		*resp = student
		return nil
	}
	return errors.New(" 未查询当相关学生信息 ")
}

func main() {
	service := micro.NewService(
		micro.Name("go.micro.srv.student"),
	)

	service.Init()
	proto.RegisterStudentServiceHandler(service.Server(), new(StudentServiceImpl))

	if err := service.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
...

// server程序进行服务的实现和服务的运行。
// REST 映射 : 现在，RPC服务已经编写完成。我们需要编程实现API的代理功能，用于处理HTTP形式的请求。
// 在rest.go文件中，实现rest的映射，详细代码如下：

type Student struct {
}

var (
	cli proto.StudentService
)

func (s *Student) GetStudent(req *restful.Request, rsp *restful.Response) {

	name := req.PathParameter("name")
	fmt.Println(name)
	response, err := cli.GetStudent(context.TODO(), &proto.Request{
		Name: name,
	})

	if err != nil {
		fmt.Println(err.Error())
		rsp.WriteError(500, err)
	}

	rsp.WriteEntity(response)
}

func main() {

	service := web.NewService(
		web.Name("go.micro.api.student"),
	)

	service.Init()

	cli = proto.NewStudentService("go.micro.srv.student", client.DefaultClient)

	student := new(Student)
	ws := new(restful.WebService)
	ws.Path("/student")
	ws.Consumes(restful.MIME_XML, restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON, restful.MIME_XML)

	ws.Route(ws.GET("/{name}").To(student.GetStudent))

	wc := restful.NewContainer()
	wc.Add(ws)

	service.Handle("/", wc)

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}