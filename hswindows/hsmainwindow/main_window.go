package hsmainwindow

import "C"
import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
   "strconv"
   "os/exec"
   "errors"
	"unicode/utf16"
	"unicode/utf8"
	"encoding/json"

	"github.com/OpenDiablo2/HellSpawner/hsbuilder"
	"github.com/OpenDiablo2/HellSpawner/hswindows/hstextfilewindow"

	"github.com/OpenDiablo2/OpenDiablo2/d2common"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2mpq"
	"github.com/gotk3/gotk3/gtk"
)

//this was taken from old commit of OpenDiablo2
const (crcByteCount = 2)
type textDictionaryHashEntry struct {
	IsActive    bool
	Index       uint16
	HashValue   uint32
	IndexString uint32
	NameString  uint32
	NameLength  uint16
}

// MainWindow represents the main window of HellSpawner
type MainWindow struct {
	*gtk.ApplicationWindow
	treeView  *gtk.TreeView
	treeStore *gtk.TreeStore
	mpqs      []*d2mpq.MPQ
	mpqPaths  map[string]*gtk.TreeIter
    mpqResourceFilepaths map[string][]string
}

//this was taken from OpenDiablo2 and modified for our purposes
// LoadTextDictionary loads the text dictionary from the given data
func LoadTextDictionary(dictionaryData []byte) (map[string]string) {
    var lookupTable map[string]string

	if lookupTable == nil {
		lookupTable = make(map[string]string)
	}

	br := d2common.CreateStreamReader(dictionaryData)

	// skip past the CRC
	br.ReadBytes(crcByteCount)

	numberOfElements := br.GetUInt16()
	hashTableSize := br.GetUInt32()

	// Version (always 0)
	if _, err := br.ReadByte(); err != nil {
		log.Fatal("Error reading Version record")
	}

	br.GetUInt32() // StringOffset
	br.GetUInt32() // When the number of times you have missed a match with a hash key equals this value, you give up because it is not there.
	br.GetUInt32() // FileSize

	elementIndex := make([]uint16, numberOfElements)
	for i := 0; i < int(numberOfElements); i++ {
		elementIndex[i] = br.GetUInt16()
	}

	hashEntries := make([]textDictionaryHashEntry, hashTableSize)
	for i := 0; i < int(hashTableSize); i++ {
		hashEntries[i] = textDictionaryHashEntry{
			br.GetByte() == 1,
			br.GetUInt16(),
			br.GetUInt32(),
			br.GetUInt32(),
			br.GetUInt32(),
			br.GetUInt16(),
		}
	}

	for idx, hashEntry := range hashEntries {
		if !hashEntry.IsActive {
			continue
		}

		br.SetPosition(uint64(hashEntry.NameString))
		nameVal := br.ReadBytes(int(hashEntry.NameLength - 1))
		value := string(nameVal)

		br.SetPosition(uint64(hashEntry.IndexString))

		key := ""

		for {
			b := br.GetByte()
			if b == 0 {
				break
			}

			key += string(b)
		}

		if key == "x" || key == "X" {
			key = "#" + strconv.Itoa(idx)
		}

		_, exists := lookupTable[key]
		if !exists {
			lookupTable[key] = value
		}
	}

    return lookupTable
}

// Create creates a new instance of MainWindow
func Create(application *gtk.Application) (*MainWindow, error) {
	builder := hsbuilder.CreateBuilderFromTemplate(template)
	result := &MainWindow{
		ApplicationWindow: hsbuilder.ExtractApplicationWindow(builder, "mainApplicationWindow", application),
		mpqs:                 make([]*d2mpq.MPQ, 0),
		mpqPaths:             make(map[string]*gtk.TreeIter),
        mpqResourceFilepaths: make(map[string][]string, 0),
	}

	result.treeStore = hsbuilder.ExtractWidget(builder, "mpqTreeStore").(*gtk.TreeStore)
	result.treeView = hsbuilder.ExtractWidget(builder, "mainTreeView").(*gtk.TreeView)

	result.wireUpMenuHandlers(builder)

	_, _ = result.treeView.Connect("row-activated", func(treeView *gtk.TreeView,
		treePath *gtk.TreePath, column *gtk.TreeViewColumn) {
		iter, _ := result.treeStore.GetIter(treePath)
		val, _ := result.treeStore.GetValue(iter, 1)
		fileName, _ := val.GetString()
		result.handleFileActivated(fileName)
	})

	_, _ = result.Connect("destroy", func() { result.onWindowDestroyed() })

	return result, nil
}

