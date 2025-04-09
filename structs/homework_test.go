package main

import (
	"math"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

const (
	BuilderGamePersonType = iota
	BlacksmithGamePersonType
	WarriorGamePersonType
)

const (
	// Name
	nameSize         = 42
	stringEnd        = 0x00

	// Mana and Type
	typeShift        = 10
	manaMask int16   = 0b0000001111111111

	// Helath, House, Gun, Family
	hgfShift         = 10
	healthMask int16 = 0b0000001111111111
	house            = 1
	gun              = 2
	family           = 4

	// Level, Expirience, Strength, Respect
	levelShift       = 12
	expShift         = 8
	strShift         = 4
	mask             = 0b0000000000001111
)

type GamePerson struct {
	x    int32
	y    int32
	z    int32
	gold int32
	name [nameSize]byte

	// [----ttmmmmmmmmmm]
	typeAndMana             int16
    // [---FGHhhhhhhhhhh]
	healthHouseGunAndFamily int16
	// [lllleeeessssrrrr]
	levExpStrResp           uint16
}

func NewGamePerson(options ...Option) GamePerson {
	person := GamePerson{}
	for _, o := range options {
		o(&person)
	}
	return person
}

func (p *GamePerson) Name() string {
	nameLength := 0
	for i:=0; i<nameSize; i++ {
		if p.name[i] == stringEnd {
			break
		}
		nameLength++
	}
	return  unsafe.String(&p.name[0], nameLength)
}

func (p *GamePerson) X() int {
	return int(p.x)
}

func (p *GamePerson) Y() int {
	return int(p.y)
}

func (p *GamePerson) Z() int {
	return int(p.z)
}

func (p *GamePerson) Gold() int {
	return int(p.gold)
}

func (p *GamePerson) Mana() int {
	return int(p.typeAndMana & manaMask)
}

func (p *GamePerson) Health() int {
	return int(p.healthHouseGunAndFamily & healthMask)
}

func (p *GamePerson) Respect() int {
	return int(p.levExpStrResp & mask)
}

func (p *GamePerson) Strength() int {
	return int((p.levExpStrResp >> strShift) & mask)
}

func (p *GamePerson) Experience() int {
	return int((p.levExpStrResp >> expShift) & mask)
}

func (p *GamePerson) Level() int {
	return int((p.levExpStrResp >> levelShift) & mask)
}

func (p *GamePerson) HasHouse() bool {
	return (p.healthHouseGunAndFamily >> hgfShift) & house != 0
}

func (p *GamePerson) HasGun() bool {
	return (p.healthHouseGunAndFamily >> hgfShift) & gun != 0
}

func (p *GamePerson) HasFamilty() bool {
	return (p.healthHouseGunAndFamily >> hgfShift) & family != 0
}

func (p *GamePerson) Type() int {
	return int(p.typeAndMana >> typeShift)
}

type Option func(*GamePerson)

func WithName(name string) func(*GamePerson) {
	return func(person *GamePerson) {
		for i:=0; i<len(name); i++ {
			person.name[i] = byte(name[i])
		}
		if len(name) < nameSize {
			person.name[len(name)] = stringEnd
		}
	}
}

func WithCoordinates(x, y, z int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.x = int32(x)
		person.y = int32(y)
		person.z = int32(z)
	}
}

func WithGold(gold int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.gold = int32(gold)
	}
}

func WithMana(mana int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.typeAndMana |= int16(mana)
	}
}

func WithHealth(health int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.healthHouseGunAndFamily |= int16(health)
	}
}

func WithRespect(respect int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.levExpStrResp |= uint16(respect)
	}
}

func WithStrength(strength int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.levExpStrResp |= uint16(strength << strShift)
	}
}

func WithExperience(experience int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.levExpStrResp |= uint16(experience << expShift)
	}
}

func WithLevel(level int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.levExpStrResp |= uint16(level << levelShift)
	}
}

func WithHouse() func(*GamePerson) {
	return func(person *GamePerson) {
		person.healthHouseGunAndFamily |= int16(house << typeShift)
	}
}

func WithGun() func(*GamePerson) {
	return func(person *GamePerson) {
		person.healthHouseGunAndFamily |= int16(gun << typeShift)
	}
}

func WithFamily() func(*GamePerson) {
	return func(person *GamePerson) {
		person.healthHouseGunAndFamily |= int16(family << typeShift)
	}
}

func WithType(personType int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.typeAndMana |= int16(personType << typeShift)
	}
}


func TestGamePerson(t *testing.T) {
	assert.LessOrEqual(t, unsafe.Sizeof(GamePerson{}), uintptr(64))

	const x, y, z = math.MinInt32, math.MaxInt32, 0
	const name = "aaaaaaaaaaaaa_bbbbbbbbbbbbb_cccccccccccccc"
	const personType = BuilderGamePersonType
	const gold = math.MaxInt32
	const mana = 1000
	const health = 1000
	const respect = 10
	const strength = 10
	const experience = 10
	const level = 10

	options := []Option{
		WithName(name),
		WithCoordinates(x, y, z),
		WithGold(gold),
		WithMana(mana),
		WithHealth(health),
		WithRespect(respect),
		WithStrength(strength),
		WithExperience(experience),
		WithLevel(level),
		WithHouse(),
		WithFamily(),
		WithType(personType),
	}

	person := NewGamePerson(options...)
	assert.Equal(t, name, person.Name())
	assert.Equal(t, x, person.X())
	assert.Equal(t, y, person.Y())
	assert.Equal(t, z, person.Z())
	assert.Equal(t, gold, person.Gold())
	assert.Equal(t, mana, person.Mana())
	assert.Equal(t, health, person.Health())
	assert.Equal(t, respect, person.Respect())
	assert.Equal(t, strength, person.Strength())
	assert.Equal(t, experience, person.Experience())
	assert.Equal(t, level, person.Level())
	assert.True(t, person.HasHouse())
	assert.True(t, person.HasFamilty())
	assert.False(t, person.HasGun())
	assert.Equal(t, personType, person.Type())
}
