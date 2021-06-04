package gocalc

// func TestInterpreter(t *testing.T) {
// 	ass := assert.New(t)
// 	input := "2 + 2\n" +
// 		"(2 + 2) * 3\n" +
// 		"a = 10 + 1\n" +
// 		"b = a * 2\n" +
// 		"b\n" +
// 		"b - 5\n"
// 	expected := "= 4\n" +
// 		"= 12.000\n" +
// 		"= 22.000\n" +
// 		"= 17.000\n"
// 	buf := &strings.Builder{}
// 	ir := NewInterpreter(strings.NewReader(input), buf)

// 	err := ir.Start()
// 	ass.NoError(err)

// 	ass.Equal(expected, buf.String())
// }
