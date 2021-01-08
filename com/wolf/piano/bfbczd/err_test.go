package bfbczd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime/debug"
	"testing"
)

type MyError struct {
	Inner      error
	Message    string
	StackTrace string
	Misc       map[string]interface{}
}

func wrapError(err error, messagef string, msgArgs ...interface{}) MyError {
	return MyError{
		Inner:      err, //1
		Message:    fmt.Sprintf(messagef, msgArgs...),
		StackTrace: string(debug.Stack()),        //2
		Misc:       make(map[string]interface{}), //3
	}
}

func (err MyError) Error() string {
	return err.Message
} //1.在这里存储我们正在包装的错误。 如果需要调查发生的事情，我们总是希望能够查看到最低级别的错误。
//2.这行代码记录了创建错误时的堆栈跟踪。
//3.这里我们创建一个杂项信息存储字段。可以存储并发ID，堆栈跟踪的hash或可能有助于诊断错误的其他上下文信息。

// "lowlevel" module

type LowLevelErr struct {
	error
}

// 最底层，抛出异常是根源，但是并不是用户明白的
func isGloballyExec(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, LowLevelErr{wrapError(err, err.Error())} // 1
	}
	return info.Mode().Perm()&0100 == 0100, nil
} //1.在这里，我们用自定义错误来封装os.Stat中的原始错误。在这种情况下，我们不会掩盖这个错误产生的信息。

// "intermediate" module

type IntermediateErr struct {
	error
}

func runJobErr(id string) error {
	const jobBinPath = "/bad/job/binary"
	isExecutable, err := isGloballyExec(jobBinPath)
	if err != nil {
		return err //1
	} else if isExecutable == false {
		return wrapError(nil, "job binary is not executable")
	}

	return exec.Command(jobBinPath, "--id="+id).Run() //1
} //1.我们传递来自 lowlevel 模块的错误，由于我们接收从其他模块传递的错误而没有将它们包装在我们自己的错误类型中，这将会产生问题。

// 中间层，调用最底层，有异常，它是给开发看的，针对用户可以更友好提示。
func runJobRight(id string) error {
	const jobBinPath = "/bad/job/binary"
	isExecutable, err := isGloballyExec(jobBinPath)
	if err != nil {
		return IntermediateErr{wrapError(err,
			"cannot run job %q: requisite binaries not available", id)} //1
	} else if isExecutable == false {
		return wrapError(
			nil,
			"cannot run job %q: requisite binaries are not executable", id,
		)
	}

	return exec.Command(jobBinPath, "--id="+id).Run()
} //1.在这里，我们现在使用自定义错误。我们想隐藏工作未运行原因的底层细节，因为这对于用户并不重要。

func handleError(key int, err error, message string) {
	// 日志
	log.SetPrefix(fmt.Sprintf("[logID: %v]: ", key))
	log.Printf("%#v", err) //3
	// 输出
	fmt.Printf("[%v] %v", key, message)
}

// 用户直接调用层，返回的错误应该是用户可明白，而日志记录完整信息便于开发人员排查
func TestErrUsing(t *testing.T) {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)

	err := runJobRight("1")
	if err != nil {
		msg := "There was an unexpected issue; please report this as a bug."
		if _, ok := err.(IntermediateErr); ok { //1
			msg = err.Error()
		}
		handleError(1, err, msg) //2
	}
} //1.在这里我们检查是否错误是预期的类型。 如果是，可以简单地将其消息传递给用户。
//2.在这一行中，将日志和错误消息与ID绑定在一起。我们可以很容易增加这个增量，或者使用一个GUID来确保一个唯一的ID。
//3.在这里我们记录完整的错误，以备需要深入了解发生了什么。
