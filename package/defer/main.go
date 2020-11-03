package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
)

//derfer的执行顺序
//当一个方法中有多个defer时， defer会将要延迟执行的方法“压栈”，当defer被触发时，将所有“压栈”的方法“出栈”并执行。所以defer的执行顺序是LIFO的。
//所以下面这段代码的输出不是1 2 3，而是3 2 1。

func example() {
	defer func() {
		fmt.Println(1)
	}()
	defer func() {
		fmt.Println(2)
	}()
	defer func() {
		fmt.Println(3)
	}()
}

//调用os.Exit时derfer不会执行
func example1() {
	//当发生panic时，所在goroutine的所有defer会被执行，但是当调用os.Exit()方法退出程序时，defer并不会被执行。
	defer func() {
		fmt.Println("hello world")
	}()
	os.Exit(0)
}

//判断执行没有err之后，再defer释放资源
//一些获取资源的操作可能会返回err参数，我们可以选择忽略返回的err参数，但是如果要使用defer进行延迟释放的的话，需要在使用defer之前先判断是否存在err，如果资源没有获取成功，即没有必要也不应该再对资源执行释放操作。如果不判断获取资源是否成功就执行释放操作的话，还有可能导致释放方法执行错误。
func example2() {
	resp, err := http.Get("https://www.baidu.com")
	// 先判断操作是否成功
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// 如果操作成功，再进行Close操作
	defer resp.Body.Close()
}

//在for循环中使用defer可能导致的性能问题
//defer在紧邻创建资源的语句后生命力，看上去逻辑没有什么问题。但是和直接调用相比，defer的执行存在着额外的开销，例如defer会对其后需要的参数进行内存拷贝，还需要对defer结构进行压栈出栈操作。所以在循环中定义defer可能导致大量的资源开销，在本例中，可以将f.Close()语句前的defer去掉，来减少大量defer导致的额外资源消耗。
func example3() {
	for i := 0; i < 5; i++ {
		f, _ := os.Open("/etc/hosts")
		defer f.Close()
	}
}

//defer在匿名返回值和命名返回值函数中的不同表现
func example41() int {
	var result int
	defer func() {
		result++
		fmt.Println("defer")
	}()
	return result

}

func example42() (result int) {
	defer func() {
		result++
		fmt.Println("defer")
	}()
	return result
}

func main() {
	//example()
	//example2()
	//example3()

	//defer在匿名返回值和命名返回值函数中的不同表现
	//上面的方法会输出0，下面的方法输出1。上面的方法使用了匿名返回值，下面的使用了命名返回值，除此之外其他的逻辑均相同，为什么输出的结果会有区别呢？
	//要搞清这个问题首先需要了解defer的执行逻辑，文档中说defer语句在方法返回“时”触发，也就是说return和defer是“同时”执行的。以匿名返回值方法举例，过程如下。
	//将result赋值给返回值（可以理解成Go自动创建了一个返回值retValue，相当于执行retValue = result）
	//然后检查是否有defer，如果有则执行
	//返回刚才创建的返回值（retValue）
	//在这种情况下，defer中的修改是对result执行的，而不是retValue，所以defer返回的依然是retValue。在命名返回值方法中，由于返回值在方法定义时已经被定义，所以没有创建retValue的过程，result就是retValue，defer对于result的修改也会被直接返回。
	fmt.Println(example41())
	fmt.Println(example42())

	//example1()

}
