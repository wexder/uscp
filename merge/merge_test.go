package merge_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wexder/uscp/merge"
)

func TestMerge2(t *testing.T) {
	type args struct {
		data  map[string]any
		patch map[string]any
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]any
		wantErr bool
	}{
		{
			name: "Simple merge",
			args: args{
				data: map[string]any{
					"key": "value",
				},
				patch: map[string]any{
					"key": "value2",
				},
			},
			want: map[string]any{
				"key": "value2",
			},
			wantErr: false,
		},
		{
			name: "Simple merge",
			args: args{
				data: map[string]any{
					"key": "value",
				},
				patch: map[string]any{
					"key":  "value2",
					"key2": 2,
				},
			},
			want: map[string]any{
				"key":  "value2",
				"key2": 2,
			},
			wantErr: false,
		},
		{
			name: "Merge inner map",
			args: args{
				data: map[string]any{
					"key": "value",
					"map": map[string]any{
						"deepKey": "deepValue",
					},
				},
				patch: map[string]any{
					"key": "value2",
					"map": map[string]any{
						"deepKey": "deepValue2",
					},
				},
			},
			want: map[string]any{
				"key": "value2",
				"map": map[string]any{
					"deepKey": "deepValue2",
				},
			},
			wantErr: false,
		},
		{
			name: "Merge inner map",
			args: args{
				data: map[string]any{
					"key": "value",
					"map": map[string]any{
						"deepKey": "deepValue",
						"map2": map[string]any{
							"deepKey2": "deepValue",
						},
					},
				},
				patch: map[string]any{
					"key": "value2",
					"map": map[string]any{
						"deepKey": "deepValue2",
					},
				},
			},
			want: map[string]any{
				"key": "value2",
				"map": map[string]any{
					"deepKey": "deepValue2",
					"map2": map[string]any{
						"deepKey2": "deepValue",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Merge inner map",
			args: args{
				data: map[string]any{
					"key": "value",
					"map": map[string]any{
						"deepKey": "deepValue",
					},
				},
				patch: map[string]any{
					"key": "value2",
					"map": map[string]any{
						"deepKey": "deepValue2",
						"map2": map[string]any{
							"deepKey2": "deepValue",
						},
					},
				},
			},
			want: map[string]any{
				"key": "value2",
				"map": map[string]any{
					"deepKey": "deepValue2",
					"map2": map[string]any{
						"deepKey2": "deepValue",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Merge inner array",
			args: args{
				data: map[string]any{
					"key":   "value",
					"array": []string{"str", "str2"},
				},
				patch: map[string]any{
					"key":   "value2",
					"array": []string{"str3", "str4"},
				},
			},
			want: map[string]any{
				"key":   "value2",
				"array": []string{"str3", "str4"},
			},
			wantErr: false,
		},
		{
			name: "Merge inner array",
			args: args{
				data: map[string]any{
					"key": "value",
					"array": []map[string]any{
						{
							"k": "str1",
						},
						{
							"k": " str2",
						},
					},
				},
				patch: map[string]any{
					"key": "value2",
					"array": []map[string]any{
						{
							"k": "str3",
						},
						{
							"k": " str4",
						},
					},
				},
			},
			want: map[string]any{
				"key": "value2",
				"array": []map[string]any{
					{
						"k": "str3",
					},
					{
						"k": " str4",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := merge.Merge(tt.args.data, tt.args.patch)
			if (err != nil) != tt.wantErr {
				t.Errorf("Merge2() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
