// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: search.proto

package protocol

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// 更新索引信息
type UpdateSearchData struct {
	Opts map[int32]*StarOpt `protobuf:"bytes,1,rep,name=opts" json:"opts,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value"`
	Sync bool               `protobuf:"varint,3,opt,name=sync,proto3" json:"sync,omitempty"`
}

func (m *UpdateSearchData) Reset()                    { *m = UpdateSearchData{} }
func (m *UpdateSearchData) String() string            { return proto.CompactTextString(m) }
func (*UpdateSearchData) ProtoMessage()               {}
func (*UpdateSearchData) Descriptor() ([]byte, []int) { return fileDescriptorSearch, []int{0} }

func (m *UpdateSearchData) GetOpts() map[int32]*StarOpt {
	if m != nil {
		return m.Opts
	}
	return nil
}

func (m *UpdateSearchData) GetSync() bool {
	if m != nil {
		return m.Sync
	}
	return false
}

// 更新活跃的星球
type UpdateRandomStar struct {
	StarID int32           `protobuf:"varint,1,opt,name=starID,proto3" json:"starID,omitempty"`
	Opt    map[int32]int32 `protobuf:"bytes,2,rep,name=opt" json:"opt,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
	Sync   bool            `protobuf:"varint,3,opt,name=sync,proto3" json:"sync,omitempty"`
}

func (m *UpdateRandomStar) Reset()                    { *m = UpdateRandomStar{} }
func (m *UpdateRandomStar) String() string            { return proto.CompactTextString(m) }
func (*UpdateRandomStar) ProtoMessage()               {}
func (*UpdateRandomStar) Descriptor() ([]byte, []int) { return fileDescriptorSearch, []int{1} }

func (m *UpdateRandomStar) GetStarID() int32 {
	if m != nil {
		return m.StarID
	}
	return 0
}

func (m *UpdateRandomStar) GetOpt() map[int32]int32 {
	if m != nil {
		return m.Opt
	}
	return nil
}

func (m *UpdateRandomStar) GetSync() bool {
	if m != nil {
		return m.Sync
	}
	return false
}

type StarOpt struct {
	Opt map[int32]int32 `protobuf:"bytes,1,rep,name=opt" json:"opt,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
}

func (m *StarOpt) Reset()                    { *m = StarOpt{} }
func (m *StarOpt) String() string            { return proto.CompactTextString(m) }
func (*StarOpt) ProtoMessage()               {}
func (*StarOpt) Descriptor() ([]byte, []int) { return fileDescriptorSearch, []int{2} }

func (m *StarOpt) GetOpt() map[int32]int32 {
	if m != nil {
		return m.Opt
	}
	return nil
}

type UpdateSearchHelpData struct {
	Uid   []int32 `protobuf:"varint,1,rep,packed,name=uid" json:"uid,omitempty"`
	Opt   int32   `protobuf:"varint,2,opt,name=opt,proto3" json:"opt,omitempty"`
	Param int32   `protobuf:"varint,3,opt,name=param,proto3" json:"param,omitempty"`
	Sync  bool    `protobuf:"varint,4,opt,name=sync,proto3" json:"sync,omitempty"`
}

func (m *UpdateSearchHelpData) Reset()                    { *m = UpdateSearchHelpData{} }
func (m *UpdateSearchHelpData) String() string            { return proto.CompactTextString(m) }
func (*UpdateSearchHelpData) ProtoMessage()               {}
func (*UpdateSearchHelpData) Descriptor() ([]byte, []int) { return fileDescriptorSearch, []int{3} }

func (m *UpdateSearchHelpData) GetUid() []int32 {
	if m != nil {
		return m.Uid
	}
	return nil
}

func (m *UpdateSearchHelpData) GetOpt() int32 {
	if m != nil {
		return m.Opt
	}
	return 0
}

func (m *UpdateSearchHelpData) GetParam() int32 {
	if m != nil {
		return m.Param
	}
	return 0
}

func (m *UpdateSearchHelpData) GetSync() bool {
	if m != nil {
		return m.Sync
	}
	return false
}

