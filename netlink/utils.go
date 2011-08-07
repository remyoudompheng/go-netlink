package netlink

import (
	"os"
	"net"
	"encoding/binary"
	"reflect"
	"bytes"
	"syscall"
)

func netlinkPadding(size int) int {
	partialChunk := size % syscall.NLMSG_ALIGNTO
	return (syscall.NLMSG_ALIGNTO - partialChunk) % syscall.NLMSG_ALIGNTO
}

func skipAlignedFromSlice(r *bytes.Buffer, dataLen int) os.Error {
	r.Next(dataLen + netlinkPadding(dataLen))
	return nil
}

func readAlignedFromSlice(r *bytes.Buffer, data interface{}, dataLen int) os.Error {
	var er os.Error
	switch dest := data.(type) {
	case nil:
		r.Next(dataLen)
	case *[]byte:
		*dest = make([]byte, dataLen)
		_, er = r.Read((*dest)[:])
	case *net.IP:
		*dest = make([]byte, dataLen)
		_, er = r.Read((*dest)[:])
	case *string:
		// Read a NULL-terminated string 
		buffer := make([]byte, dataLen)
		_, er = r.Read(buffer[:])
		*dest = string(buffer[:len(buffer)-1])
	default:
		// Read a binary struct
		er = binary.Read(r, systemEndianness, data)
		realLen := sizeof(data)
		r.Next(dataLen - realLen)
	}
	if er != nil {
		return er
	}
	// advance by the padding size
	r.Next(netlinkPadding(dataLen))
	return nil
}

func putAttribute(w *bytes.Buffer, attrtype uint16, data interface{}) os.Error {
	var attr Attr
	switch data := data.(type) {
	case []byte:
		attr = Attr{Len: uint16(len(data)), Type: attrtype}
		binary.Write(w, systemEndianness, attr)
		binary.Write(w, systemEndianness, data)
	case string:
		attr = Attr{Len: uint16(len(data) + 1), Type: attrtype}
		binary.Write(w, systemEndianness, attr)
		binary.Write(w, systemEndianness, []byte(data))
		w.WriteByte(0)
	default:
		attr = Attr{Len: uint16(sizeof(data)), Type: attrtype}
		binary.Write(w, systemEndianness, attr)
		binary.Write(w, systemEndianness, data)
	}
	for i := 0; i < netlinkPadding(int(attr.Len)); i++ {
		w.WriteByte(0)
	}
	return nil
}

func sizeof(data interface{}) int {
	var v reflect.Value
	switch d := reflect.ValueOf(data); d.Kind() {
	case reflect.Ptr:
		v = d.Elem()
	case reflect.Slice:
		v = d
	default:
		v = d
	}
	return binary.TotalSize(v)
}
