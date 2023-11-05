package main

import (
	"fmt"
	"strings"
)

type Filter struct {
	Ids     []string `json:"ids,omitempty"`
	Kinds   []int    `json:"kinds,omitempty"`
	Authors []string `json:"authors,omitempty"`
	Tags    []string `json:"-,omitempty"`
	Since   int `json:"since,omitempty"`
	Until   int `json:"until,omitempty"`
	Limit   int `json:"limit,omitempty"`
}

type Filters []Filter

func (filter Filter) writeFilter(result *strings.Builder) {
	sectionAdded := false
	result.WriteString("{")
	if len(filter.Ids) > 0 {
		result.WriteString("\"ids\":")
		writeStrArrToJson(result, filter.Ids)
		sectionAdded = true
	}
	if len(filter.Kinds) > 0 {
		if sectionAdded {
			result.WriteString(",")
		}
		result.WriteString("\"kinds\":")
		writeIntArrToJson(result, filter.Kinds)
		sectionAdded = true
	}
	if len(filter.Authors) > 0 {
		if sectionAdded {
			result.WriteString(",")
		}
		result.WriteString("\"authors\":")
		writeStrArrToJson(result, filter.Authors)
		sectionAdded = true
	}
	//if len(filter.Tags) > 0 { // @MarkFix probably need to define FilterTag type?
	//	if sectionsAdded {
	//		result.WriteString(",")
	//	}
	//	result.WriteString("\"#e\":")
	//	//result.WriteString("\"#p\":")
	//	writeStrArrToJson(&result, filter.Tags))
	//	sectionsAdded = true
	//}
	if filter.Since > 0 {
		if sectionAdded {
			result.WriteString(",")
		}
		result.WriteString(fmt.Sprintf("\"since\":%d", filter.Since))
		sectionAdded = true
	}
	if filter.Until > 0 {
		if sectionAdded {
			result.WriteString(",")
		}
		result.WriteString(fmt.Sprintf("\"until\":%d", filter.Until))
		sectionAdded = true
	}
	if filter.Limit > 0 {
		if sectionAdded {
			result.WriteString(",")
		}
		result.WriteString(fmt.Sprintf("\"limit\":%d", filter.Limit))
	}
	result.WriteString("}")
}

func (filters Filters) ToJson() string {
	var result strings.Builder

	for i, filter := range filters {
		if i > 1 {
			result.WriteString(",")
		}
		filter.writeFilter(&result)
	}

	DevBuildValidJson(result.String())

	return result.String()
}
