package trace

import (
	"fmt"
	"io"
)

// Tracer 인터페이스는 코드 전체에서 이벤트를 추적할 수 있는 객체를 설명한다.
type Tracer interface { // 시작이 대문자이므로 public 타입
	Trace(...interface{}) // Trace 메소드가 어떤 타입의 인수를 0개 이상 허용한다는 것을 나타낸다.
}

type tracer struct {
	out io.Writer
}
type nilTracer struct{}

func (t *tracer) Trace(a ...interface{}) { // 빈 인터페이스는 어떠한 타입도 담을 수 있는 컨테이너라고 볼 수 있다.
	fmt.Fprint(t.out, a...) // t에 쓰고
	fmt.Fprintln(t.out)     // 줄바꿈
}
func (t *nilTracer) Trace(a ...interface{}) {}

func New(w io.Writer) Tracer {
	return &tracer{out: w}
}

// Off는 Trace에 대한 호출을 무시할 Tracer를 생성한다. (처음 채팅방을 생성할 때 초기화해야 하는데 그 때는 추적할 필요가 없으니 Off사용)
func Off() Tracer {
	return &nilTracer{}
}