type RandomHelpTarget struct {
	Uid      int32 `protobuf:"varint,1,opt,name=uid,proto3" json:"uid,omitempty"`
	StarType int32 `protobuf:"varint,2,opt,name=starType,proto3" json:"starType,omitempty"`
	ItemID   int32 `protobuf:"varint,3,opt,name=itemID,proto3" json:"itemID,omitempty"`
	Count    int32 `protobuf:"varint,4,opt,name=count,proto3" json:"count,omitempty"`
}

func (m *RandomHelpTarget) Reset()                    { *m = RandomHelpTarget{} }
func (m *RandomHelpTarget) String() string            { return proto.CompactTextString(m) }
func (*RandomHelpTarget) ProtoMessage()               {}
func (*RandomHelpTarget) Descriptor() ([]byte, []int) { return fileDescriptorSearch, []int{4} }

func (m *RandomHelpTarget) GetUid() int32 {
	if m != nil {
		return m.Uid
	}
	return 0
}

func (m *RandomHelpTarget) GetStarType() int32 {
	if m != nil {
		return m.StarType
	}
	return 0
}

func (m *RandomHelpTarget) GetItemID() int32 {
	if m != nil {
		return m.ItemID
	}
	return 0
}

func (m *RandomHelpTarget) GetCount() int32 {
	if m != nil {
		return m.Count
	}
	return 0
}

type RandomHelpTargetRet struct {
	Targets []int32 `protobuf:"varint,1,rep,packed,name=targets" json:"targets,omitempty"`
}

func (m *RandomHelpTargetRet) Reset()                    { *m = RandomHelpTargetRet{} }
func (m *RandomHelpTargetRet) String() string            { return proto.CompactTextString(m) }
func (*RandomHelpTargetRet) ProtoMessage()               {}
func (*RandomHelpTargetRet) Descriptor() ([]byte, []int) { return fileDescriptorSearch, []int{5} }

func (m *RandomHelpTargetRet) GetTargets() []int32 {
	if m != nil {
		return m.Targets
	}
	return nil
}

func init() {
	proto.RegisterType((*UpdateSearchData)(nil), "protocol.update_search_data")
	proto.RegisterType((*UpdateRandomStar)(nil), "protocol.update_random_star")
	proto.RegisterType((*StarOpt)(nil), "protocol.star_opt")
	proto.RegisterType((*UpdateSearchHelpData)(nil), "protocol.update_search_help_data")
	proto.RegisterType((*RandomHelpTarget)(nil), "protocol.random_help_target")
	proto.RegisterType((*RandomHelpTargetRet)(nil), "protocol.random_help_target_ret")
}
func (m *UpdateSearchData) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *UpdateSearchData) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Opts) > 0 {
		for k, _ := range m.Opts {
			dAtA[i] = 0xa
			i++
			v := m.Opts[k]
			msgSize := 0
			if v != nil {
				msgSize = v.Size()
				msgSize += 1 + sovSearch(uint64(msgSize))
			}
			mapSize := 1 + sovSearch(uint64(k)) + msgSize
			i = encodeVarintSearch(dAtA, i, uint64(mapSize))
			dAtA[i] = 0x8
			i++
			i = encodeVarintSearch(dAtA, i, uint64(k))
			if v != nil {
				dAtA[i] = 0x12
				i++
				i = encodeVarintSearch(dAtA, i, uint64(v.Size()))
				n1, err := v.MarshalTo(dAtA[i:])
				if err != nil {
					return 0, err
				}
				i += n1
			}
		}
	}
	if m.Sync {
		dAtA[i] = 0x18
		i++
		if m.Sync {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i++
	}
	return i, nil
}

func (m *UpdateRandomStar) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *UpdateRandomStar) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.StarID != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintSearch(dAtA, i, uint64(m.StarID))
	}
	if len(m.Opt) > 0 {
		for k, _ := range m.Opt {
			dAtA[i] = 0x12
			i++
			v := m.Opt[k]
			mapSize := 1 + sovSearch(uint64(k)) + 1 + sovSearch(uint64(v))
			i = encodeVarintSearch(dAtA, i, uint64(mapSize))
			dAtA[i] = 0x8
			i++
			i = encodeVarintSearch(dAtA, i, uint64(k))
			dAtA[i] = 0x10
			i++
			i = encodeVarintSearch(dAtA, i, uint64(v))
		}
	}
	if m.Sync {
		dAtA[i] = 0x18
		i++
		if m.Sync {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i++
	}
	return i, nil
}

