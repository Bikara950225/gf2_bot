package bot_cmd_handler

import (
	"reflect"
	"testing"
)

func Test_parseInput(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   []string
		wantErr bool
	}{
		{
			name: "error path1, 空的命令",
			args: args{
				input: "",
			},
			want:    "",
			want1:   nil,
			wantErr: true,
		},
		{
			name: "error path2, 非法格式",
			args: args{
				input: "角色面板",
			},
			want:    "",
			want1:   nil,
			wantErr: true,
		},
		{
			name: "happy path1",
			args: args{
				input: "/角色面板 比优卡",
			},
			want:    "角色面板",
			want1:   []string{"比优卡"},
			wantErr: false,
		},
		{
			name: "happy path2",
			args: args{
				input: "/角色面板",
			},
			want:    "角色面板",
			want1:   []string{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := parseInput(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseInput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseInput() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("parseInput() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
