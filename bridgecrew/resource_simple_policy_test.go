package bridgecrew

import (
	"context"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Test_resourceSimplePolicy(t *testing.T) {
	tests := []struct {
		name string
		want *schema.Resource
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := resourceSimplePolicy(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("resourceSimplePolicy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_resourceSimplePolicyCreate(t *testing.T) {
	type args struct {
		ctx context.Context
		d   *schema.ResourceData
		m   interface{}
	}
	tests := []struct {
		name string
		args args
		want diag.Diagnostics
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := resourceSimplePolicyCreate(tt.args.ctx, tt.args.d, tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("resourceSimplePolicyCreate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_setSimplePolicy(t *testing.T) {
	type args struct {
		d *schema.ResourceData
	}
	tests := []struct {
		name    string
		args    args
		want    Policy
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := setSimplePolicy(tt.args.d)
			if (err != nil) != tt.wantErr {
				t.Errorf("setSimplePolicy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("setSimplePolicy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_setConditions(t *testing.T) {
	type args struct {
		d *schema.ResourceData
	}
	tests := []struct {
		name    string
		args    args
		want    []Conditions
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := setConditions(tt.args.d)
			if (err != nil) != tt.wantErr {
				t.Errorf("setConditions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("setConditions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_setBenchmark(t *testing.T) {
	type args struct {
		d *schema.ResourceData
	}
	tests := []struct {
		name string
		args args
		want Benchmark
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := setBenchmark(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("setBenchmark() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_resourceSimplePolicyRead(t *testing.T) {
	type args struct {
		ctx context.Context
		d   *schema.ResourceData
		m   interface{}
	}
	tests := []struct {
		name string
		args args
		want diag.Diagnostics
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := resourceSimplePolicyRead(tt.args.ctx, tt.args.d, tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("resourceSimplePolicyRead() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_resourceSimplePolicyUpdate(t *testing.T) {
	type args struct {
		ctx context.Context
		d   *schema.ResourceData
		m   interface{}
	}
	tests := []struct {
		name string
		args args
		want diag.Diagnostics
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := resourceSimplePolicyUpdate(tt.args.ctx, tt.args.d, tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("resourceSimplePolicyUpdate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_simplepolicyChange(t *testing.T) {
	type args struct {
		d *schema.ResourceData
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := simplepolicyChange(tt.args.d); got != tt.want {
				t.Errorf("simplepolicyChange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_resourceSimplePolicyDelete(t *testing.T) {
	type args struct {
		ctx context.Context
		d   *schema.ResourceData
		m   interface{}
	}
	tests := []struct {
		name string
		args args
		want diag.Diagnostics
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := resourceSimplePolicyDelete(tt.args.ctx, tt.args.d, tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("resourceSimplePolicyDelete() = %v, want %v", got, tt.want)
			}
		})
	}
}