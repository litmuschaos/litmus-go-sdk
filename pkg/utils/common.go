/*
Copyright © 2025 The LitmusChaos Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package utils

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"regexp"
	"strings"

	"github.com/litmuschaos/litmus-go-sdk/pkg/logger"
	"gopkg.in/yaml.v2"
)

func Scanner() string {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		return scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		logger.Errorf("reading standard input: %v", err)
	}
	return ""
}

func PrintError(err error) {
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func PrintInJsonFormat(inf interface{}) {
	var out bytes.Buffer
	byt, err := json.Marshal(inf)
	PrintError(err)

	err = json.Indent(&out, byt, "", "  ")
	PrintError(err)

	logger.Info(out.String())
}

func PrintInYamlFormat(inf interface{}) {
	byt, err := yaml.Marshal(inf)
	PrintError(err)

	logger.Info(string(byt))
}

func GenerateRandomString(n int) (string, error) {
	if n <= 0 {
		return "", fmt.Errorf("length should not be negative")
	}
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret), nil
}

func CheckKeyValueFormat(str string) bool {
	selectors := strings.Split(str, ",")

	for _, el := range selectors {
		kv := strings.Split(el, "=")
		if len(kv) != 2 {
			logger.Error("nodeselector is not correct. Correct format: \"key1=value2,key2=value2\"")
			return false
		}

		if strings.Contains(kv[0], "\"") || strings.Contains(kv[1], "\"") {
			logger.Error("nodeselector contains escape character(s). Correct format: \"key1=value2,key2=value2\"")
			return false
		}
	}
	return true
}

func GenerateNameID(in string) string {
	// Replace spaces and special characters with underscore
	reg := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	replaced := reg.ReplaceAllString(in, "_")

	// Remove hyphens
	noHyphens := strings.ReplaceAll(replaced, "-", "")

	// Convert everything to lowercase
	nameID := strings.ToLower(noHyphens)

	return nameID
}
