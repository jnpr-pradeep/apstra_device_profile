package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

type Port struct {
	ConnectorType    string           `json:"connector_type"`
	PanelID          int              `json:"panel_id"`
	Transformations  []Transformation `json:"transformations"`
	ColumnID         int              `json:"column_id"`
	PortID           int              `json:"port_id"`
	RowID            int              `json:"row_id"`
	FailureDomainID  int              `json:"failure_domain_id"`
	DisplayID        int              `json:"display_id"`
	SlotID           int              `json:"slot_id"`
}

type Transformation struct {
	TransformationID int         `json:"transformation_id"`
	IsDefault        bool        `json:"is_default"`
	Interfaces       []Interface `json:"interfaces"`
}

type Interface struct {
	Name      string `json:"name"`
	InterfaceID int `json:"interface_id"`
	Speed     Speed  `json:"speed"`
	State     string `json:"state"`
	Setting   string `json:"setting"`
}

type Speed struct {
	Value int    `json:"value"`
	Unit  string `json:"unit"`
}

// Below PortData required for the internal processing
type PortData struct {
	Name				string `json:"name"`
	TransformationID	int		`json:"transformation_id"`
}

func processFixedChassisDevice(deviceProfileData map[string]interface{}) map[string]PortData {
	interfaceNamePartMap := map[string]PortData{}
	ports := deviceProfileData["ports"].([]interface{})
	for _, p := range ports {
		port := p.(map[string]interface{})
		transformations := port["transformations"].([]interface{})
		for _, t := range transformations {
			transformation := t.(map[string]interface{})
			interfaces := transformation["interfaces"].([]interface{})
			for _, i := range interfaces {
				ifc := i.(map[string]interface{})
				name := ifc["name"].(string)
				// interfaceID := ifc["interface_id"].(float64)
				speed := ifc["speed"].(map[string]interface{})
				speedValue := speed["value"].(float64)
				// speedUnit := speed["unit"].(string)
				transformationID := int(transformation["transformation_id"].(float64))
				ifcPartName := RemoveChannalisation(name)
				ifcSearchName := constructInterfaceString(ifcPartName, len(interfaces), int(speedValue))
				// fmt.Printf("Interface name: %s, InterfaceID: %v, Speed: %v %v, Transformation ID: %v\n", ifcPartName, interfaceID, speedValue, speedUnit, transformationID)
				_, ok := interfaceNamePartMap[ifcSearchName] 
				if !ok {
					interfaceNamePartMap[ifcSearchName] = PortData{Name: name, TransformationID: transformationID}
				}
			}
		}
	}
	return interfaceNamePartMap
}

// RemoveChannalisation removes the string after ":"
// It will take the input as "et-0/0/0:1" and returns output as -0/0/0"
func RemoveChannalisation(str2 string) string {
	if strings.Contains(str2, ":") {
        str2 = strings.Split(str2, ":")[0]
	}
	if strings.Contains(str2, "-") {
        str2 = strings.Split(str2, "-")[1]
	}
	// fmt.Printf("%v\n", str2)
	return str2
}

func constructInterfaceString(ifcPartName string, channalisationCount int, speed int) string {
	return fmt.Sprintf("%v$%v$%v", ifcPartName, channalisationCount, speed)
}

func main() {
	// data, err := ioutil.ReadFile("device_profile.json")
	data, err := ioutil.ReadFile("mc_device_profile.json")

	if err != nil {
		fmt.Printf("Error reading file: %v", err)
	}

	var dataMap map[string]interface{}
	err = json.Unmarshal(data, &dataMap)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	interfaceNamePartMap := processFixedChassisDevice(dataMap)
	// fmt.Println("InterfaceNamePartMap=%v", interfaceNamePartMap)

	ifcString := "1/0/0"
	// Search for 0/0/0, 1, 100
	searchIfc := constructInterfaceString(ifcString, 1, 100)
	fmt.Println("InterfaceNamePartMap[%v]=%v", searchIfc, interfaceNamePartMap[searchIfc])

	// Search for 0/0/0, 1, 40
	searchIfc = constructInterfaceString(ifcString, 1, 40)
	fmt.Println("InterfaceNamePartMap[%v]=%v", searchIfc, interfaceNamePartMap[searchIfc])

	// Search for 0/0/0, 2, 50
	searchIfc = constructInterfaceString(ifcString, 4, 10)
	fmt.Println("InterfaceNamePartMap[%v]=%v", searchIfc, interfaceNamePartMap[searchIfc])
}