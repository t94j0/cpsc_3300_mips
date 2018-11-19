package machine

import (
	"fmt"
	"io"
	"os"
)

func (m *Machine) LoadFromStdin() error {
	return m.LoadFromReader(os.Stdin)
}

func (m *Machine) LoadFromReader(r io.Reader) error {
	m.memory = make([]uint32, 0)

	for {
		var line uint32
		_, err := fmt.Fscanf(r, "%x\n", &line)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		m.memory = append(m.memory, line)
	}

	return nil
}