// Append a row to the list store for the tree view
func (m *MainWindow) addRow(parent *gtk.TreeIter, file, path string) *gtk.TreeIter {
	// Get an iterator for a new row at the end of the list store
	iter := m.treeStore.Append(parent)

	// Set the contents of the list store row that the iterator represents
	_ = m.treeStore.SetValue(iter, 0, file)
	_ = m.treeStore.SetValue(iter, 1, path)

	return iter
}

func (m *MainWindow) onWindowDestroyed() {
	gtk.MainQuit()
}

func (m *MainWindow) wireUpMenuHandlers(builder *gtk.Builder) {
	miFileExit := hsbuilder.ExtractWidget(builder, "miFileExit").(*gtk.MenuItem)
	miFileAddExistingMPQ := hsbuilder.ExtractWidget(builder, "miFileAddExistingMPQ").(*gtk.MenuItem)

	_, _ = miFileExit.Connect("activate", func() { m.onFileExit() })
	_, _ = miFileAddExistingMPQ.Connect("activate", func() { m.onFileAddExistingMPQ() })
}

func (m *MainWindow) onFileExit() {
	gtk.MainQuit()
}

func (m *MainWindow) onFileAddExistingMPQ() {
	chooser, _ := gtk.FileChooserNativeDialogNew("Select MPQ(s)...", m, gtk.FILE_CHOOSER_ACTION_OPEN,
		"Open", "Cancel")

	fileFilter, _ := gtk.FileFilterNew()
	fileFilter.AddPattern("*.mpq")

	chooser.AddFilter(fileFilter)
	chooser.SetModal(true)
	chooser.SetSelectMultiple(true)

	if chooser.Run() == int(gtk.RESPONSE_ACCEPT) {
		fileNames, fileNamesErr := chooser.GetFilenames()

        if fileNamesErr != nil {
            fmt.Println("onFileAddExistingMPQ :: fileNamesErr :: ", fileNamesErr)
            log.Fatal(fileNamesErr)
        }

		for fileNameIdx := range fileNames {
			mpq, err := d2mpq.Load(fileNames[fileNameIdx])

			if err != nil {
				continue
			}

			m.mpqs = append(m.mpqs, mpq.(*d2mpq.MPQ))
			mpqFileName := filepath.Base(fileNames[fileNameIdx])
			mpqFiles := m.readMPQFiles(mpq.(*d2mpq.MPQ))

			for idx := range mpqFiles {                
				filePath := filepath.Clean(mpqFileName + "\\" + strings.ToLower(mpqFiles[idx]))
				parentNode := m.getFolderNode(filePath)
				fileParts := strings.Split(mpqFiles[idx], "\\")
				m.addRow(parentNode, fileParts[len(fileParts)-1], fileNames[fileNameIdx]+":"+mpqFiles[idx])
			}
		}
	}
}

func (m *MainWindow) buildResourcePath(inputMPQFile string, inputFilepath string) (string, error) {
    _, existsInMap := m.mpqResourceFilepaths[inputMPQFile]

	if !existsInMap {
        fmt.Println("main_window :: buildResourcePath :: unexpected :: inputMPQFile :: ", inputMPQFile)
        return "", errors.New("main_window :: buildResourcePath :: unexpected :: inputMPQFile :: "+inputMPQFile)
    }

    indexOfMatch := -1

    for i := 0; i < len(m.mpqResourceFilepaths[inputMPQFile]); i++ {
        if strings.EqualFold(m.mpqResourceFilepaths[inputMPQFile][i], strings.Replace(inputFilepath, "\\", "/", -1)) == true {
            indexOfMatch = i
            break
        }
    }

    if indexOfMatch == -1 {
        fmt.Println("main_window :: buildResourcePath :: unexpected :: inputFilepath :: ", inputFilepath)
        return "", errors.New("main_window :: buildResourcePath :: unexpected :: inputFilepath :: "+inputFilepath)
    }

    return m.mpqResourceFilepaths[inputMPQFile][indexOfMatch], nil
}

func (m *MainWindow) readMPQFiles(mpq *d2mpq.MPQ) []string {
	// Read listfile
	listfile, _ := mpq.GetFileList()
	// Search through using known contents
	s := bufio.NewScanner(strings.NewReader(rawListfile))

	for s.Scan() {
		if mpq.FileExists(s.Text()) {
			listfile = append(listfile, s.Text())
		}
	}

	fileMap := make(map[string]bool)
	var result []string

	for _, file := range listfile {
        var filenameBuffer strings.Builder
        m.addMPQResourcePath(mpq.Path(), strings.Replace(file, "\\", "/", -1))
        filenameSplitArr := strings.Split(strings.ToLower(file), "\\")

        for i := 0; i < len(filenameSplitArr); i++ {
            if i < (len(filenameSplitArr)-1) {
                filenameBuffer.WriteString(" ")
                filenameBuffer.WriteString(filenameSplitArr[i])
                filenameBuffer.WriteString("\\")
            } else {
                filenameBuffer.WriteString(filenameSplitArr[i])
            }
        }

        filename := filenameBuffer.String()
		_, ok := fileMap[filename]

		if !ok {
			fileMap[filename] = true
			result = append(result, filename)
		}
	}

	sort.Strings(result)
	return result
}

