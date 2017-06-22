package models

type StringArray struct {
	array []string
}

func NewStringArray() StringArray {
	out := StringArray{}
	out.array = []string{}
	return out
}

func StringArrayFromSlice(slice []string) StringArray {
	out := StringArray{}
	out.array = slice
	return out
}

func (s *StringArray) Add(str string) {
	s.array = append(s.array, str)
}

func (s StringArray) Join(delimiter string) string {
	out := ""
	for i, v := range s.array {
		out += v
		if i != len(v) - 1 {
			out += delimiter
		}
	}
	return out
}
