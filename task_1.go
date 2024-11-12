package main

import (
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

var file, _ = os.Create("buf.txt")

// параметры в качестве входа, файлы сохранить все необходимые данные кривых
// откуда знать, что все перебрали
// показывать все шаги и проверять все шаги
// нужно хранить все перебранные варианты простых чисел елки-палки это же елементарно до 20
func main() {
	defer file.Close()
	var n, m int
	fmt.Scan(&n)
	fmt.Scan(&m)
	file.WriteString("\n" + "n = " + strconv.Itoa(n))
	file.WriteString("\n" + "m = " + strconv.Itoa(m))

	one_big := big.NewInt(1)
	list_prime_num := gen_prim_num(n, one_big)
	if len(list_prime_num) == 0 {
		file.WriteString("\n" + "Нет числа заданной длинны")
		fmt.Println("Нет числа заданной длинны")
		return
	}

	var i_num, a, b, n_big, r, order *big.Int
	for i := 0; ; i++ {
		defer func() {
			if err := recover(); err != nil {
				file.WriteString("\n" + "На этапе проверки условий N = 1 + T + p, не нашлось подходящих простых чисел")
				return
			}
		}()
		i_num = list_prime_num[i]
		fmt.Println(list_prime_num[0])
		a, b, _ = decomposeToSquares(i_num)

		n_big, r, order = test_2_conditions(a, b, one_big, i_num)
		if n_big != nil && r != nil && order != nil {
			break
		}
	}

	file.WriteString("\n" + "Проверка условий N = 1 + T + p")
	file.WriteString("\n" + "N, r, k = " + n_big.String() + ", " + r.String() + ", " + order.String())
	var prov string

	if check_m(r, i_num, one_big, m) {
		prov = "проходит"
	} else {
		prov = "не проходит"
	}
	file.WriteString("\n" + "Результат: " + prov)

	a_cof := foo(n_big, order)
	x, y := find_points(a_cof, i_num)
	rx, ry := scalarMultiply(x, y, n_big, a_cof, i_num)
	fmt.Printf("Результат умножения точки на скаляр: (%s, %s)\n", rx.String(), ry.String())
	rx, ry = scalarMultiply(x, y, order, a_cof, i_num)
	fmt.Printf("Результат умножения точки на скаляр: (%s, %s)\n", rx.String(), ry.String())
}

func gen_prim_num(n int, one_big *big.Int) []*big.Int {
	k := 15
	i_num := gen_n_bit(n, "0")
	for_big := big.NewInt(4)
	sum := new(big.Int)
	var numPrime []*big.Int

	for {
		sum.Add(i_num, one_big)
		prime := isPrime(i_num, k)

		if prime && (new(big.Int).Mod(i_num, for_big).Cmp(one_big) == 0) {
			if !given_length(n, i_num) {
				break
			}
			numPrime = append(numPrime, new(big.Int).Set(i_num))
		}
		i_num.Set(sum)
	}

	return numPrime
}

// Функция для чтения последней строки из файла
func find_points(a_cof, p *big.Int) (*big.Int, *big.Int) {
	x := big.NewInt(1)
	y := big.NewInt(1)
	for x.Cmp(p) != 0 {
		for y.Cmp(p) != 0 {
			if new(big.Int).Div(new(big.Int).Sub(new(big.Int).Mul(y, y), new(big.Int).Mul(new(big.Int).Mul(x, x), x)), x).Cmp(a_cof) == 0 {
				fmt.Println(x)
				fmt.Println(y)
				return x, y
			}
			y.Add(y, big.NewInt(1))
		}
		x.Add(x, big.NewInt(1))
	}
	return nil, nil
}

// Символ Якоби
func jacobi(a, n *big.Int) int {
	if n.Cmp(big.NewInt(0)) <= 0 {
		fmt.Println(n)
		return 0
	}
	j := 1
	aMod := new(big.Int).Mod(a, n)
	for aMod.Cmp(big.NewInt(0)) != 0 {
		for aMod.Bit(0) == 0 {
			aMod.Rsh(aMod, 1)
			nMod8 := new(big.Int).Mod(n, big.NewInt(8)).Int64()
			if nMod8 == 3 || nMod8 == 5 {
				j = -j
			}
		}
		fmt.Println("a, n", n, aMod)
		aMod, n = n, aMod
		fmt.Println("a, n", n, aMod)
		if aMod.Mod(aMod, big.NewInt(4)).Cmp(big.NewInt(3)) == 0 &&
			n.Mod(n, big.NewInt(4)).Cmp(big.NewInt(3)) == 0 {
			j = -j
		}
		aMod.Mod(aMod, n)
	}
	if n.Cmp(big.NewInt(1)) == 0 {
		return j
	}
	return 0
}

// Вычет ли
func isQuadraticResidueJacobi(a, p *big.Int) bool {
	return jacobi(a, p) == 1
}

// Невычет ли
func isQuadraticNonResidue(a, p *big.Int) bool {
	return jacobi(a, p) == -1
}

// Функция для поиска квадратичного вычета по модулю p (составного числа)
func findQuadraticResidueComposite(p, order *big.Int) (*big.Int, error) {
	a := big.NewInt(2)
	if order.Cmp(big.NewInt(2)) == 0 {
		for a.Cmp(p) < 0 {
			if isQuadraticNonResidue(a, p) {
				file.WriteString("\n" + "A невычет для N = 2r " + "\n" + "A = " + a.String())
				return a, nil
			}
			a.Add(a, big.NewInt(1))
		}
		return a, nil
	} else {
		if order.Cmp(big.NewInt(4)) == 0 {
			for a.Cmp(p) < 0 {
				if isQuadraticResidueJacobi(a, p) {
					file.WriteString("\n" + "A вычет для N = 4r " + "\n" + "A = " + a.String())
					return a, nil
				}
				a.Add(a, big.NewInt(1))
			}
		} else {
			return nil, errors.New("Нет квадратичного вычета и невычета")
		}
	}
	return nil, nil

}

// Генерация коэффициента A вычета для 4r и невычета для 2r
func foo(n_big, order *big.Int) *big.Int {
	a, err := findQuadraticResidueComposite(n_big, order)
	if err != nil {
		file.WriteString("\n" + "Нет квадратичного вычета и невычета")
		return nil
	} else {
		fmt.Printf("A: %d\n", a)
		return a
	}
}

func check_m(r, i_num, one_big *big.Int, m int) bool {
	file.WriteString("\n" + "Проверка (p^i != 1)")
	p_pow_i := new(big.Int)
	if r.Cmp(i_num) != 0 {
		for i := 1; i <= m; i++ {
			i_big, _ := new(big.Int).SetString(strconv.Itoa(i), 10)
			if p_pow_i.Mod(new(big.Int).Xor(i_num, i_big), r).Cmp(new(big.Int).Mod(one_big, r)) != 0 {
				file.WriteString("\n" + p_pow_i.Mod(new(big.Int).Xor(i_num, i_big), r).String() + " != " + new(big.Int).Mod(one_big, r).String())
				return true
			} else {
				return false
			}
		}
	}

	return false
}

func test_2_conditions(a, b, one_big, p *big.Int) (*big.Int, *big.Int, *big.Int) {
	two_big, _ := new(big.Int).SetString("2", 10)
	n_1 := new(big.Int).Mul(a, two_big)
	n_2 := new(big.Int).Mul(b, two_big)
	var t_arr [2]big.Int = [2]big.Int{*n_1, *n_2}
	n_sum := new(big.Int)
	n_sub := new(big.Int)
	p_plus_one := new(big.Int).Add(p, one_big)
	// fmt.Println(p_plus_one)
	k := 10

	for _, t_el := range t_arr {
		n_sum.Add(p_plus_one, &t_el)
		n_sub.Sub(p_plus_one, &t_el)
		r, mult, flag := if_true(n_sub, two_big, k)
		if flag {
			return n_sub, r, mult
		}
		r, mult, flag = if_true(n_sum, two_big, k)
		if flag {
			return n_sub, r, mult
		} else {
			return nil, nil, nil
		}
	}
	return nil, nil, nil
}

func if_true(n_some, r *big.Int, k int) (*big.Int, *big.Int, bool) {
	prime_num := new(big.Int).Div(n_some, r)
	if isPrime(prime_num, k) {
		return prime_num, r, true
	} else {
		if isPrime(prime_num.Div(prime_num, r), k) {
			return prime_num, r.Mul(r, r), true
		} else {
			return nil, nil, false
		}
	}
}

func decomposeToSquares(p *big.Int) (a, b *big.Int, found bool) {
	fmt.Println(p)
	a = new(big.Int)
	b = new(big.Int)
	one := big.NewInt(1)

	// Проверка, что p ≡ 1 (mod 4)
	mod4 := new(big.Int).Mod(p, big.NewInt(4))
	if mod4.Cmp(one) != 0 {
		file.WriteString("\np !≡ 1 (mod 4)")
		return nil, nil, false

	}

	for a.SetInt64(1); a.Cmp(p) < 0; a.Add(a, one) {
		aSquared := new(big.Int).Mul(a, a)
		bSquared := new(big.Int).Sub(p, aSquared)
		if bSquared.Sign() < 0 {
			break
		}

		b.Sqrt(bSquared)
		if new(big.Int).Mul(b, b).Cmp(bSquared) == 0 {
			file.WriteString("\na, b = " + a.String() + ", " + b.String())
			return a, b, true
		}
	}

	return nil, nil, false
}

func gen_n_bit(n int, fill string) *big.Int {
	var prime_n []string
	for i := 0; i < n; i++ {
		prime_n = append(prime_n, fill)
	}
	prime_n[0] = "1"
	str_num := strings.Join(prime_n, "")
	i_num, _ := new(big.Int).SetString(str_num, 2)
	return i_num
}

func given_length(n int, i_num *big.Int) bool {
	var prime_new *big.Int = gen_n_bit(n, "1")
	if i_num.Cmp(prime_new) < 1 {
		file.WriteString("\n" + "p = " + i_num.String())
		return true
	} else {
		// file.WriteString("\n" + "Нет числа заданной длинны ")
		return false
	}

}

func modExp(base, exponent, modulus *big.Int) *big.Int {
	result := big.NewInt(1)
	base.Mod(base, modulus)
	exp := new(big.Int).Set(exponent)

	for exp.Cmp(big.NewInt(0)) > 0 {
		if new(big.Int).And(exp, big.NewInt(1)).Cmp(big.NewInt(1)) == 0 {
			result.Mul(result, base)
			result.Mod(result, modulus)
		}

		exp.Rsh(exp, 1)
		base.Mul(base, base)
		base.Mod(base, modulus)
	}
	return result
}

func millerRabinTest(n, d *big.Int, rng *rand.Rand) bool {
	a := big.NewInt(0).Rand(rng, new(big.Int).Sub(n, big.NewInt(4)))
	a.Add(a, big.NewInt(2))
	x := modExp(a, d, n)

	if x.Cmp(big.NewInt(1)) == 0 || x.Cmp(new(big.Int).Sub(n, big.NewInt(1))) == 0 {
		return true
	}

	for {
		x.Mul(x, x).Mod(x, n)

		if x.Cmp(big.NewInt(1)) == 0 {
			return false
		}

		if x.Cmp(new(big.Int).Sub(n, big.NewInt(1))) == 0 {
			return true
		}

		d.Mul(d, big.NewInt(2))
		if d.Cmp(new(big.Int).Sub(n, big.NewInt(1))) >= 0 {
			break
		}
	}

	return false
}

// Основная функция для проверки числа на простоту
func isPrime(n *big.Int, k int) bool {
	if n.Cmp(big.NewInt(2)) == 0 || n.Cmp(big.NewInt(3)) == 0 {
		return true
	}
	if n.Cmp(big.NewInt(1)) <= 0 || new(big.Int).Mod(n, big.NewInt(2)).Cmp(big.NewInt(0)) == 0 {
		return false
	}

	d := new(big.Int).Sub(n, big.NewInt(1))
	for new(big.Int).Mod(d, big.NewInt(2)).Cmp(big.NewInt(0)) == 0 {
		d.Div(d, big.NewInt(2))
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < k; i++ {

		if !millerRabinTest(n, new(big.Int).Set(d), rng) {
			return false
		}
	}
	return true
}
