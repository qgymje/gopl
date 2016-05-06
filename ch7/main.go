package main

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

const doc = `
0. How to think in interface? 

1. Interface types express generalizations or abstactions about the behaviors of other types. By generalizing, interfaces let us write functions that are more flexible and adaptable because they are not tied to the details of one particular implementation.
简单来说, 接口就是为了抽象与一般性, 想想拖线版, usb, 电脑主板
但是代码里, 必须规定接口里的方法(参数类型与返回类型), 所以弱类型语言有接口有啥用? 比如PHP的接口? 像js/ruby/python都没有接口

2. Many object-oriented languages have some notion of interfaces, but what makes Go's interfaces so distincitve is that the are satified implicitly. In other words, there's no need to declare all the interfaces that a given concrete type satisfis; simple processing the necessary methods is enough. This design lets you create new interfaces that are satisfied by existing concrete types without changing the existing types, which is particularyly useful for types defined in packages that you don't control.
隐式实现, 也就是一个对象, 实现多个了接口里的方法, 但不需要声明接口类型, 编译器会检查
如果想让一个第三方包里的对象也符合我自己写的接口, 我无法改它的源代码, 不能在它的声明里加上implements MyInterface, 这个时候, 这个特性就是牛逼的, 同是它也有像Java里需要声明的作用

In this chapter, we'll start by looking at the basic mechanics of interface types and their values. Along the way, we'll study serval important interfaces from the standard library. Many Go programs make as much use of standard interfaces as they do of their own ones. Finally, we'll look at type assertions and type switches and see how they enable a different kind of generality.
标准库里的接口很好用, 如果想让自己的一个自定义的struct实现io.Reader或者io.Writer, 最简单的方法是让bytes.Buffer作为一个field嵌入, 就能拥有其功能
type assertion v.(float64)
type switch swiftch v.(type)


Interfaces as Contracts

ALl the types we've looked at so far have been concrete types. Aconcrete type specifies the exact representation of its values and exposes the intrinsic operations of that representation, such as arithmetic for numbers, or indexing , append, and range for slices. A concrete type may also provide additional behaviors throught its methods. When you have a value of a concrete type, you know exactly what it is and what you can do with it.
concrete type 实体, 会有相对应的方法, 是一个十分具体的对象, 我们知道它有哪些用

There is another kind of type in Go called an interface type. An interface is an abstract type. It doesn't expose the representation or internal structure of its values, or the set of basic operations they support; it reveals only some of their methods. When you have a value of an interface type, you know nothing about what it is; you know only what it can do, or more precisely, what behaviors are provided by its mehtods.
interface type 虚体, 只有行为, 没有实现, 是实体的抽象层, 像虚拟世界, 像电影, 像艺术

Thoughout the book, we've been using two similar functoins for string formating: fmt.Printf, which writes the result to the standard output (a file), and fmt.Sprintf, which returns the result as a string. It would be unfortunate if the hard part, formatting the result, had to be duplicated because of these superficial differences in how the result is used. Thanks to interfaces, it does not. Both of these function are, in effect, wrappers around a third functin, fmt.Fprintf, that is agnostic aboult what happens to the result it computes:
package fmt

func Fptinrf(w io.Writer, format string, args ...interface{}) (int, error)

func Printf(format string, args ...interface{}) (int, err) {
	return Fprintf(os.Stdout, format, args...)
}

func Sprintf(format string, args ...interface{}) string {
	var buf bytes.Buffer
	Fprintf(&buf, format, args...)
	return buf.String()
}
一个interface类型, 是一组定义的协议, 不是实现者, 它们作为参数传入/传出时候, 承载这些协议概念的是实现了这些行为的实体(concrete type)

7.2 Inerface Types
An interface type specifies a set of methods that a concrete type must possess to be considered an instance of that interface.

The io.Writer type is one of the most widely used interfaces because it provides an abstraction of all the types to which bytes can be written, which includes files, memory buffers, network connections, HTTP clients, archivers, hashers, and so on. The io package defines many oter useful interfaces. A reader represents any type from which you can read byets, and a Closer is any value that you can close, such as a file or a network connecitn.(By now you've probably noticed the naming convention for many of Go's single-method interfaces.)

package io

type Reader interface{
	Read(p []byte)(n int, err error)
}
// 注意看这个接口, 以前从来没有仔细看过
type Closer interface{
	Close() error
}

7.3 Interface Satisfaction
A type satifies an interface if it possesses all the methods the interface requries. For example, an *os.File satisfiles io.Reader, Writer, Closer, and ReadWriter. A *bytes.Buffer satisfies Reader, Writer, and ReadWriter, but does not satisfy Closer because it does not have a Close method. As a shorthand, Go programmers often say that a concrete type "is a" particular interface type, meaning that it satisfies the itnerface. For example, a *bytes.Buffer is an io.Writer; an *os.File is an io.ReadWriter.

"它是一个狗"

The assignability rule for interfaces is very simple: an expression may be assigned to an interface only if its type satisfies the interface.

写代码的过程, 就是分类的过程, 再给每一个类添加功能, 就好像让它活起来了.
写代码的过程, 也是提炼共同的过程, 在不同分类的事物中发现相同的特性, 将其指出形成接口.

`

func init() {
	fmt.Println(doc)
}

//WordCounter 是一个单词计数器
type WordCounter int

func (w *WordCounter) Write(p []byte) (n int, err error) {
	n, _, err = bufio.ScanWords(p, false)
	*w += WordCounter(n)
	return n, err
}

//MyType 内嵌了bytes.Buffer, 以实现io.Reader以及io.Writer
type MyType struct {
	bytes.Buffer
}

func main() {
	mt := MyType{}
	fmt.Fprintf(&mt, "%s - %s", "hello", "world")
	fmt.Println("mt.String() as a reader :", mt.String())

	var wc WordCounter
	wc.Write([]byte("who am i who am i"))
	fmt.Println("word counter: ", wc)

	// hasher
	hash := md5.Sum([]byte("hello world"))
	fmt.Printf("md5 = %s\n", hex.EncodeToString(hash[:]))

	hash2 := md5.New()
	hash2.Write([]byte("hello world"))
	fmt.Printf("md5 = %s\n", hex.EncodeToString(hash2.Sum(nil)))
}
