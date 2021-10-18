package main

import "testing"

func TestConvert(t *testing.T) {
	columns := []Column{
		{
			ColumnName: "id",
			DataType:   "bigint",
			IsNullable: "NO",
			ColumnKey:  "",
			IsList:     false,
		},
		{
			ColumnName: "remaining_uploads",
			DataType:   "bigint",
			IsNullable: "NO",
			ColumnKey:  "",
			IsList:     false,
		},
		{
			ColumnName: "users",
			DataType:   "users",
			IsNullable: "YES",
			ColumnKey:  "",
			IsList:     true,
		},

		{
			ColumnName: "created_at",
			DataType:   "datetime",
			IsNullable: "YES",
			ColumnKey:  "",
			IsList:     false,
		},
	}
	expecteds := []ClassProperties{
		{
			VarName:             "id",
			VarType:             "int",
			IsNullable:          false,
			Declaration:         "final int id;",
			ConstructorLine:     "required this.id,",
			JsonKey:             "id",
			FromMapLine:         "id: map['id'],",
			CopyWithParamLine:   "int? id,",
			CopyWithPassingLine: "id: id ?? this.id,",
			ToMapLine:           "'id': id,",
			ToStringLine:        "id: $id",
			EqualityLine:        "other.id == id",
			HashCodeLine:        "id.hashCode",
		},
		{
			VarName:             "remainingUpload",
			VarType:             "int",
			IsNullable:          false,
			Declaration:         "final int remainingUpload;",
			ConstructorLine:     "required this.remainingUpload,",
			JsonKey:             "remainingUpload",
			FromMapLine:         "remainingUpload: map['remainingUpload'],",
			CopyWithParamLine:   "int? remainingUpload,",
			CopyWithPassingLine: "remainingUpload: remainingUpload ?? this.remainingUpload,",
			ToMapLine:           "'remainingUpload': remainingUpload,",
			ToStringLine:        "remainingUpload: $remainingUpload",
			EqualityLine:        "other.remainingUpload == remainingUpload",
			HashCodeLine:        "remainingUpload.hashCode",
		},
		{
			VarName:             "users",
			VarType:             "Users",
			IsNullable:          true,
			Declaration:         "final List<Users>? users;",
			ConstructorLine:     "this.users,",
			JsonKey:             "users",
			FromMapLine:         "users: map['users'] == null ? null : List<Users>.from(map['users']?.map((x) => Users.fromMap(x))),",
			CopyWithParamLine:   "List<Users>? users,",
			CopyWithPassingLine: "users: users ?? this.users,",
			ToMapLine:           "'users': users?.map((x) => x.toMap()).toList(),",
			ToStringLine:        "users: $users",
			EqualityLine:        "DeepCollectionEquality().equals(other.users, users)",
			HashCodeLine:        "users.hashCode",
		},
		{
			VarName:             "createdAt",
			VarType:             "DateTime",
			IsNullable:          true,
			Declaration:         "final DateTime? createdAt;",
			ConstructorLine:     "this.createdAt,",
			JsonKey:             "createdAt",
			FromMapLine:         "createdAt: map['createdAt'] == null ? null : DateTime.parse(map['createdAt']),",
			CopyWithParamLine:   "DateTime? createdAt,",
			CopyWithPassingLine: "createdAt: createdAt ?? this.createdAt,",
			ToMapLine:           "'createdAt': createdAt?.toIso8601String(),",
			ToStringLine:        "createdAt: $createdAt",
			EqualityLine:        "other.createdAt == createdAt",
			HashCodeLine:        "createdAt.hashCode",
		},
	}
	if len(expecteds) != len(columns) {
		t.Error("len is not same")
		t.FailNow()
	}
	for i, v := range columns {
		prop := convertColumnToProperties(v)

		if expecteds[i].VarName != prop.VarName {
			t.Errorf("VarName expected: %s value : %s", expecteds[i].VarName, prop.VarName)
		}
		if expecteds[i].VarType != prop.VarType {
			t.Errorf("VarType expected: %s value : %s", expecteds[i].VarType, prop.VarType)
		}
		if expecteds[i].IsNullable != prop.IsNullable {
			t.Errorf("IsNullable expected: %v value : %v", expecteds[i].IsNullable, prop.IsNullable)
		}
		if expecteds[i].Declaration != prop.Declaration {
			t.Errorf("Declaration expected: %s value : %s", expecteds[i].Declaration, prop.Declaration)
		}
		if expecteds[i].ConstructorLine != prop.ConstructorLine {
			t.Errorf("ConstructorLine expected: %s value : %s", expecteds[i].ConstructorLine, prop.ConstructorLine)
		}
		if expecteds[i].JsonKey != prop.JsonKey {
			t.Errorf("JsonKey expected: %s value : %s", expecteds[i].JsonKey, prop.JsonKey)
		}
		if expecteds[i].FromMapLine != prop.FromMapLine {
			t.Errorf("FromMapLine expected: %s value : %s", expecteds[i].FromMapLine, prop.FromMapLine)
		}
		if expecteds[i].CopyWithParamLine != prop.CopyWithParamLine {
			t.Errorf("CopyWithParamLine expected: %s value : %s", expecteds[i].CopyWithParamLine, prop.CopyWithParamLine)
		}
		if expecteds[i].CopyWithPassingLine != prop.CopyWithPassingLine {
			t.Errorf("CopyWithPassingLine expected: %s value : %s", expecteds[i].CopyWithPassingLine, prop.CopyWithPassingLine)
		}
		if expecteds[i].ToMapLine != prop.ToMapLine {
			t.Errorf("ToMapLine expected: %s value : %s", expecteds[i].ToMapLine, prop.ToMapLine)
		}
		if expecteds[i].ToStringLine != prop.ToStringLine {
			t.Errorf("ToStringLine expected: %s value : %s", expecteds[i].ToStringLine, prop.ToStringLine)
		}
		if expecteds[i].EqualityLine != prop.EqualityLine {
			t.Errorf("EqualityLine expected: %s value : %s", expecteds[i].EqualityLine, prop.EqualityLine)
		}
		if expecteds[i].HashCodeLine != prop.HashCodeLine {
			t.Errorf("HashCodeLine expected: %s value : %s", expecteds[i].HashCodeLine, prop.HashCodeLine)
		}
	}
}
