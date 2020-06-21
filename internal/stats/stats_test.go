package stats

import "testing"

func TestLaunchStats(t *testing.T) {

	tableTest := [][]int{
		//{[3,4,5], 14},
		{[5,6,7], 55}
	}

	var res int
	for _, test := range tableTest {
		res = LaunchStats(test[0])
		if res != test[1] {
			t.Fatal()
		}
		t.Logf("%d == %d\n", res, test[1])
	}

}
