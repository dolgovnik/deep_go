package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

func forEach(e *entity, f func(int, int)){
	if e == nil {
		return
	}
	forEach(e.left, f)
	f(e.k, e.v)
	forEach(e.right, f)
}

func minSearch(curr *entity) *entity {
	if curr.left == nil {
		return curr
	}
	return minSearch(curr.left)
}

func del(curr *entity, key int) (*entity, bool) {
	// если не нашли ключ, то ничего не удаляем
	if curr == nil {
		return nil, false
	}

	deleted := false

	// нашли ключ
	if key == curr.k {
		// и он - лист - просто удаляем
		if curr.left == nil && curr.right == nil {
			return  nil, true
		}
		// у него есть левый - ставим этот левый вместо искомого
		if curr.right == nil {
			return curr.left, true
		}
		// у него есть правый - ставим этот правый вместо искомого
		if curr.left == nil {
			return curr.right, true
		}

		// есть оба - ищем минимальный, ставим его вместо искомого, затем удаляем минимальный
		change := minSearch(curr.right)
		curr.k = change.k
		curr.v = change.v
		curr.right, deleted = del(curr.right, change.k)

	// ключ больше текущего - идем вправо
	} else if key > curr.k {
		curr.right, deleted = del(curr.right, key)
	// ключ меньше текущего - идем влево
	} else if key < curr.k {
		curr.left, deleted = del(curr.left, key)
	}

	return curr, deleted
}

type entity struct {
	k int
	v int
	left *entity
	right *entity
}

func newEntity(key int, value int) *entity{
	return &entity{
			k: key,
			v: value,
			}
}

type OrderedMap struct {
	root *entity
	size int
}

func NewOrderedMap() OrderedMap {
	return OrderedMap{} // need to implement
}

func (m *OrderedMap) Insert(key, value int) {
	if m.root == nil {
		m.root = newEntity(key, value)
		m.size++
		return
	}

	curr := m.root
	for curr != nil {
		if key == curr.k {
			curr.v = value
			return
		}
		if key > curr.k {
			if curr.right == nil {
				curr.right = newEntity(key, value)
				m.size++
				return
			} else {
				curr = curr.right
				continue
			}
		} else if key < curr.k {
			if curr.left == nil {
				curr.left = newEntity(key, value)
				m.size++
				return
			} else {
				curr = curr.left
				continue
			}
		}
	}
}

func (m *OrderedMap) Erase(key int) {
	deleted := false
	m.root, deleted = del(m.root, key)
	if deleted {
		m.size--
	}	
}

func (m *OrderedMap) Contains(key int) bool {
	curr :=  m.root
	for curr != nil {
		if key == curr.k {
			return true
		}
		if key > curr.k {
			curr = curr.right
			continue
		}
		if key < curr.k {
			curr = curr.left
			continue
		}
	}
	return false 
}

func (m *OrderedMap) Size() int {
	return m.size
}

func (m *OrderedMap) ForEach(action func(int, int)) {
	forEach(m.root, action)
}

func TestCircularQueue(t *testing.T) {
	data := NewOrderedMap()
	assert.Zero(t, data.Size())

	data.Insert(10, 10)
	data.Insert(5, 5)
	data.Insert(15, 15)
	data.Insert(2, 2)
	data.Insert(4, 4)
	data.Insert(12, 12)
	data.Insert(14, 14)

	assert.Equal(t, 7, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(3))
	assert.False(t, data.Contains(13))

	var keys []int
	expectedKeys := []int{2, 4, 5, 10, 12, 14, 15}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))

	data.Erase(15)
	data.Erase(14)
	data.Erase(2)

	assert.Equal(t, 4, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(2))
	assert.False(t, data.Contains(14))

	keys = nil
	expectedKeys = []int{4, 5, 10, 12}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))
}
