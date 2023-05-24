package scanner

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var f = `
Once upon a time in a small village nestled in the lush countryside, there lived a young girl named Lily. She had always dreamed of exploring distant lands and embarking on grand adventures. But her days were filled with mundane chores, and her spirit longed for something more.

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
	lines, err := findSubstring([]byte(f), "Lily")
	if err != nil {
		assert.Nil(err, "failed find ")
	}
	for _, v := range lines {
		t.Logf("indexes : %+v, line: %d", v.Indexes, v.LineNum)
	}
}
