// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package lookslike

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// pathComponentType indicates the type of pathComponent.
type pathComponentType int

const (
	// pcMapKey is the Type for map keys.
	pcMapKey pathComponentType = 1 + iota
	// pcSliceIdx is the Type for slice indices.
	pcSliceIdx
	//
	pcScalar
)

func (pct pathComponentType) String() string {
	if pct == pcMapKey {
		return "map"
	} else if pct == pcSliceIdx {
		return "slice"
	} else if pct == pcScalar {
		return "scalar"
	} else {
		// This should never happen, but we don't want to return an
		// error since that would unnecessarily complicate the fluid API
		return "<unknown>"
	}
}

// pathComponent structs represent one breadcrumb in a Path.
type pathComponent struct {
	Type  pathComponentType // One of pcMapKey or pcSliceIdx
	Key   string            // Populated for maps
	Index int               // Populated for slices
}

func (pc pathComponent) String() string {
	if pc.Type == pcSliceIdx {
		return fmt.Sprintf("[%d]", pc.Index)
	}
	return pc.Key
}

// Path represents the Path within a nested set of maps.
type Path []pathComponent

// ExtendSlice is used to add a new pathComponent of the pcSliceIdx type.
func (p Path) ExtendSlice(index int) Path {
	return p.Extend(
		pathComponent{pcSliceIdx, "", index},
	)
}

// ExtendMap adds a new pathComponent of the pcMapKey type.
func (p Path) ExtendMap(key string) Path {
	return p.Extend(
		pathComponent{pcMapKey, key, -1},
	)
}

func (p Path) Extend(pc pathComponent) Path {
	out := make(Path, len(p)+1)
	copy(out, p)
	out[len(p)] = pc
	return out
}

// Concat combines two paths into a new Path without modifying any existing paths.
func (p Path) Concat(other Path) Path {
	out := make(Path, 0, len(p)+len(other))
	out = append(out, p...)
	return append(out, other...)
}

func (p Path) String() string {
	out := make([]string, len(p))
	for idx, pc := range p {
		out[idx] = pc.String()
	}
	return strings.Join(out, ".")
}

// Last returns a pointer to the Last pathComponent in this Path. If the Path empty,
// a nil pointer is returned.
func (p Path) Last() *pathComponent {
	idx := len(p) - 1
	if idx < 0 {
		return nil
	}
	return &p[len(p)-1]
}

// GetFrom takes a map and fetches the given Path from it.
func (p Path) GetFrom(m interface{}) (value interface{}, exists bool) {
	value = m
	exists = true
	for _, pc := range p {
		rt := reflect.TypeOf(value)
		switch rt.Kind() {
		case reflect.Map:
			converted := interfaceToMap(value)
			value, exists = converted[pc.Key]
		case reflect.Slice:
			converted := sliceToSliceOfInterfaces(value)
			if pc.Index < len(converted) {
				exists = true
				value = converted[pc.Index]
			} else {
				exists = false
				value = nil
			}
		default:
			// If this case has been reached this means the expected type, say a map,
			// is actually something else, like a string or an array. In this case we
			// simply say the value doesn't exist. From a practical perspective this is
			// the right behavior since it will cause validation to fail.
			return nil, false
		}

		if exists == false {
			return nil, exists
		}
	}

	return value, exists
}

var arrMatcher = regexp.MustCompile("\\[(\\d+)\\]")

// InvalidPathString is the error type returned from unparseable paths.
type InvalidPathString string

func (ps InvalidPathString) Error() string {
	return fmt.Sprintf("Invalid Path: %#v", ps)
}

// ParsePath parses a Path of form key.[0].otherKey.[1] into a Path object.
func ParsePath(in string) (p Path, err error) {
	keyParts := strings.Split(in, ".")

	p = make(Path, len(keyParts))
	for idx, part := range keyParts {
		r := arrMatcher.FindStringSubmatch(part)
		pc := pathComponent{Index: -1}
		if len(r) > 0 {
			pc.Type = pcSliceIdx
			// Cannot fail, validated by regexp already
			pc.Index, err = strconv.Atoi(r[1])
			if err != nil {
				return p, err
			}
		} else if len(part) > 0 {
			pc.Type = pcMapKey
			pc.Key = part
		} else {
			return nil, InvalidPathString(in)
		}

		p[idx] = pc
	}

	return p, nil
}

// MustParsePath is a convenience method for parsing paths that have been previously validated
func MustParsePath(in string) Path {
	out, err := ParsePath(in)
	if err != nil {
		panic(err)
	}
	return out
}
