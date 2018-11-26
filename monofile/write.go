package monofile

import (
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func WriteFile(filename string, m Monofile) error {
	// Dump to YAML
	data, err := Dump(m)
	if err != nil {
		return err
	}
	// Write to file
	return ioutil.WriteFile(filename, data, 0644)
}

func Write(writer io.Writer, m Monofile) error {
	// Dump to YAML
	data, err := Dump(m)
	if err != nil {
		return err
	}
	// Write to output
	if _, err := writer.Write(data); err != nil {
		return err
	}
	return nil
}

// Lockfile

func WriteLockFile(filename string, m LockedMonofile) error {
	// Dump to YAML
	data, err := Dump(m)
	if err != nil {
		return err
	}
	// Write to file
	return ioutil.WriteFile(filename, data, 0644)
}

func WriteLock(writer io.Writer, m LockedMonofile) error {
	// Dump to YAML
	data, err := Dump(m)
	if err != nil {
		return err
	}
	// Write to output
	if _, err := writer.Write(data); err != nil {
		return err
	}
	return nil
}

// Dump returns a list of records in YAML
func Dump(v interface{}) ([]byte, error) {
	return yaml.Marshal(v)
}
