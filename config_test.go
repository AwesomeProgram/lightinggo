package lightinggo

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

type configTestType struct {
	filename  string
	fn        EventOnChangeFunc
	option    ConfigOption
	wantfalse bool
	wanttrue  bool
}

var fnc EventOnChangeFunc

var configOption ConfigOption

var ConfigTestData map[string]configTestType

func TestMain(m *testing.M) {
	fmt.Println("Package config unit-test begin...") // 测试之前的做一些设置

	ConfigTestData = map[string]configTestType{
		"WithConfigFile": {
			filename:  "config.yaml",
			option:    configOption,
			wantfalse: false,
			wanttrue:  true,
		},
		"WithConfigWatcher": {
			fn:        fnc,
			option:    configOption,
			wantfalse: false,
			wanttrue:  true,
		},
		"LoadVariable": {
			option:    configOption,
			wantfalse: false,
			wanttrue:  true,
		},
	}
	// 如果 TestMain 使用了 flags，这里应该加上flag.Parse()
	retCode := m.Run()                              // 执行测试
	fmt.Println("Package config unit-test over...") // 测试之后做一些拆卸工作
	os.Exit(retCode)                                // 退出测试
}

func TestWithConfigFile(t *testing.T) {
	// 是否并行
	t.Parallel()
	type args struct {
		filename string
	}

	tests := []struct {
		name    string
		args    args
		want    ConfigOption
		wantErr bool
	}{
		{
			"normal",
			args{ConfigTestData["WithConfigFile"].filename},
			ConfigTestData["WithConfigFile"].option,
			ConfigTestData["WithConfigFile"].wantfalse,
		},
		{
			"abnormal",
			args{ConfigTestData["WithConfigFile"].filename},
			ConfigTestData["WithConfigFile"].option,
			ConfigTestData["WithConfigFile"].wanttrue,
		},
	}

	// teardownTestCase := setupTestCase(t)
	// defer teardownTestCase(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WithConfigFile(tt.args.filename)

			if reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf("WithConfigFile().valueType = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkWithConfigFile(b *testing.B) {
	filename := ConfigTestData["WithConfigFile"].filename
	for i := 0; i < b.N; i++ {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = WithConfigFile(filename)
			}
		})
	}
}

func TestWithConfigWatcher(t *testing.T) {
	t.Parallel()
	type args struct {
		fn EventOnChangeFunc
	}
	tests := []struct {
		name    string
		args    args
		want    ConfigOption
		wantErr bool
	}{
		{
			"normal",
			args{ConfigTestData["WithConfigWatchxer"].fn},
			ConfigTestData["WithConfigFile"].option,
			false,
		},
		{
			"abnormal", // directory is not exists.
			args{ConfigTestData["WithConfigWatchxer"].fn},
			ConfigTestData["WithConfigFile"].option,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WithConfigWatcher(tt.args.fn)

			if reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf("WithConfigWatcher() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkWithConfigWatcher(b *testing.B) {
	fn := ConfigTestData["WithConfigWatchxer"].fn
	for i := 0; i < b.N; i++ {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = WithConfigWatcher(fn)
			}
		})
	}
}

func TestLoadVariable(t *testing.T) {
	t.Parallel()
	type args struct {
		configOption ConfigOption
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			"normal1",
			args{WithConfigFile(ConfigTestData["WithConfigFile"].filename)},
			"",
			false,
		},
		{
			"normal2",
			args{WithConfigWatcher(ConfigTestData["WithConfigWatchxer"].fn)},
			"",
			false,
		},
		{
			"abnormal", // filepath is not exists.
			args{WithConfigWatcher(ConfigTestData["WithConfigWatchxer"].fn)},
			"",
			false,
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if i == 0 {
				err := LoadVariable()
				if (err != nil) != tt.wantErr {
					t.Errorf("LoadVariable() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			} else {
				err := LoadVariable(tt.args.configOption)
				if (err != nil) != tt.wantErr {
					t.Errorf("LoadVariable() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
		})
	}
}

func BenchmarkLoadVariable(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = LoadVariable(ConfigTestData["LoadVariable"].option)
			}
		})
	}
}
