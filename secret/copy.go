package secret

import (
	"bufio"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	marecmd "github.com/femnad/mare/cmd"
)

const (
	bufferSize      = 8192
	concealedType   = "CONCEALED"
	defaultCategory = "LOGIN"
	emailKey        = "email"
	passwordId      = "password"
	regularType     = "STRING"
	tempDir         = "/dev/shm"
	tempPattern     = "pws"
	usernameId      = "username"
)

type secret struct {
	password string
	data     map[string]string
}

type field struct {
	Id      string `json:"id"`
	Type    string `json:"type"`
	Purpose string `json:"purpose"`
	Label   string `json:"label"`
	Value   string `json:"value"`
}

type template struct {
	Title    string `json:"title"`
	Category string `json:"category"`
	Fields   []field
}

func shred(filename string) error {
	file, err := os.OpenFile(filename, os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	for offset := int64(0); offset < info.Size(); offset += bufferSize {
		buffer := make([]byte, bufferSize)
		_, err = rand.Read(buffer)
		if err != nil {
			return err
		}

		_, err = file.WriteAt(buffer, offset)
		if err != nil && err != io.EOF {
			return err
		}
	}

	err = file.Truncate(0)
	if err != nil {
		return err
	}

	return os.Remove(filename)
}

func Copy(secretName string) error {
	in := marecmd.Input{Command: fmt.Sprintf("pass %s", secretName)}
	out, err := marecmd.RunOutBuffer(in)
	if err != nil {
		return err
	}

	data := make(map[string]string)
	var parsedFirstLine bool
	var s secret

	scanner := bufio.NewScanner(&out.Stdout)
	for scanner.Scan() {
		line := scanner.Text()

		if !parsedFirstLine {
			s.password = line
			parsedFirstLine = true
			continue
		}

		fields := strings.Split(line, ": ")
		if len(fields) != 2 {
			return fmt.Errorf("unexpected number of field for line %s", line)
		}

		data[fields[0]] = fields[1]
	}
	s.data = data

	file, err := os.CreateTemp(tempDir, tempPattern)
	if err != nil {
		return err
	}
	defer func() {
		shredErr := shred(file.Name())
		if shredErr != nil {
			log.Fatalf("error shredding temp file %v", shredErr)
		}
	}()

	var fields []field
	fields = append(fields, field{
		Id:      "password",
		Type:    concealedType,
		Purpose: strings.ToUpper(passwordId),
		Label:   passwordId,
		Value:   s.password,
	})
	username, ok := data[usernameId]
	if !ok {
		email, exists := data[emailKey]
		if exists {
			username = email
			delete(data, emailKey)
		}
	}

	if username != "" {
		fields = append(fields, field{
			Id:      usernameId,
			Type:    regularType,
			Purpose: strings.ToUpper(usernameId),
			Label:   usernameId,
			Value:   username,
		})
	}

	cmd := []string{"op", "item", "create", "--template", file.Name()}
	for k, v := range data {
		if k == usernameId {
			continue
		}

		cmd = append(cmd, fmt.Sprintf("%s=%s", k, v))
		continue
	}

	t := template{
		Title:    secretName,
		Category: defaultCategory,
		Fields:   fields,
	}

	enc := json.NewEncoder(file)
	err = enc.Encode(t)
	if err != nil {
		return err
	}

	in = marecmd.Input{CmdSlice: cmd}
	return marecmd.RunErrOnly(in)
}
