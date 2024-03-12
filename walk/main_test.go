package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createTempDir(t *testing.T, files map[string]int) (dirname string, cleanup func()) {
	tempDir, err := os.MkdirTemp("", "walktest")
	if err != nil {
		t.Fatal(err)
	}
	for file, fileCount := range files {
		for j := 1; j <= fileCount; j++ {
			fname := fmt.Sprintf("file%d%s", j, file)
			fpath := filepath.Join(tempDir, fname)
			if err := os.WriteFile(fpath, []byte("dummy"), 0644); err != nil {
				t.Fatal(err)
			}
		}
	}
	return tempDir, func() { os.RemoveAll(tempDir) }
}

func TestRunArchive(t *testing.T) {
	testCases := []struct {
		name         string
		cfg          config
		extNoArchive string
		nArchive     int
		nNoArchive   int
	}{
		{name: "ArchiveExtentionNoMatch", cfg: config{ext: ".log"}, extNoArchive: ".gz", nArchive: 0, nNoArchive: 10},
		{name: "ArchiveExtentionMatch", cfg: config{ext: ".log"}, extNoArchive: "", nArchive: 10, nNoArchive: 0},
		{name: "ArchiveExtentionMixed", cfg: config{ext: ".log"}, extNoArchive: ".gz", nArchive: 5, nNoArchive: 5},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			buffer := bytes.Buffer{}
			tempDir, cleanup := createTempDir(t, map[string]int{
				testCase.cfg.ext:      testCase.nArchive,
				testCase.extNoArchive: testCase.nNoArchive,
			})
			defer cleanup()
			archiveDir, cleanupArchive := createTempDir(t, nil)
			defer cleanupArchive()
			testCase.cfg.archive = archiveDir
			assert.Nil(t, run(tempDir, &buffer, testCase.cfg))
			pattern := filepath.Join(tempDir, fmt.Sprintf("*%s", testCase.cfg.ext))
			expFiles, err := filepath.Glob(pattern)
			assert.Nil(t, err)
			expOut := strings.Join(expFiles, "\n")
			res := strings.TrimSpace(buffer.String())
			assert.Equal(t, expOut, res)

			filesArchived, err := os.ReadDir(archiveDir)
			assert.Nil(t, err)
			assert.Equal(t, testCase.nArchive, len(filesArchived))

		})
	}
}
func TestRunDelExtension(t *testing.T) {
	testCases := []struct {
		name        string
		cfg         config
		extNoDelete string
		nDelete     int
		nNoDelete   int
		expected    string
	}{
		{
			name:        "DeleteExtensionNoMatch",
			cfg:         config{ext: ".log", del: true},
			extNoDelete: ".gz",
			nDelete:     0,
			nNoDelete:   10,
			expected:    "",
		},
		{
			name:        "DeleteExtensionMatch",
			cfg:         config{ext: ".log", del: true},
			extNoDelete: "",
			nDelete:     10,
			nNoDelete:   0,
			expected:    "",
		},
		{
			name:        "DeleteExtensionMixed",
			cfg:         config{ext: ".log", del: true},
			extNoDelete: ".gz",
			nDelete:     5,
			nNoDelete:   5,
			expected:    "",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			buffer := bytes.Buffer{}
			logBuffer := bytes.Buffer{}
			tempDir, cleanup := createTempDir(t, map[string]int{
				testCase.cfg.ext:     testCase.nDelete,
				testCase.extNoDelete: testCase.nNoDelete,
			})
			defer cleanup()

			testCase.cfg.wLog = &logBuffer

			if err := run(tempDir, &buffer, testCase.cfg); err != nil {
				t.Fatal(err)
			}
			res := buffer.String()
			assert.Equal(t, testCase.expected, res)

			filesLeft, err := os.ReadDir(tempDir)
			assert.Nil(t, err)
			assert.Equal(t, testCase.nNoDelete, len(filesLeft))

			if testCase.cfg.del {
				lines := bytes.Split(logBuffer.Bytes(), []byte("\n"))
				assert.Equal(t, testCase.nDelete+1, len(lines))
			}
		})
	}
}
func TestFilterOut(t *testing.T) {
	testCases := []struct {
		name     string
		file     string
		ext      string
		minSize  int64
		expected bool
	}{
		{"FilterNotExtension", "testdata/dir.log", "", 0, false},
		{"FilterExtensionsMatch", "testdata/dir.log", ".log", 0, false},
		{"FilterExtensionNoMatch", "testdata/dir.log", ".sh", 0, true},
		{"FilterExtensionSizeMatch", "testdata/dir.log", ".log", 10, false},
		{"FilterExtensionSizeNoMatch", "testdata/dir.log", ".log", 20, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			info, err := os.Stat(tc.file)
			assert.Nil(t, err)
			f := filterOut(tc.file, tc.ext, tc.minSize, info)
			assert.Equal(t, tc.expected, f)
		})
	}
}

func TestRun(t *testing.T) {
	testCases := []struct {
		name     string
		root     string
		cfg      config
		expected string
	}{
		{"NoFilter", "testdata", config{ext: "", size: 0, list: true}, "testdata/dir.log\ntestdata/dir2/script.sh\n"},
		{"FilterExtensionMatch", "testdata", config{ext: ".log", size: 0, list: true}, "testdata/dir.log\n"},
		{"FilterExtensionSizeMatch", "testdata", config{ext: ".log", size: 10, list: true}, "testdata/dir.log\n"},
		{"FilterExtensionSizeNoMatch", "testdata", config{ext: ".log", size: 20, list: true}, ""},
		{"FilterExtensionNoMatch", "testdata", config{ext: ".gz", size: 0, list: true}, ""},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			buffer := bytes.Buffer{}
			err := run(testCase.root, &buffer, testCase.cfg)
			assert.Nil(t, err)
			assert.Equal(t, testCase.expected, buffer.String())
		})
	}
}

func TestArrays(t *testing.T) {
	// append
	x := []int{1, 2, 3}
	x = append(x, 19)
	x = append(x, x...)
	assert.Equal(t, x, []int{1, 2, 3, 19, 1, 2, 3, 19})
}

func TestCapacity(t *testing.T) {
	// capacity
	y := []int{1}
	assert.Equal(t, cap(y), 1)

	y = append(y, 1)
	assert.Equal(t, cap(y), 2)

	y = append(y, 1)
	assert.Equal(t, cap(y), 4)

	y = append(y, 1, 1)
	assert.Equal(t, len(y), 5)
	assert.Equal(t, cap(y), 8)
}
