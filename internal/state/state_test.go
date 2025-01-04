package state

// func Test_StateMarshalJSON(t *testing.T) {
// 	type args struct {
// 		state *State
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    string
// 		wantErr bool
// 	}{
// 		{"zero", args{New()}, "{}", false},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := json.Marshal(tt.args.state)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("State.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			assert.Equal(t, tt.want, string(got))

// 		})
// 	}
// }
