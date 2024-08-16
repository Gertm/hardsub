package main

import "testing"

func Test_sendNotification(t *testing.T) {
	type args struct {
		msg    string
		title  string
		config Config
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"simple",
			args{
				"Dit is de msg",
				"Hardsub",
				Config{},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if err := sendNotification(tt.args.msg, tt.args.title, &tt.args.config); (err != nil) != tt.wantErr {
			// 	t.Errorf("sendNotification() error = %v, wantErr %v", err, tt.wantErr)
			// }
		})
	}
}