func (m *StarOpt) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *StarOpt) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Opt) > 0 {
		for k, _ := range m.Opt {
			dAtA[i] = 0xa
			i++
			v := m.Opt[k]
			mapSize := 1 + sovSearch(uint64(k)) + 1 + sovSearch(uint64(v))
			i = encodeVarintSearch(dAtA, i, uint64(mapSize))
			dAtA[i] = 0x8
			i++
			i = encodeVarintSearch(dAtA, i, uint64(k))
			dAtA[i] = 0x10
			i++
			i = encodeVarintSearch(dAtA, i, uint64(v))
		}
	}
	return i, nil
}

func (m *UpdateSearchHelpData) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *UpdateSearchHelpData) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Uid) > 0 {
		dAtA3 := make([]byte, len(m.Uid)*10)
		var j2 int
		for _, num1 := range m.Uid {
			num := uint64(num1)
			for num >= 1<<7 {
				dAtA3[j2] = uint8(uint64(num)&0x7f | 0x80)
				num >>= 7
				j2++
			}
			dAtA3[j2] = uint8(num)
			j2++
		}
		dAtA[i] = 0xa
		i++
		i = encodeVarintSearch(dAtA, i, uint64(j2))
		i += copy(dAtA[i:], dAtA3[:j2])
	}
	if m.Opt != 0 {
		dAtA[i] = 0x10
		i++
		i = encodeVarintSearch(dAtA, i, uint64(m.Opt))
	}
	if m.Param != 0 {
		dAtA[i] = 0x18
		i++
		i = encodeVarintSearch(dAtA, i, uint64(m.Param))
	}
	if m.Sync {
		dAtA[i] = 0x20
		i++
		if m.Sync {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i++
	}
	return i, nil
}

func (m *RandomHelpTarget) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *RandomHelpTarget) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Uid != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintSearch(dAtA, i, uint64(m.Uid))
	}
	if m.StarType != 0 {
		dAtA[i] = 0x10
		i++
		i = encodeVarintSearch(dAtA, i, uint64(m.StarType))
	}
	if m.ItemID != 0 {
		dAtA[i] = 0x18
		i++
		i = encodeVarintSearch(dAtA, i, uint64(m.ItemID))
	}
	if m.Count != 0 {
		dAtA[i] = 0x20
		i++
		i = encodeVarintSearch(dAtA, i, uint64(m.Count))
	}
	return i, nil
}

func (m *RandomHelpTargetRet) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *RandomHelpTargetRet) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Targets) > 0 {
		dAtA5 := make([]byte, len(m.Targets)*10)
		var j4 int
		for _, num1 := range m.Targets {
			num := uint64(num1)
			for num >= 1<<7 {
				dAtA5[j4] = uint8(uint64(num)&0x7f | 0x80)
				num >>= 7
				j4++
			}
			dAtA5[j4] = uint8(num)
			j4++
		}
		dAtA[i] = 0xa
		i++
		i = encodeVarintSearch(dAtA, i, uint64(j4))
		i += copy(dAtA[i:], dAtA5[:j4])
	}
	return i, nil
}

func encodeVarintSearch(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *UpdateSearchData) Size() (n int) {
	var l int
	_ = l
	if len(m.Opts) > 0 {
		for k, v := range m.Opts {
			_ = k
			_ = v
			l = 0
			if v != nil {
				l = v.Size()
				l += 1 + sovSearch(uint64(l))
			}
			mapEntrySize := 1 + sovSearch(uint64(k)) + l
			n += mapEntrySize + 1 + sovSearch(uint64(mapEntrySize))
		}
	}
	if m.Sync {
		n += 2
	}
	return n
}

