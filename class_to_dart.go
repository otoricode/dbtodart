package main

import "fmt"

func classToDart(c Class) string {
	// import class
	s := "import 'dart:convert';\n\n"
	if c.HasList {
		s += "import 'package:collection/collection.dart';\n"
	}
	for _, v := range c.Import {
		s += "import '" + v + ".dart';"
	}
	// class
	s += fmt.Sprintf("class %s {\n", c.Name)
	// declaration
	for _, v := range c.Properties {
		s += "  " + fmt.Sprintln(v.Declaration)
	}
	s += "\n"
	// constructor
	s += fmt.Sprintf("  %s({\n", c.Name)
	for _, v := range c.Properties {
		s += "    " + v.ConstructorLine + "\n"
	}
	s += "  });\n\n"
	// checking null safety
	s += fmt.Sprintf("  factory %s.fromMap(Map<String, dynamic> map) {\n", c.Name)
	for _, v := range c.Properties {
		if v.IsNullable {
			continue
		}
		s += fmt.Sprintf("    if (map['%s'] == null) {\n", v.JsonKey)
		s += fmt.Sprintf("      throw '%s is required';\n", v.JsonKey)
		s += fmt.Sprintln("    }")
	}
	s += fmt.Sprintf("    return %s(\n", c.Name)
	for _, v := range c.Properties {
		s += "      " + v.FromMapLine + "\n"
	}
	s += "    );\n  }\n\n"
	// copyWith
	s += fmt.Sprintf("  %s copyWith({\n", c.Name)
	for _, v := range c.Properties {
		s += "    " + v.CopyWithParamLine + "\n"
	}
	s += fmt.Sprintf("  }) {\n    return %s(\n", c.Name)
	for _, v := range c.Properties {
		s += "      " + v.CopyWithPassingLine + "\n"
	}
	s += "    );\n"
	s += "  }\n\n"
	// toMap
	s += "  Map<String, dynamic> toMap() {\n"
	s += "    return {\n"
	for _, v := range c.Properties {
		s += "      " + v.ToMapLine + "\n"
	}

	s += "    };\n"
	s += "  }\n\n"
	// toJson
	s += "  String toJson() => json.encode(toMap());\n\n"
	// fromJson
	s += fmt.Sprintf("  factory %s.fromJson(String source) =>%s.fromMap(json.decode(source));\n\n", c.Name, c.Name)
	// toString

	s += "  @override\n"
	s += "  String toString() {\n"
	s += fmt.Sprintf("    return '%s(", c.Name)
	for i, v := range c.Properties {
		if i != len(c.Properties)-1 {
			s += fmt.Sprintf(" %s,", v.ToStringLine)
		} else {
			s += fmt.Sprintf("%s)';\n", v.ToStringLine)
		}
	}
	s += " }\n\n"

	// equality
	s += "  @override"
	s += "  bool operator ==(Object other) {"
	s += "    if (identical(this, other)) return true;"
	s += "\n"
	s += fmt.Sprintf("    return other is %s &&\n", c.Name)
	for i, v := range c.Properties {
		if i < len(c.Properties)-1 {
			s += "        " + v.EqualityLine + "&&\n"

		} else {
			s += "        " + v.EqualityLine + ";\n"

		}
	}
	s += "  }\n\n"

	// hashCode
	s += "  @override\n"
	s += "  int get hashCode {\n"
	for i, v := range c.Properties {
		if i == 0 {
			s += "    return " + v.HashCodeLine + " ^\n"
		} else if i == len(c.Properties)-1 {
			s += "        " + v.HashCodeLine + ";\n"
		} else {
			s += "        " + v.HashCodeLine + " ^\n"

		}
	}
	s += "  }\n\n"

	// class
	s += "}"
	return s
}
