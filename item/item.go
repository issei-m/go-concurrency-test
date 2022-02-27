package item

type Item int

func CreateItems(num int) []Item {
	items := make([]Item, num)
	for i := 0; i < num; i++ {
		items[i] = Item(i + 1)
	}

	return items
}
