package hsconfig

import (
	"path/filepath"
	"testing"
)

const testConfig = "./testdata/config.json"

func Test_Config_AddToRecentProjects(t *testing.T) {
	config := generateDefaultConfig(testConfig)

	path1 := "/path/to/project1"
	path2 := "/path/to/project2"

	config.AddToRecentProjects(path1)

	if config.RecentProjects[0] != path1 {
		t.Fatal("unexpected project path was added")
	}

	// second path should be added on index 0
	config.AddToRecentProjects(path2)

	if config.RecentProjects[0] != path2 {
		t.Fatal("unexpected project path was added")
	}

	// but first path still should exist in recentProjects
	if config.RecentProjects[1] != path1 {
		t.Fatal("unexpected project path was added")
	}

	// when we'll try to add path1 second time, it should be moved to the 'top'
	config.AddToRecentProjects(path1)

	if config.RecentProjects[0] != path1 {
		t.Fatal("unexpected project path was added")
	}

	if config.RecentProjects[1] != path2 {
		t.Fatal("unexpected project path was added")
	}

	// there should be still one instance of path1
	if len(config.RecentProjects) != 2 {
		t.Fatal("Unexpected recentProjects len after tests")
	}
}

func Test_Config_GetAuxMPQs(t *testing.T) {
	config := generateDefaultConfig(testConfig)

	if len(config.GetAuxMPQs()) != 0 {
		t.Fatal("Wrong mpqs list len on start")
	}

	path, err := filepath.Abs("./testdata/")
	if err != nil {
		t.Error(err)
	}

	config.AuxiliaryMpqPath = path

	mpqs := config.GetAuxMPQs()

	if len(mpqs) != 2 {
		t.Fatal("Unexpected mpqs read")
	}

	for i := range mpqs {
		if mpqs[i] == "invalid.qpm" {
			t.Fatal("Unexpected mpq list read")
		}
	}
}
