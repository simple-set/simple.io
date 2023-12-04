package simpleHttp

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestHeader_Has(t *testing.T) {
	header := NewHeader()
	header.Set("k1", "v1")
	header.Add("k1", "v1.1")
	k1 := header.Get("k1")
	fmt.Println(k1)
	hasK1 := header.Has("k1")
	fmt.Println(hasK1)

	header.Set("k2", "")
	hasK2 := header.Has("test")
	fmt.Println(hasK2)
	fmt.Println(header.Has("k3"))
	//type fields struct {
	//	Header http.Header
	//	name   string
	//	value  []string
	//}
	//type args struct {
	//	key string
	//}
	//tests := []struct {
	//	name   string
	//	fields fields
	//	args   args
	//	want   bool
	//}{
	//	// TODO: Add test cases.
	//}
	//for _, tt := range tests {
	//	t.Run(tt.name, func(t *testing.T) {
	//		h := &Header{
	//			Header: tt.fields.Header,
	//			name:   tt.fields.name,
	//			value:  tt.fields.value,
	//		}
	//		if got := h.Has(tt.args.key); got != tt.want {
	//			t.Errorf("Has() = %v, want %v", got, tt.want)
	//		}
	//	})
	//}
}

func TestNewHeader(t *testing.T) {
	header := NewHeader()
	header.Add("host", "localhost")
	header.Add("Server", "simple.io/1.1")

	for i, v := range header.Header {
		fmt.Println(i, v)
		for n, k := range v {
			fmt.Println(n, k)
		}
	}
	//tests := []struct {
	//	name string
	//	want *Header
	//}{
	//	// TODO: Add test cases.
	//}
	//for _, tt := range tests {
	//	t.Run(tt.name, func(t *testing.T) {
	//		if got := NewHeader(); !reflect.DeepEqual(got, tt.want) {
	//			t.Errorf("NewHeader() = %v, want %v", got, tt.want)
	//		}
	//	})
	//}
}

func TestHeader_Has1(t *testing.T) {
	type fields struct {
		Header http.Header
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Header{
				Header: tt.fields.Header,
			}
			if got := h.Has(tt.args.key); got != tt.want {
				t.Errorf("Has() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewHeader1(t *testing.T) {
	tests := []struct {
		name string
		want *Header
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHeader(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}