func (m *UpdateRandomStar) Size() (n int) {
	var l int
	_ = l
	if m.StarID != 0 {
		n += 1 + sovSearch(uint64(m.StarID))
	}
	if len(m.Opt) > 0 {
		for k, v := range m.Opt {
			_ = k
			_ = v
			mapEntrySize := 1 + sovSearch(uint64(k)) + 1 + sovSearch(uint64(v))
			n += mapEntrySize + 1 + sovSearch(uint64(mapEntrySize))
		}
	}
	if m.Sync {
		n += 2
	}
	return n
}

func (m *StarOpt) Size() (n int) {
	var l int
	_ = l
	if len(m.Opt) > 0 {
		for k, v := range m.Opt {
			_ = k
			_ = v
			mapEntrySize := 1 + sovSearch(uint64(k)) + 1 + sovSearch(uint64(v))
			n += mapEntrySize + 1 + sovSearch(uint64(mapEntrySize))
		}
	}
	return n
}

func (m *UpdateSearchHelpData) Size() (n int) {
	var l int
	_ = l
	if len(m.Uid) > 0 {
		l = 0
		for _, e := range m.Uid {
			l += sovSearch(uint64(e))
		}
		n += 1 + sovSearch(uint64(l)) + l
	}
	if m.Opt != 0 {
		n += 1 + sovSearch(uint64(m.Opt))
	}
	if m.Param != 0 {
		n += 1 + sovSearch(uint64(m.Param))
	}
	if m.Sync {
		n += 2
	}
	return n
}

func (m *RandomHelpTarget) Size() (n int) {
	var l int
	_ = l
	if m.Uid != 0 {
		n += 1 + sovSearch(uint64(m.Uid))
	}
	if m.StarType != 0 {
		n += 1 + sovSearch(uint64(m.StarType))
	}
	if m.ItemID != 0 {
		n += 1 + sovSearch(uint64(m.ItemID))
	}
	if m.Count != 0 {
		n += 1 + sovSearch(uint64(m.Count))
	}
	return n
}

func (m *RandomHelpTargetRet) Size() (n int) {
	var l int
	_ = l
	if len(m.Targets) > 0 {
		l = 0
		for _, e := range m.Targets {
			l += sovSearch(uint64(e))
		}
		n += 1 + sovSearch(uint64(l)) + l
	}
	return n
}

