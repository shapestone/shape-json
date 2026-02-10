package json

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

// encoderFunc appends the JSON encoding of rv to buf, returning the extended buffer.
type encoderFunc func(buf []byte, rv reflect.Value) ([]byte, error)

// encoderCache stores compiled encoders keyed by reflect.Type.
// Uses copy-on-write map behind atomic.Value for lock-free reads.
var encoderCache atomic.Value // holds map[reflect.Type]encoderFunc
var encoderMu sync.Mutex      // protects writes only

func init() {
	encoderCache.Store(make(map[reflect.Type]encoderFunc))
}

// Pre-computed reflect types for special handling.
var (
	marshalerType = reflect.TypeOf((*Marshaler)(nil)).Elem()
	timeType      = reflect.TypeOf(time.Time{})
	durationType  = reflect.TypeOf(time.Duration(0))
)

// encoderForType returns a cached encoder for the given type, building one if needed.
func encoderForType(t reflect.Type) encoderFunc {
	// Fast path: lock-free read
	m := encoderCache.Load().(map[reflect.Type]encoderFunc)
	if enc, ok := m[t]; ok {
		return enc
	}

	// Slow path: build encoder
	encoderMu.Lock()

	// Double-check after lock
	m = encoderCache.Load().(map[reflect.Type]encoderFunc)
	if enc, ok := m[t]; ok {
		encoderMu.Unlock()
		return enc
	}

	// Placeholder for recursive types — allows sub-type builds to find this
	// type in the cache without deadlocking
	var wg sync.WaitGroup
	wg.Add(1)
	var realEnc encoderFunc
	placeholder := func(buf []byte, rv reflect.Value) ([]byte, error) {
		wg.Wait()
		return realEnc(buf, rv)
	}

	// Store placeholder and release lock before building
	newM := make(map[reflect.Type]encoderFunc, len(m)+1)
	for k, v := range m {
		newM[k] = v
	}
	newM[t] = placeholder
	encoderCache.Store(newM)
	encoderMu.Unlock()

	// Build the real encoder (may recursively call encoderForType for sub-types)
	realEnc = buildEncoder(t)

	// Replace placeholder with real encoder
	encoderMu.Lock()
	m = encoderCache.Load().(map[reflect.Type]encoderFunc)
	newM2 := make(map[reflect.Type]encoderFunc, len(m))
	for k, v := range m {
		newM2[k] = v
	}
	newM2[t] = realEnc
	encoderCache.Store(newM2)
	encoderMu.Unlock()
	wg.Done()

	return realEnc
}

// buildEncoder creates an encoder for the given type.
func buildEncoder(t reflect.Type) encoderFunc {
	// Check Marshaler interface on value type
	if t.Implements(marshalerType) {
		return marshalerEnc
	}
	// Check Marshaler on pointer-to-type
	if t.Kind() != reflect.Ptr && reflect.PointerTo(t).Implements(marshalerType) {
		return buildAddrMarshalerEnc(t)
	}

	// Special types
	if t == timeType {
		return timeEnc
	}
	if t == durationType {
		return durationEnc
	}

	switch t.Kind() {
	case reflect.Ptr:
		return buildPtrEncoder(t)
	case reflect.Interface:
		return interfaceEnc
	case reflect.String:
		return stringEnc
	case reflect.Bool:
		return boolEnc
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return intEnc
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return uintEnc
	case reflect.Float32:
		return float32Enc
	case reflect.Float64:
		return float64Enc
	case reflect.Struct:
		return buildStructEncoder(t)
	case reflect.Map:
		return buildMapEncoder(t)
	case reflect.Slice:
		return buildSliceEncoder(t)
	case reflect.Array:
		return buildArrayEncoder(t)
	default:
		return unsupportedEnc(t)
	}
}

// ================================
// Primitive Encoders (zero allocation)
// ================================

func boolEnc(buf []byte, rv reflect.Value) ([]byte, error) {
	if rv.Bool() {
		return append(buf, "true"...), nil
	}
	return append(buf, "false"...), nil
}

