package monofile

import (
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func Read(reader io.Reader) (*Monofile, error) {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return Parse(data)
}

func ReadFile(filename string) (*Monofile, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return Parse(data)
}

func Parse(data []byte) (*Monofile, error) {
	m := Monofile{}
	err := yaml.Unmarshal(data, &m)
	if err != nil {
		return nil, err
	}

	return &m, nil
}

// Lockfle

func ReadLock(reader io.Reader) (*LockedMonofile, error) {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return ParseLock(data)
}

func ReadLockFile(filename string) (*LockedMonofile, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return ParseLock(data)
}

func ParseLock(data []byte) (*LockedMonofile, error) {
	m := LockedMonofile{}
	err := yaml.Unmarshal(data, &m)
	if err != nil {
		return nil, err
	}

	return &m, nil
}