func sovSearch(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozSearch(x uint64) (n int) {
	return sovSearch(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *UpdateSearchData) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSearch
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: update_search_data: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: update_search_data: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Opts", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSearch
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthSearch
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Opts == nil {
				m.Opts = make(map[int32]*StarOpt)
			}
			var mapkey int32
			var mapvalue *StarOpt
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowSearch
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					wire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				fieldNum := int32(wire >> 3)
				if fieldNum == 1 {
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowSearch
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						mapkey |= (int32(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
				} else if fieldNum == 2 {
					var mapmsglen int
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowSearch
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						mapmsglen |= (int(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					if mapmsglen < 0 {
						return ErrInvalidLengthSearch
					}
					postmsgIndex := iNdEx + mapmsglen
					if mapmsglen < 0 {
						return ErrInvalidLengthSearch
					}
					if postmsgIndex > l {
						return io.ErrUnexpectedEOF
					}
					mapvalue = &StarOpt{}
					if err := mapvalue.Unmarshal(dAtA[iNdEx:postmsgIndex]); err != nil {
						return err
					}
					iNdEx = postmsgIndex
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipSearch(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if skippy < 0 {
						return ErrInvalidLengthSearch
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.Opts[mapkey] = mapvalue
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Sync", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSearch
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Sync = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipSearch(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthSearch
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *UpdateRandomStar) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSearch
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: update_random_star: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: update_random_star: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field StarID", wireType)
			}
			m.StarID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSearch
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.StarID |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Opt", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSearch
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthSearch
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Opt == nil {
				m.Opt = make(map[int32]int32)
			}
			var mapkey int32
			var mapvalue int32
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowSearch
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					wire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				fieldNum := int32(wire >> 3)
				if fieldNum == 1 {
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowSearch
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						mapkey |= (int32(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
				} else if fieldNum == 2 {
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowSearch
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						mapvalue |= (int32(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipSearch(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if skippy < 0 {
						return ErrInvalidLengthSearch
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.Opt[mapkey] = mapvalue
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Sync", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSearch
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Sync = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipSearch(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthSearch
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *StarOpt) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSearch
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: star_opt: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: star_opt: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Opt", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSearch
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthSearch
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Opt == nil {
				m.Opt = make(map[int32]int32)
			}
			var mapkey int32
			var mapvalue int32
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowSearch
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					wire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				fieldNum := int32(wire >> 3)
				if fieldNum == 1 {
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowSearch
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						mapkey |= (int32(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
				} else if fieldNum == 2 {
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowSearch
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						mapvalue |= (int32(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipSearch(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if skippy < 0 {
						return ErrInvalidLengthSearch
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.Opt[mapkey] = mapvalue
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipSearch(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthSearch
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *UpdateSearchHelpData) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSearch
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: update_search_help_data: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: update_search_help_data: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType == 0 {
				var v int32
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowSearch
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					v |= (int32(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				m.Uid = append(m.Uid, v)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowSearch
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					packedLen |= (int(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if packedLen < 0 {
					return ErrInvalidLengthSearch
				}
				postIndex := iNdEx + packedLen
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				for iNdEx < postIndex {
					var v int32
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowSearch
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						v |= (int32(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					m.Uid = append(m.Uid, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field Uid", wireType)
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Opt", wireType)
			}
			m.Opt = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSearch
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Opt |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Param", wireType)
			}
			m.Param = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSearch
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Param |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Sync", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSearch
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Sync = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipSearch(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthSearch
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *RandomHelpTarget) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSearch
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: random_help_target: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: random_help_target: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Uid", wireType)
			}
			m.Uid = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSearch
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Uid |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field StarType", wireType)
			}
			m.StarType = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSearch
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.StarType |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ItemID", wireType)
			}
			m.ItemID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSearch
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ItemID |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Count", wireType)
			}
			m.Count = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSearch
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Count |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipSearch(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthSearch
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *RandomHelpTargetRet) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSearch
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: random_help_target_ret: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: random_help_target_ret: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType == 0 {
				var v int32
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowSearch
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					v |= (int32(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				m.Targets = append(m.Targets, v)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowSearch
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					packedLen |= (int(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if packedLen < 0 {
					return ErrInvalidLengthSearch
				}
				postIndex := iNdEx + packedLen
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				for iNdEx < postIndex {
					var v int32
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowSearch
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						v |= (int32(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					m.Targets = append(m.Targets, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field Targets", wireType)
			}
		default:
			iNdEx = preIndex
			skippy, err := skipSearch(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthSearch
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipSearch(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowSearch
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowSearch
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowSearch
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			iNdEx += length
			if length < 0 {
				return 0, ErrInvalidLengthSearch
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowSearch
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipSearch(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthSearch = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowSearch   = fmt.Errorf("proto: integer overflow")
)

func init() { proto.RegisterFile("search.proto", fileDescriptorSearch) }

var fileDescriptorSearch = []byte{
	// 380 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x51, 0xcf, 0x4a, 0xc3, 0x30,
	0x18, 0x27, 0xeb, 0x3a, 0xeb, 0xb7, 0x1d, 0x46, 0x18, 0xb3, 0x4c, 0x18, 0x65, 0xa0, 0xf4, 0x62,
	0x0f, 0x13, 0x54, 0x76, 0x94, 0x79, 0x18, 0x1e, 0x84, 0xe2, 0xbd, 0xc4, 0x36, 0x6c, 0xc3, 0xad,
	0x89, 0x69, 0x2a, 0xf4, 0xb1, 0x7c, 0x0b, 0x8f, 0x3e, 0x82, 0xec, 0x49, 0x24, 0x49, 0xbb, 0x75,
	0x74, 0x27, 0x4f, 0xfd, 0x7e, 0xe5, 0x97, 0xef, 0xf7, 0xe7, 0x83, 0x5e, 0x46, 0x89, 0x88, 0x57,
	0x01, 0x17, 0x4c, 0x32, 0xec, 0xe8, 0x4f, 0xcc, 0x36, 0x93, 0x2f, 0x04, 0x38, 0xe7, 0x09, 0x91,
	0x34, 0x32, 0x8c, 0x28, 0x21, 0x92, 0xe0, 0x19, 0xb4, 0x19, 0x97, 0x99, 0x8b, 0x3c, 0xcb, 0xef,
	0x4e, 0xaf, 0x83, 0x8a, 0x1f, 0x34, 0xb9, 0xc1, 0x0b, 0x97, 0xd9, 0x53, 0x2a, 0x45, 0x11, 0xea,
	0x37, 0x18, 0x43, 0x3b, 0x2b, 0xd2, 0xd8, 0xb5, 0x3c, 0xe4, 0x3b, 0xa1, 0x9e, 0x47, 0xcf, 0x70,
	0xbe, 0xa7, 0xe1, 0x3e, 0x58, 0xef, 0xb4, 0x70, 0x91, 0x87, 0x7c, 0x3b, 0x54, 0x23, 0xf6, 0xc1,
	0xfe, 0x24, 0x9b, 0x9c, 0xba, 0x2d, 0x0f, 0xf9, 0xdd, 0x29, 0x3e, 0xe8, 0x65, 0x92, 0x88, 0x88,
	0x71, 0x19, 0x1a, 0xc2, 0xac, 0xf5, 0x80, 0xea, 0x9e, 0x05, 0x49, 0x13, 0xb6, 0x8d, 0x14, 0x0b,
	0x0f, 0xa1, 0xa3, 0xbe, 0x8b, 0x79, 0xb9, 0xb9, 0x44, 0xf8, 0x1e, 0x2c, 0xc6, 0xa5, 0xdb, 0xd2,
	0x51, 0xae, 0x1a, 0x51, 0x6a, 0x2b, 0x54, 0x14, 0x93, 0x44, 0xbd, 0x38, 0x19, 0xe4, 0x0e, 0x9c,
	0x8a, 0x74, 0x22, 0xc7, 0xa0, 0x9e, 0xc3, 0xae, 0x7b, 0xfe, 0x00, 0xa7, 0x8a, 0x82, 0x6f, 0x8c,
	0x21, 0xd3, 0xed, 0x65, 0x33, 0xeb, 0xb1, 0x8d, 0x7f, 0x4b, 0x2e, 0xe1, 0xe2, 0xf8, 0x5a, 0x2b,
	0xba, 0xe1, 0xe6, 0xbc, 0x7d, 0xb0, 0xf2, 0x75, 0xa2, 0x1d, 0xd8, 0xa1, 0x1a, 0xd5, 0x1f, 0x53,
	0x92, 0x5e, 0xac, 0x5c, 0x0e, 0xc0, 0xe6, 0x44, 0x90, 0xad, 0x8e, 0x6f, 0x87, 0x06, 0xec, 0x3b,
	0x69, 0x1f, 0x3a, 0x99, 0x70, 0xc0, 0x65, 0x89, 0x5a, 0x41, 0x12, 0xb1, 0xa4, 0xf2, 0xa0, 0x81,
	0x2a, 0x8d, 0x91, 0xe9, 0xe0, 0xb5, 0xe0, 0x95, 0xdb, 0x3d, 0x56, 0xc7, 0x5b, 0x4b, 0xba, 0x5d,
	0xcc, 0x4b, 0xb9, 0x12, 0x29, 0x17, 0x31, 0xcb, 0x53, 0xa9, 0x05, 0xed, 0xd0, 0x80, 0xc9, 0x14,
	0x86, 0x4d, 0xc5, 0x48, 0x50, 0x89, 0x5d, 0x38, 0x33, 0x28, 0x2b, 0xd3, 0x55, 0xf0, 0xb1, 0xf7,
	0xbd, 0x1b, 0xa3, 0x9f, 0xdd, 0x18, 0xfd, 0xee, 0xc6, 0xe8, 0xad, 0xa3, 0x5b, 0xbf, 0xfd, 0x0b,
	0x00, 0x00, 0xff, 0xff, 0xc3, 0xc3, 0x3c, 0xa1, 0x18, 0x03, 0x00, 0x00,
}