func intEnc(buf []byte, rv reflect.Value) ([]byte, error) {
	return strconv.AppendInt(buf, rv.Int(), 10), nil
}

func uintEnc(buf []byte, rv reflect.Value) ([]byte, error) {
	return strconv.AppendUint(buf, rv.Uint(), 10), nil
}

func float32Enc(buf []byte, rv reflect.Value) ([]byte, error) {
	return strconv.AppendFloat(buf, rv.Float(), 'g', -1, 32), nil
}

func float64Enc(buf []byte, rv reflect.Value) ([]byte, error) {
	return strconv.AppendFloat(buf, rv.Float(), 'g', -1, 64), nil
}

func stringEnc(buf []byte, rv reflect.Value) ([]byte, error) {
	buf = append(buf, '"')
	buf = appendEscapedString(buf, rv.String())
	buf = append(buf, '"')
	return buf, nil
}

// ================================
// Special Type Encoders
// ================================

func timeEnc(buf []byte, rv reflect.Value) ([]byte, error) {
	t := rv.Interface().(time.Time)
	buf = append(buf, '"')
	buf = t.AppendFormat(buf, time.RFC3339Nano)
	buf = append(buf, '"')
	return buf, nil
}

func durationEnc(buf []byte, rv reflect.Value) ([]byte, error) {
	d := time.Duration(rv.Int())
	buf = append(buf, '"')
	buf = appendISO8601Duration(buf, d)
	buf = append(buf, '"')
	return buf, nil
}

// ================================
// Marshaler Interface Encoders
// ================================

func marshalerEnc(buf []byte, rv reflect.Value) ([]byte, error) {
	if rv.Kind() == reflect.Ptr && rv.IsNil() {
		return append(buf, "null"...), nil
	}
	m := rv.Interface().(Marshaler)
	b, err := m.MarshalJSON()
	if err != nil {
		return buf, err
	}
	return append(buf, b...), nil
}

func buildAddrMarshalerEnc(t reflect.Type) encoderFunc {
	// Fallback encoder for when we can't take address
	fallback := buildEncoderNoMarshaler(t)
	return func(buf []byte, rv reflect.Value) ([]byte, error) {
		if rv.CanAddr() {
			m := rv.Addr().Interface().(Marshaler)
			b, err := m.MarshalJSON()
			if err != nil {
				return buf, err
			}
			return append(buf, b...), nil
		}
		return fallback(buf, rv)
	}
}

// buildEncoderNoMarshaler builds an encoder skipping the Marshaler check.
func buildEncoderNoMarshaler(t reflect.Type) encoderFunc {
	if t == timeType {
		return timeEnc
	}
	if t == durationType {
		return durationEnc
	}
	switch t.Kind() {
	case reflect.Struct:
		return buildStructEncoder(t)
	case reflect.String:
		return stringEnc
	case reflect.Bool:
		return boolEnc
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return intEnc
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return uintEnc
	case reflect.Float32:
		return float32Enc
	case reflect.Float64:
		return float64Enc
	default:
		return unsupportedEnc(t)
	}
}

// ================================
// Pointer / Interface Encoders
// ================================

func buildPtrEncoder(t reflect.Type) encoderFunc {
	elemEnc := encoderForType(t.Elem())
	return func(buf []byte, rv reflect.Value) ([]byte, error) {
		if rv.IsNil() {
			return append(buf, "null"...), nil
		}
		return elemEnc(buf, rv.Elem())
	}
}

func interfaceEnc(buf []byte, rv reflect.Value) ([]byte, error) {
	if rv.IsNil() {
		return append(buf, "null"...), nil
	}
	// Try the fast path (type switch) before falling back to reflect
	v := rv.Interface()
	buf, err := appendInterface(buf, v)
	if err == errNeedReflect {
		elem := rv.Elem()
		enc := encoderForType(elem.Type())
		return enc(buf, elem)
	}
	return buf, err
}

// ================================
// Struct Encoder
// ================================

