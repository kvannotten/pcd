// Copyright © 2018 Kristof Vannotten <kristof@vannotten.be>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package cmd

import (
	"reflect"
	"testing"
)

func TestEpisodeRangeArgs(t *testing.T) {
	cases := make(map[string][]int)
	cases["1"] = []int{1}
	cases["2-5"] = []int{2, 3, 4, 5}
	cases["2-5,!4"] = []int{2, 3, 5}
	cases["1,2,!99,!102,97-105,!103"] = []int{1, 2, 97, 98, 100, 101, 104, 105}
	cases["!25,25-27,!26"] = []int{27}
	cases["22,12,10-13,!11,!22,!12"] = []int{10, 12, 13, 22}
	cases["101-106,7,!105,!104"] = []int{7, 101, 102, 103, 106}
	cases[""] = nil

	for arg, want := range cases {
		got, err := parseRangeArg(arg)
		if err != nil {
			t.Error(err)
		} else if reflect.DeepEqual(want, got) == false {
			t.Errorf("missmatch for %s: got %v want %v", arg, got, want)
		}
	}
}
