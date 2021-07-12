package vgui

import (
    . "github.com/ahmetb/go-linq"
)

type NameFlagKey struct {
    Name       string
    Flags      []string
    FlagValues []bool
}

type NameFlagKeyPool struct {
    pool []*NameFlagKey
}

type Object struct {
    Name        string
    IsValue     bool
    Value       string
    Properties  map[string]*Object
    Flags       map[string]bool
    flagProps   map[*NameFlagKey]*Object
    flagKeyPool NameFlagKeyPool
}

func NewNameFlagKey(name string, flags map[string]bool) *NameFlagKey {
    var flagKeys []string
    var flagValues []bool
    for k, v := range flags {
        flagKeys = append(flagKeys, k)
        flagValues = append(flagValues, v)
    }
    return &NameFlagKey{name, flagKeys, flagValues}
}

func (k *NameFlagKey) Equal(other *NameFlagKey) bool {
    if other == nil {
        return false
    }
    return k.Name == other.Name && From(k.Flags).SequenceEqual(From(other.Flags)) &&
        From(k.FlagValues).SequenceEqual(From(other.FlagValues))
}

func (p *NameFlagKeyPool) GetOrCreate(name string, flags map[string]bool) *NameFlagKey {
    newKey := NewNameFlagKey(name, flags)
    for _, key := range p.pool {
        if newKey.Equal(key) {
            return key
        }
    }
    return newKey
}

func NewObject(name string) *Object {
    return &Object{name, false, "", map[string]*Object{}, map[string]bool{}, map[*NameFlagKey]*Object{}, NameFlagKeyPool{}}
}

func NewObjectFromValue(name string, value string) *Object {
    return &Object{name, true, value, map[string]*Object{}, map[string]bool{}, map[*NameFlagKey]*Object{}, NameFlagKeyPool{}}
}

func (o *Object) Get(name string) (*Object, bool) {
    obj, ok := o.Properties[name]
    return obj, ok
}

func (o *Object) mergeOrAddProperty(other *Object) {
    // if flagProps contains a matching NameFlagKey based on the iterated prop,
    // then merge. Otherwise, add to both Properties and flagProps
    key := other.getNameFlagKey()
    if prop, ok := o.flagProps[key]; ok && !prop.IsValue && !other.IsValue {
        prop.tryMerge(other)
    } else {
        o.flagProps[key] = other
        o.Properties[other.Name] = other
    }
}

func (o *Object) tryMerge(other *Object) bool {
    if o.Name == other.Name || !o.compareFlags(other) || o.IsValue || other.IsValue {
        return false
    }
    for _, prop := range other.Properties {
        o.mergeOrAddProperty(prop)
    }
    return true
}

func (o *Object) compareFlags(other *Object) bool {
    for flagKey, flagValue := range o.Flags {
        if value, ok := other.Flags[flagKey]; !ok || flagValue != value {
            return false
        }
    }
    return true
}

func (o *Object) getNameFlagKey() *NameFlagKey {
    return o.flagKeyPool.GetOrCreate(o.Name, o.Flags)
}