// structField holds pre-computed info for a single struct field.
type structField struct {
	index     int                      // field index in struct
	nameBytes []byte                   // pre-encoded `"fieldName":` including quotes and colon
	encoder   encoderFunc              // pre-resolved encoder for this field's type
	omitEmpty bool                     // whether to skip empty values
	emptyFn   func(reflect.Value) bool // pre-resolved empty checker (nil if !omitEmpty)
}

func buildStructEncoder(t reflect.Type) encoderFunc {
	var fields []structField

	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		if sf.PkgPath != "" { // unexported
			continue
		}

		info := getFieldInfo(sf)
		if info.skip {
			continue
		}

		// Pre-encode the JSON key: "fieldName":
		nameBytes := make([]byte, 0, len(info.name)+4)
		nameBytes = append(nameBytes, '"')
		nameBytes = appendEscapedString(nameBytes, info.name)
		nameBytes = append(nameBytes, '"', ':')

		enc := encoderForType(sf.Type)
		if info.asString {
			enc = wrapStringEncoder(enc, sf.Type.Kind())
		}

		f := structField{
			index:     i,
			nameBytes: nameBytes,
			encoder:   enc,
			omitEmpty: info.omitEmpty,
		}

		if info.omitEmpty {
			f.emptyFn = emptyFuncForKind(sf.Type)
		}

		fields = append(fields, f)
	}

	// Sort fields by name ONCE at build time
	sort.Slice(fields, func(i, j int) bool {
		return string(fields[i].nameBytes) < string(fields[j].nameBytes)
	})

	return func(buf []byte, rv reflect.Value) ([]byte, error) {
		buf = append(buf, '{')
		first := true
		for i := range fields {
			f := &fields[i]
			fv := rv.Field(f.index)

			if f.omitEmpty && f.emptyFn(fv) {
				continue
			}

			if !first {
				buf = append(buf, ',')
			}
			first = false

			buf = append(buf, f.nameBytes...)

			var err error
			buf, err = f.encoder(buf, fv)
			if err != nil {
				return buf, err
			}
		}
		buf = append(buf, '}')
		return buf, nil
	}
}

// emptyFuncForKind returns a specialized empty checker for the given type.
func emptyFuncForKind(t reflect.Type) func(reflect.Value) bool {
	switch t.Kind() {
	case reflect.Bool:
		return func(v reflect.Value) bool { return !v.Bool() }
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return func(v reflect.Value) bool { return v.Int() == 0 }
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return func(v reflect.Value) bool { return v.Uint() == 0 }
	case reflect.Float32, reflect.Float64:
		return func(v reflect.Value) bool { return v.Float() == 0 }
	case reflect.String:
		return func(v reflect.Value) bool { return v.Len() == 0 }
	case reflect.Slice, reflect.Map, reflect.Array:
		return func(v reflect.Value) bool { return v.Len() == 0 }
	case reflect.Ptr, reflect.Interface:
		return func(v reflect.Value) bool { return v.IsNil() }
	default:
		return func(v reflect.Value) bool { return false }
	}
}

// wrapStringEncoder wraps an encoder to output the value as a JSON string.
func wrapStringEncoder(inner encoderFunc, kind reflect.Kind) encoderFunc {
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return func(buf []byte, rv reflect.Value) ([]byte, error) {
			buf = append(buf, '"')
			buf = strconv.AppendInt(buf, rv.Int(), 10)
			buf = append(buf, '"')
			return buf, nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return func(buf []byte, rv reflect.Value) ([]byte, error) {
			buf = append(buf, '"')
			buf = strconv.AppendUint(buf, rv.Uint(), 10)
			buf = append(buf, '"')
			return buf, nil
		}
	case reflect.Float32:
		return func(buf []byte, rv reflect.Value) ([]byte, error) {
			buf = append(buf, '"')
			buf = strconv.AppendFloat(buf, rv.Float(), 'g', -1, 32)
			buf = append(buf, '"')
			return buf, nil
		}
	case reflect.Float64:
		return func(buf []byte, rv reflect.Value) ([]byte, error) {
			buf = append(buf, '"')
			buf = strconv.AppendFloat(buf, rv.Float(), 'g', -1, 64)
			buf = append(buf, '"')
			return buf, nil
		}
	case reflect.Bool:
		return func(buf []byte, rv reflect.Value) ([]byte, error) {
			if rv.Bool() {
				return append(buf, `"true"`...), nil
			}
			return append(buf, `"false"`...), nil
		}
	default:
		return inner
	}
}

