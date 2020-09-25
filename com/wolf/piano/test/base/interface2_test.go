package base

import (
	"fmt"
	"testing"
)

// interface 是GO语言的基础特性之一。可以理解为一种类型的规范或者约定
//
//interface 是方法声明的集合
//任何类型的对象实现了在interface 接口中声明的全部方法，则表明该类型实现了该接口。
//interface 可以作为一种数据类型，实现了该接口的任何对象都可以给对应的接口类型变量赋值。
//a. interface 可以被任意对象实现，一个类型/对象也可以实现多个 interface
//b. 方法不能重载，如 eat(), eat(s string) 不能同时存在

type Phone interface {
	call()
}

type NokiaPhone struct {
}

func (nokiaPhone NokiaPhone) call() {
	fmt.Println("I am Nokia, I can call you!")
}

type ApplePhone struct {
}

func (iPhone ApplePhone) call() {
	fmt.Println("I am Apple Phone, I can call you!")
}

// 基本使用接口，承接不同实现
func TestInterfaceBase2(t *testing.T) {
	var phone Phone
	phone = new(NokiaPhone)
	phone.call()

	phone = new(ApplePhone)
	phone.call()
}

// 面向对象中的开闭原则
// 作为interface数据类型，他存在的意义在哪呢？ 实际上是为了满足一些面向对象的编程思想。
// 软件设计的最高目标就是高内聚，低耦合。
// 随着我们的程序写的越来越复杂，银行业务员可能就要增加方法，会导致业务员模块越来越臃肿。
// 这样的设计(一个Banker类包含所有方法)会导致，当我们去给Banker添加新的业务的时候，会直接修改原有的Banker代码，那么Banker模块的功能会越来越多，
// 出现问题的几率也就越来越大，所有的业务都在一个Banker类里，他们的耦合度太高，Banker的职责也不够单一，代码的维护成本随着业务的复杂正比成倍增大。
// 开闭原则定义:
// 一个软件实体如类、模块和函数应该对扩展开放，对修改关闭。
// 简单的说就是在修改需求的时候，应该尽量通过扩展来实现变化，而不是通过修改已有代码来实现变化。

// 如果我们拥有接口, interface这个东西，那么我们就可以抽象一层出来，制作一个抽象的Banker模块，然后提供一个抽象的方法。 分别根据这个抽象模块，去实现支付Banker（实现支付方法）,转账Banker（实现转账方法）
// 依然可以搞定程序的需求。
// 当我们想要给Banker添加额外功能的时候，之前我们是直接修改Banker的内容，现在我们可以单独定义一个股票Banker(实现股票方法)，到这个系统中。 而且股票Banker的实现成功或者失败都不会影响之前的稳定系统，他很单一，而且独立。
// 当我们给一个系统添加一个功能的时候，不是通过修改代码，而是通过增添代码来完成，那么就是开闭原则的核心思想了。所以要想满足上面的要求，是一定需要interface来提供一层抽象的接口的。
//抽象的银行业务员
type AbstractBanker interface {
	DoBusi() //抽象的处理业务接口
}

//存款的业务员
type SaveBanker struct {
	//AbstractBanker
}

func (sb *SaveBanker) DoBusi() {
	fmt.Println("进行了存款")
}

//转账的业务员
type TransferBanker struct {
	//AbstractBanker
}

func (tb *TransferBanker) DoBusi() {
	fmt.Println("进行了转账")
}

//支付的业务员
type PayBanker struct {
	//AbstractBanker
}

func (pb *PayBanker) DoBusi() {
	fmt.Println("进行了支付")
}

func TestInterfaceOpenClose(t *testing.T) {
	//进行存款
	sb := &SaveBanker{}
	sb.DoBusi()

	//进行转账
	tb := &TransferBanker{}
	tb.DoBusi()

	//进行支付
	pb := &PayBanker{}
	pb.DoBusi()
}

// 上面有些组织逻辑需要我们进行手动组织处理
// 也可以根据AbstractBanker设计一个小框架
//实现架构层(基于抽象层进行业务封装-针对interface接口进行封装)
func BankerBusiness(banker AbstractBanker) {
	//通过接口来向下调用，(多态现象)
	banker.DoBusi()
}

func TestInterfaceOpenCloseAbstract(t *testing.T) {
	//进行存款
	BankerBusiness(&SaveBanker{})

	//进行存款
	BankerBusiness(&TransferBanker{})

	//进行存款
	BankerBusiness(&PayBanker{})
}