func (m *MainWindow) addMPQResourcePath(mpqFile string, path string) {
    _, existsInMap2 := m.mpqResourceFilepaths[mpqFile]

	if !existsInMap2 {
        var newArr []string
		m.mpqResourceFilepaths[mpqFile] = newArr
	}

    alreadyContains := false

    for i := 0; i < len(m.mpqResourceFilepaths[mpqFile]); i++ {
        if m.mpqResourceFilepaths[mpqFile][i] == path {
            alreadyContains = true
            break
        }
    }

    if alreadyContains == false {
        m.mpqResourceFilepaths[mpqFile] = append(m.mpqResourceFilepaths[mpqFile], path)
    }
}

func (m *MainWindow) getFolderNode(path string) *gtk.TreeIter {
	pathParts := strings.Split(path, "\\")
	pathParts = pathParts[:len(pathParts)-1]
	fullPath := ""
	parentPath := ""

	// Ensure folder structure
	for idx := range pathParts {
		fullPath += pathParts[idx] + "\\"
		_, ok := m.mpqPaths[fullPath]

		if !ok {
			m.mpqPaths[fullPath] = m.addRow(m.mpqPaths[parentPath], pathParts[idx], parentPath)
		}

		parentPath += pathParts[idx] + "\\"
	}

	return m.mpqPaths[fullPath]
}

func (m *MainWindow) handleFileActivated(name string) {
	parts := strings.Split(name, ":")

	if len(parts) != 2 {
		return
	}

	mpqPath := parts[0]
	filePath := parts[1]
    filePath = strings.Replace(filePath, " ", "", -1)
	fileExt := strings.ToLower(filepath.Ext(filePath))

	if fileExt == "" {
		return
	}

	switch fileExt {
	case ".txt":
		m.openTextFileWindow(mpqPath, filePath)
	case ".tbl":
		m.openTBLFileWindow(mpqPath, filePath)
    case ".dc6":
        actualFilepath, actualFilepathErr := m.buildResourcePath(mpqPath, filePath)

        if actualFilepathErr != nil {
            fmt.Println("main_window :: handleFileActivated :: actualFilepathErr :: ", actualFilepathErr)
            log.Fatal(actualFilepathErr)
        }
        
        m.spawnDC6FileViewer(mpqPath, filePath, actualFilepath)
	}

	log.Printf("Opening file for %s", fileExt)
}

func GetCallCommandArguments(command string) (string, []string) {
    commandParts := strings.Split(command, " ")
    var commandArgs []string

    for i := 1; i < len(commandParts); i++ {
        commandArgs = append(commandArgs, commandParts[i])
    }
    
    return commandParts[0], commandArgs
}

func (m *MainWindow) spawnDC6FileViewer(mpqPath, filePath, actualFilePath string) (error) {
    var commandBuffer strings.Builder
		
		if runtime.GOOS == "windows" {
			commandBuffer.WriteString(`dc6viewer.exe -mpq `)
		} else {
			commandBuffer.WriteString(`dc6viewer -mpq `)
		}
		
    commandBuffer.WriteString(mpqPath)
    commandBuffer.WriteString(` -asset `)
    commandBuffer.WriteString(actualFilePath)
    callName, callArgs := GetCallCommandArguments(commandBuffer.String())
    cmd := exec.Command(callName, callArgs...)
    err := cmd.Start()

    if err != nil {
        if strings.Contains(err.Error(), "exit status 1") == false {
            fmt.Println(err)
        }
    }
    
    return nil
}

func (m *MainWindow) openTBLFileWindow(mpqPath, filePath string) (error) {
	mpq := m.getMpqFromPath(mpqPath)

	if mpq == nil {
		return errors.New("main_window :: openTBLFileWindow :: unexpected mpq == nil")
	}

	data, dataErr := mpq.ReadFile(filePath)

	if dataErr != nil {
		log.Printf("Error reading file. dataErr :: ", dataErr)
        return dataErr
	}

	if len(data) == 0 {
		return errors.New("main_window :: openTBLFileWindow :: unexpected len(data) == 0")
	}

    strings := LoadTextDictionary(data)
	json, _ := json.MarshalIndent(strings, "", " ")
	window, windowErr := hstextfilewindow.Create(filePath, string(json))

    if windowErr != nil {
        return windowErr
    }

	window.ShowAll()
	window.ActivateFocus()
    return nil
}

