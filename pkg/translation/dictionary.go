// CopyRight all Right are Reserved

// Package translation this package is created for the
// sake of translating the language embedded in the Templates to needed language choise of the User
package translation

import "strings"

// DICTIONARY FOR SAVING DIFFERENT LANGUAGES
var DICTIONARY = map[string]map[string]string{
	"amh": map[string]string{
		"login":                "giba",
		"logout":               "Wuta",
		"register":             "temezgeb",
		"full name":            "mulu sim",
		"sex":                  "tsota",
		"age":                  "edme",
		"acadamic status":      "yetmhrt dereja",
		"language":             "quankua",
		"city":                 "ketema",
		"kebele":               "kebele",
		"phone":                "silk kutr",
		"address":              "Adrasha",
		"region":               "kilil",
		"category":             "zerf",
		"id":                   "metawekia",
		"trainers id card":     "ye Seltagn metawekia",
		"disclaimer":           "mastawesha",
		"this id is valid for": "Yih Metawekia yemiageleglew le",
		"months only":          "bicha new ",
		"given date":           "yetesetebet ken",
	},
	"oromifa": map[string]string{
		"login":            "giba",
		"logout":           "chufa",
		"register":         "temezgeb",
		"Full name":        "mulu sim",
		"Sex":              "tsota",
		"Age":              "edme",
		"Acadamic Status":  "yetmhrt dereja",
		"Language":         "quankua",
		"City":             "ketema",
		"Kebele":           "kebele",
		"Phone":            "silk kutr",
		"address":          "Adrasha",
		"region":           "kilil",
		"category":         "zerf",
		"id":               "metawekia",
		"trainers id card": "ye Seltagn metawekia",
	},
}

// Translate  function to change the word to the needed Language Representation
func Translate(lang string, sentence string) string {
	switch strings.ToLower(lang) {
	case "en", "eng":
		return sentence
	case "amh", "am", "amharic", "amhara":
		return strings.ToTitle((DICTIONARY["amh"])[strings.ToLower(sentence)])
	case "oro", "or", "oromifa", "oromo":
		return strings.ToTitle((DICTIONARY["oromifa"])[strings.ToLower(sentence)])
	}
	return sentence
}
