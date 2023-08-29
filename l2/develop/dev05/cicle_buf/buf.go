package ciclebuf

type Buf struct {
	data  [][]byte
	start int // позиция от куда можно прочитать самый старый элемент
	end   int // позиция, куда можно записать новый элемент\
	len   int // количество элементов записанных
}

func NewBuf(n int) *Buf {
	return &Buf{
		data: make([][]byte, n),
	}
}

func (b *Buf) Push(arr []byte) {
	if len(b.data) == 0 {
		return
	}
	// Если нужно будет перезаписать существующий элемент
	// То мы перемещаем начало массива вперед
	if b.end == b.start && b.len != 0 {
		b.start = (b.start + 1) % len(b.data)
	} else {
		b.len++
	}

	b.data[b.end] = arr
	b.end = (b.end + 1) % len(b.data)
}

func (b *Buf) Pop() ([]byte, bool) {
	// В буфере нет элементов
	if b.len == 0 {
		return nil, false
	}

	res := b.data[b.start]                // Находим элемент
	b.start = (b.start + 1) % len(b.data) // Инкрементируемся
	b.len--

	return res, true
}