func (m *MainWindow) openTextFileWindow(mpqPath, filePath string) (error) {
	mpq := m.getMpqFromPath(mpqPath)

	if mpq == nil {
        log.Printf("Error reading mpq.")
		return errors.New("main_window :: openTextFileWindow :: unexpected mpq == nil")
	}

	data, dataErr := mpq.ReadFile(filePath)

	if dataErr != nil {
		log.Printf("Error reading file. :: dataErr :: ", dataErr)
		return dataErr
	}

	if len(data) == 0 {
        return errors.New("main_window :: openTextFileWindow :: unexpected len(data) == 0")
	}

	textData := ""

	if data[0] == 255 && data[1] == 254 {
		// UTF16 apparently
        var textDataErr error
		textData, textDataErr = decodeUTF16(data)

        if textDataErr != nil {
            return textDataErr
        }
	} else {
		textData = string(data)
	}

	if textData == "" {
		return errors.New("main_window :: openTextFileWindow :: unexpected textData == \"\"")
	}

	window, windowErr := hstextfilewindow.Create(filePath, textData)

    if windowErr != nil {
        return windowErr
    }

	window.ShowAll()
	window.ActivateFocus()
    return nil
}

func decodeUTF16(b []byte) (string, error) {
	if len(b)%2 != 0 {
		return "", fmt.Errorf("Must have even length byte slice")
	}

	u16s := make([]uint1
6, 1)
	ret := &bytes.Buffer{}
	b8buf := make([]byte, 4)
	lb := len(b)
	for i := 0; i < lb; i += 2 {
		u16s[0] = uint16(b[i]) + (uint16(b[i+1]) << 8)
		r := utf16.Decode(u16s)
		n := utf8.EncodeRune(b8buf, r[0])
		ret.Write(b8buf[:n])
	}

	return ret.String(), nil
}

func (m *MainWindow) getMpqFromPath(mpqPath string) *d2mpq.MPQ {
	var mpq *d2mpq.MPQ

	for idx := range m.mpqs {
		if !strings.EqualFold(m.mpqs[idx].Path(), mpqPath) {
			continue
		}

		mpq = m.mpqs[idx]

		break
	}

	return mpq
}

const template = `
	<?xml version="1.0" encoding="UTF-8"?>
	<interface>
		<requires lib="gtk+" version="3.20"/>
		<object class="GtkTreeStore" id="mpqTreeStore">
			<columns>
				<column type="gchararray"/>
				<column type="gchararray"/>
			</columns>
		</object>
		<object class="GtkApplicationWindow" id="mainApplicationWindow">
			<property name="title" translatable="yes">HellSpawner Toolset</property>
			<property name="default-width">300</property>
			<property name="default-height">500</property>
			<child>
				<object class="GtkBox">
					<property name="orientation">vertical</property>
					<child>
						<object class="GtkMenuBar" id="mainMenuBar">
							<child>
								<object class="GtkMenuItem">
									<property name="label">File</property>
									<child type="submenu">
										<object class="GtkMenu">
											<child>
												<object class="GtkMenuItem" id="miFileAddExistingMPQ">
													<property name="label">Add Existing MPQ...</property>
												</object>
											</child>
											<child>
												<object class="GtkMenuItem" id="miFileExit">
													<property name="label">Exit</property>
												</object>
											</child>
										</object>
									</child>
								</object>
							</child>
						</object>
					</child>
					<child>
						<object class="GtkScrolledWindow">
							<child>
								<object class="GtkTreeView" id="mainTreeView">
									<property name="model">mpqTreeStore</property>
									<property name="enable-search">True</property>
									<property name="enable-tree-lines">False</property>
									<property name="headers-visible">False</property>
									<child>
										<object class="GtkTreeViewColumn" id="test-column">
											<property name="title">Name</property>
											<child>
												<object class="GtkCellRendererText" id="test-renderer"/>
												<attributes>
												<attribute name="text">0</attribute>
												</attributes>
											</child>
										</object>
									</child>
								</object>
							</child>
						</object>
						<packing>
							<property name="expand">True</property>
							<property name="fill">True</property>
						</packing>
					</child>
					<child>
						<object class="GtkHSeparator"></object>
					</child>
					<child>
						<object class="GtkToolbar">
							<child>
								<object class="GtkToolButton">
									<property name="icon-name">list-add</property>
								</object>
							</child>
							<child>
								<object class="GtkToolButton">
									<property name="icon-name">list-remove</property>
								</object>
							</child>
							<child>
								<object class="GtkSeparatorToolItem"></object>
							</child>
							<child>
								<object class="GtkToolButton">
									<property name="icon-name">document-properties</property>
								</object>
							</child>
						</object>
					</child>
				</object>
			</child>
		</object>
	</interface>
`