// 接口的最大的意义就是实现多态的思想，就是我们可以根据interface类型来设计API接口，那么这种API接口的适应能力不仅能适应当下所实现的全部模块，也适应未来实现的模块来进行调用。 调用未来可能就是接口的最大意义所在吧，

// 面向对象中的依赖倒转原则
//耦合度极高的模块关系设计

// === > 奔驰汽车 <===
type Benz struct {
}

func (this *Benz) Run() {
	fmt.Println("Benz is running...")
}

// === > 宝马汽车  <===
type BMW struct {
}

func (this *BMW) Run() {
	fmt.Println("BMW is running ...")
}

//===> 司机张三  <===
type Zhang3 struct {
	//...
}

func (zhang3 *Zhang3) DriveBenZ(benz *Benz) {
	fmt.Println("zhang3 Drive Benz")
	benz.Run()
}

func (zhang3 *Zhang3) DriveBMW(bmw *BMW) {
	fmt.Println("zhang3 drive BMW")
	bmw.Run()
}

//===> 司机李四 <===
type Li4 struct {
	//...
}

func (li4 *Li4) DriveBenZ(benz *Benz) {
	fmt.Println("li4 Drive Benz")
	benz.Run()
}

func (li4 *Li4) DriveBMW(bmw *BMW) {
	fmt.Println("li4 drive BMW")
	bmw.Run()
}

// 代码中每个模块之间的依赖关系，实际上并没有用到任何的interface接口层的代码，显然最后我们的两个业务:张三开奔驰, 李四开宝马，程序中也都实现了。
// 但是这种设计的问题就在于，小规模没什么问题，但是一旦程序需要扩展，比如我现在要增加一个丰田汽车 或者 司机王五， 那么模块和模块的依赖关系将成指数级递增，想蜘蛛网一样越来越难维护和捋顺。
func TestCoupling(t *testing.T) {
	//业务1 张3开奔驰
	benz := &Benz{}
	zhang3 := &Zhang3{}
	zhang3.DriveBenZ(benz)

	//业务2 李四开宝马
	bmw := &BMW{}
	li4 := &Li4{}
	li4.DriveBMW(bmw)
}

// 面向抽象层依赖倒转
// 如果我们在设计一个系统的时候，将模块分为3个层次，业务逻辑层-->抽象层<--实现层、。那么，我们首先将抽象层的模块和接口定义出来，这里就需要了interface接口的设计，
// 然后我们依照抽象层，依次实现每个实现层的模块，
// 在我们写实现层代码的时候，实际上我们只需要参考对应的抽象层实现就好了，实现每个模块，也和其他的实现的模块没有关系，这样也符合了上面介绍的开闭原则。
// 这样实现起来每个模块只依赖对象的接口，而和其他模块没关系，依赖关系单一。系统容易扩展和维护。
// 我们在实现业务逻辑也是一样，只需要参考抽象层的接口来业务就好了，抽象层暴露出来的接口就是我们业务层可以使用的方法，然后可以通过多态的功能，接口指针指向哪个实现模块，
// 调用了就是具体的实现方法，这样我们业务逻辑层也是依赖抽象成编程。
// 我们就将这种的设计原则叫做依赖倒转原则。都依赖接口而非实现

// ===== >   抽象层  < ========
type Car interface {
	Run()
}

type Driver interface {
	Drive(car Car)
}

// ===== >   实现层  < ========
type BenZ struct {
	//...
}

func (benz *BenZ) Run() {
	fmt.Println("Benz is running...")
}

type Bmw struct {
	//...
}

func (bmw *Bmw) Run() {
	fmt.Println("Bmw is running...")
}

type Zhang_3 struct {
	//...
}

func (zhang3 *Zhang_3) Drive(car Car) {
	fmt.Println("Zhang3 drive car")
	car.Run()
}

type Li_4 struct {
	//...
}

func (li4 *Li_4) Drive(car Car) {
	fmt.Println("li4 drive car")
	car.Run()
}

// 业务逻辑层，只面向接口开发
func TestUsingInterface(t *testing.T) {
	//张3 开 宝马
	var bmw Car
	bmw = &Bmw{}

	var zhang3 Driver
	zhang3 = &Zhang_3{}

	zhang3.Drive(bmw)

	//李4 开 奔驰
	var benz Car
	benz = &BenZ{}

	var li4 Driver
	li4 = &Li_4{}

	li4.Drive(benz)
}
