package trace

import (
	"fmt"
	"io"
)

// Tracer는 코드 전체에서 이벤트를 추적할 수 있는 객체를 설명하는 인터페이스다.
type Tracer interface { // 시작이 대문자이므로 public 타입
	Trace(...interface{}) // Trace 메소드가 어떤 타입의 인수를 0개 이상 허용한다는 것을 나타낸다.

}

type tracer struct {
	out io.Writer
}

func (t *tracer) Trace(a ...interface{}) { // 빈 인터페이스는 어떠한 타입도 담을 수 있는 컨테이너라고 볼 수 있다.
	fmt.Fprint(t.out, a...)
	fmt.Fprintln(t.out)
}

func New(w io.Writer) Tracer {
	return &tracer{out: w}
}
