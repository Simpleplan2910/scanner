package scanner

import (
	"scanner/pkg/db"
	"testing"

	"github.com/stretchr/testify/assert"
)

var f = `
Once upon a time in a small village nestled in the lush countryside, there lived a young girl named Lily. She had always dreamed of exploring distant lands and embarking on grand adventures. But her days were filled with mundane chores, and her spirit longed for something more Lily.

One day, as Lily was gathering berries in the forest, she stumbled upon a hidden pathway. Curiosity sparked within her, and she decided to follow it. The path led her deep into the woods until she arrived at a peculiar old tree.

To her astonishment, the tree had a small door carved into its trunk. Lily hesitated for a moment before bravely opening it. Inside, she discovered a magical realm filled with enchanting creatures and breathtaking landscapes.

Lily's heart swelled with joy as she befriended talking animals, encountered mythical beings, and traversed sparkling waterfalls. She realized that her dreams had come true, and she was living her very own fairy tale.

Days turned into weeks, and weeks into months, but Lily's love for the magical realm never wavered. She became a revered explorer, sharing her tales with wide-eyed villagers who marveled at her extraordinary experiences.

As time went on, Lily's adventures became legends, inspiring generations to embrace their own dreams and seek the magic that lay hidden within their own hearts.

And so, in that small village, the name Lily became synonymous with courage, curiosity, and the boundless spirit of adventure.

The end.
`

func TestFindSubstring(t *testing.T) {
	assert := assert.New(t)

	type testTable []struct {
		name             string
		isExpectNilError bool
		substr           string
		str              string
		lines            []db.Line
	}
	tt := testTable{
		{
			name:             "Lily test cases",
			isExpectNilError: true,
			substr:           "Lily",
			str:              f,
			lines: []db.Line{
				{
					LineNum: 1,
					Indexes: []int{100, 278},
				},
				{
					LineNum: 3,
					Indexes: []int{12},
				},
				{
					LineNum: 5,
					Indexes: []int{70},
				},
				{
					LineNum: 7,
					Indexes: []int{0},
				},
				{
					LineNum: 9,
					Indexes: []int{51},
				},
				{
					LineNum: 11,
					Indexes: []int{17},
				},
				{
					LineNum: 13,
					Indexes: []int{40},
				},
			},
		},
		{
			name:             "short test case 1",
			isExpectNilError: true,
			substr:           "abc",
			str:              "abc dd abc ff",
			lines: []db.Line{
				{
					LineNum: 0,
					Indexes: []int{0, 7},
				},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			lines, err := findSubstring([]byte(tc.str), tc.substr)
			if tc.isExpectNilError {
				assert.Nil(err, "error should be nil")
			} else {
				assert.NotNil(err, "error shouldn't be nil")
			}

			assert.True(isEqual(lines, tc.lines), "lines is not equal")

		})
	}
}

func isEqual(l1, l2 []db.Line) bool {
	if len(l1) != len(l2) {
		return false
	}
	m1 := make(map[uint32][]int, len(l1))
	for _, v := range l1 {
		m1[v.LineNum] = v.Indexes
	}

	for _, v := range l2 {
		if _, ok := m1[v.LineNum]; !ok {
			return false
		}
		if !isEqualIntSlice(m1[v.LineNum], v.Indexes) {
			return false
		}
		delete(m1, v.LineNum)
	}
	return len(m1) == 0
}

func isEqualIntSlice(x, y []int) bool {
	if len(x) != len(y) {
		return false
	}
	diff := make(map[int]int, len(x))
	for _, _x := range x {
		diff[_x]++
	}
	for _, _y := range y {
		if _, ok := diff[_y]; !ok {
			return false
		}
		diff[_y] -= 1
		if diff[_y] == 0 {
			delete(diff, _y)
		}
	}
	return len(diff) == 0
}
