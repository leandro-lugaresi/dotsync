package steps

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestCopyConfigs_Run(t *testing.T) {
	_, err := ioutil.TempDir(os.TempDir(), "dotfiles")
	if err != nil {
		t.Fatal(err)
	}

}

var copyFileTests = []struct {
	src         string
	dst         string
	perm        os.FileMode
	expectError bool
}{
	{"step_copy_configs_test.go", "temp", 0644, false},
	{"step_copy_configs_test.go", "doesnotexist", 0640, false},
	{"step_copy_configs_test.go", "", 0660, true},
	{"", "temp", 0644, true},
	{"", "", 0644, true},
}

func TestCopyFile(t *testing.T) {
	for _, tt := range copyFileTests {
		dst := ""
		if tt.dst != "" {
			tempFile, err := ioutil.TempFile("", tt.dst)
			if err != nil {
				t.Fatalf("TempFile %s", err)
			}
			tempFile.Close()
			dst = tempFile.Name()
			// CopyFile should create the dst file if it does not exist.
			if tt.dst == "doesnotexist" {
				os.Remove(dst)
			}
			defer os.Remove(dst)
		}
		err := CopyFile(dst, tt.src, tt.perm)
		switch {
		case tt.expectError:
			if err == nil {
				t.Fatalf("CopyFile %s, %s: error expected, none found", dst, tt.src)
			}
			continue
		case !tt.expectError:
			if err != nil {
				t.Fatalf("CopyFile %s", err)
			}
		}
		dstFi, err := os.Stat(dst)
		if err != nil {
			t.Fatalf("Stat %s", err)
		}
		if dstFi.Mode() != tt.perm {
			t.Errorf("expected: %v, got %v", tt.perm, dstFi.Mode())
		}
		dstContents, err := ioutil.ReadFile(dst)
		if err != nil {
			t.Fatalf("ReadFile %s", err)
		}
		srcContents, err := ioutil.ReadFile(tt.src)
		if err != nil {
			t.Fatalf("ReadFile %s", err)
		}
		if string(dstContents) != string(srcContents) {
			t.Errorf("expected: %s, got %s", srcContents, dstContents)
		}
	}
}