// ================================
// Map Encoder
// ================================

// mapKV holds a key-value pair for sorted map encoding.
type mapKV struct {
	key string
	val reflect.Value
}

// mapKVPool pools []mapKV slices for map key sorting to reduce allocations.
var mapKVPool = sync.Pool{}

func buildMapEncoder(t reflect.Type) encoderFunc {
	if t.Key().Kind() != reflect.String {
		return func(buf []byte, rv reflect.Value) ([]byte, error) {
			return buf, fmt.Errorf("json: unsupported map key type %s", t.Key())
		}
	}
	valEnc := encoderForType(t.Elem())

	return func(buf []byte, rv reflect.Value) ([]byte, error) {
		if rv.IsNil() {
			return append(buf, "null"...), nil
		}

		buf = append(buf, '{')

		n := rv.Len()
		if n == 0 {
			buf = append(buf, '}')
			return buf, nil
		}

		// Get or create a kv slice from pool
		var pairs []mapKV
		if v := mapKVPool.Get(); v != nil {
			pairs = *v.(*[]mapKV)
			pairs = pairs[:0]
		}
		if cap(pairs) < n {
			pairs = make([]mapKV, 0, n)
		}

		// Collect key-value pairs in a single pass (avoids re-lookup)
		iter := rv.MapRange()
		for iter.Next() {
			pairs = append(pairs, mapKV{key: iter.Key().String(), val: iter.Value()})
		}
		sort.Slice(pairs, func(i, j int) bool {
			return pairs[i].key < pairs[j].key
		})

		for i := range pairs {
			if i > 0 {
				buf = append(buf, ',')
			}
			// Write key
			buf = append(buf, '"')
			buf = appendEscapedString(buf, pairs[i].key)
			buf = append(buf, '"', ':')

			// Write value (no re-lookup needed)
			var err error
			buf, err = valEnc(buf, pairs[i].val)
			if err != nil {
				// Clear refs before returning to pool
				for j := range pairs {
					pairs[j].val = reflect.Value{}
				}
				mapKVPool.Put(&pairs)
				return buf, err
			}
		}

		// Clear reflect.Value refs before returning to pool (avoid retaining references)
		for i := range pairs {
			pairs[i].val = reflect.Value{}
		}
		mapKVPool.Put(&pairs)

		buf = append(buf, '}')
		return buf, nil
	}
}

// ================================
// Slice / Array Encoders
// ================================

func buildSliceEncoder(t reflect.Type) encoderFunc {
	elemEnc := encoderForType(t.Elem())

	return func(buf []byte, rv reflect.Value) ([]byte, error) {
		if rv.IsNil() {
			return append(buf, "null"...), nil
		}

		buf = append(buf, '[')
		n := rv.Len()
		for i := 0; i < n; i++ {
			if i > 0 {
				buf = append(buf, ',')
			}
			var err error
			buf, err = elemEnc(buf, rv.Index(i))
			if err != nil {
				return buf, err
			}
		}
		buf = append(buf, ']')
		return buf, nil
	}
}

func buildArrayEncoder(t reflect.Type) encoderFunc {
	elemEnc := encoderForType(t.Elem())

	return func(buf []byte, rv reflect.Value) ([]byte, error) {
		buf = append(buf, '[')
		n := rv.Len()
		for i := 0; i < n; i++ {
			if i > 0 {
				buf = append(buf, ',')
			}
			var err error
			buf, err = elemEnc(buf, rv.Index(i))
			if err != nil {
				return buf, err
			}
		}
		buf = append(buf, ']')
		return buf, nil
	}
}

// ================================
// Error Encoder
// ================================

func unsupportedEnc(t reflect.Type) encoderFunc {
	return func(buf []byte, rv reflect.Value) ([]byte, error) {
		return buf, fmt.Errorf("json: unsupported type %s", t)
	}
}
