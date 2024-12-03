package main

import (
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

var file, _ = os.Create("buf.txt")
var file_2, _ = os.Create("group.txt")

func main() {
	defer file.Close()
	defer file_2.Close()
	var n, m int
	fmt.Scan(&n)

	fmt.Scan(&m)
	file.WriteString("\n" + "n = " + strconv.Itoa(n))
	file.WriteString("\n" + "m = " + strconv.Itoa(m))
	var list_p []*big.Int
	one_big := big.NewInt(1)
	list_prime_num, _ := parseBigIntArray(gen_prim_num(n, one_big))

	if len(list_prime_num) == 0 {
		fmt.Println("Нет числа заданной длинны")
		return
	}

	if n < 4 {
		file.WriteString("\n" + "Маленькое n")
		return
	}
	var i_num, a, b, n_big, r, order *big.Int
	for i := 0; ; i++ {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
				file.WriteString("\n" + "На этапе проверки условий N = 1 + T + p, не нашлось подходящих простых чисел")
				return
			}
		}()

		var list_prime_num_array []int
		var f int
		for {
			f = rand.Intn(len(list_prime_num))
			found := slices.Contains(list_prime_num_array, f)
			if !found {
				list_prime_num_array = append(list_prime_num_array, f)
			} else {
				break
			}
		}
		// надо сделать округление вверх
		file.WriteString("\nprime p = " + list_prime_num[f].String())
		list_p = append(list_p, list_prime_num[f])
		fmt.Println(list_p)
		i_num = list_prime_num[f]

		a, b, _ = decomposeToSquares(i_num)
		if a == nil || b == nil {
			continue
		}
		res1, res2, res3, flag := test_2_conditions(a, b, one_big, i_num)
		if flag {
			n_big, r, order = res1, res2, res3
			fmt.Println("последнее p = ", i_num)
			break
		} else {
			file.WriteString("\n" + "На этапе проверки условий N = 1 + T + p, не нашлось подходящих чисел")
		}
	}

	file.WriteString("\n" + "Проверка условий N = 1 + T + p")
	file.WriteString("\n" + "N, r, k = " + n_big.String() + ", " + r.String() + ", " + order.String())

	if check_m(r, i_num, one_big, m) {
		file.WriteString("\n" + "Результат: " + "проходит")
	} else {
		file.WriteString("\n" + "Результат: " + "не проходит")
		return
	}

	a_cof_list := gen_a_cof(n_big, order)

	var x, y, a_cof *big.Int
	for _, el := range a_cof_list {
		x, y = find_points(el, i_num)
		if x != nil && y != nil {
			a_cof = el
			break
		}
	}
	file.WriteString("\nКоэффициент A: " + a_cof.String())
	file.WriteString("\n" + "Сгенерированы следующие точки \nx, y = " + x.String() + ", " + y.String())
	rx, ry, _ := scalarMultiply(x, y, n_big, a_cof, i_num)
	fmt.Printf("Результат проверки бесконечности: (%s, %s)\n", rx.String(), ry.String())

	file.WriteString("\n" + "N, r, k = " + n_big.String() + ", " + r.String() + ", " + order.String())

	aboba := new(big.Int).Div(n_big, r)
	rx, ry, _ = scalarMultiply(x, y, aboba, a_cof, i_num)

	xy_string := "(" + rx.String() + ", " + ry.String() + ")"
	file.WriteString("\nТочка Q = " + xy_string)

	file.WriteString("\nКривая (P, A ,Q, r) = " + i_num.String() + ", " + a_cof.String() + ", " + xy_string + ", " + r.String())

	var arr_points_x string
	var arr_points_y string
	arr_points_x = arr_points_x + rx.String() + ","
	arr_points_y = arr_points_y + ry.String() + ","

	i := big.NewInt(2)
	for {
		if i.Cmp(new(big.Int).Add(r, big.NewInt(1))) == 1 {
			break
		}
		q_qx, q_qy, _ := scalarMultiply(rx, ry, i, a_cof, i_num)
		i = new(big.Int).Add(i, big.NewInt(1))
		arr_points_x = arr_points_x + q_qx.String() + ","
		arr_points_y = arr_points_y + q_qy.String() + ","
	}

	file_2.WriteString("\nЦиклическая группа " + "\n" + arr_points_x + "\n" + arr_points_y)
	file.WriteString("\nЦиклическая группа " + "\n" + arr_points_x + "\n" + arr_points_y)
}

func checkPointOnCurve(x, y, a, p *big.Int) bool {
	left := new(big.Int).Mul(y, y)
	left.Mod(left, p)

	right := new(big.Int).Mul(x, x)
	right.Mul(right, x)
	right.Add(right, new(big.Int).Mul(a, x))
	right.Mod(right, p)

	return left.Cmp(right) == 0
}

func parseBigIntArray(input string) ([]*big.Int, error) {
	parts := strings.Split(input, ",")
	result := make([]*big.Int, 0, len(parts))

	for _, part := range parts {

		part = strings.TrimSpace(part)

		num := new(big.Int)
		_, ok := num.SetString(part, 10)
		if !ok {
			return nil, fmt.Errorf("не удалось преобразовать %s в big.Int", part)
		}
		result = append(result, num)
	}

	return result, nil
}

func gen_prim_num(n int, one_big *big.Int) string {
	k := 10
	i_num := gen_n_bit(n, "0")
	for_big := big.NewInt(4)
	sum := new(big.Int)
	var numP string
	fmt.Println(i_num)
	for {
		if !given_length(n, i_num) {
			break
		}
		sum.Add(i_num, one_big)
		prime := isPrime(i_num, k)
		if prime && (new(big.Int).Mod(i_num, for_big).Cmp(one_big) == 0) {
			numP = numP + ", " + i_num.String()
		}
		i_num.Set(sum)
	}
	return numP[2:]
}

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

func findQuadraticResiduesAndNonResidues(p, order *big.Int) ([]*big.Int, []*big.Int) {
	residueMap := make(map[string]bool)
	var residues []*big.Int
	var nonResidues []*big.Int

	one := big.NewInt(1)

	for i := new(big.Int).Set(one); i.Cmp(p) == -1; i.Add(i, one) {
		square := new(big.Int).Exp(i, big.NewInt(2), p)
		squareStr := square.String()

		if !residueMap[squareStr] {
			residues = append(residues, new(big.Int).Set(square))
			residueMap[squareStr] = true
		}
	}

	for i := new(big.Int).Set(one); i.Cmp(p) == -1; i.Add(i, one) {
		if !residueMap[i.String()] {
			nonResidues = append(nonResidues, new(big.Int).Set(i))
		}
	}

	return residues, nonResidues
}

func gen_a_cof(n_big, order *big.Int) []*big.Int {
	var ans []*big.Int
	if order.Cmp(big.NewInt(2)) == 0 {
		_, ans = findQuadraticResiduesAndNonResidues(n_big, order)
	}
	if order.Cmp(big.NewInt(4)) == 0 {
		ans, _ = findQuadraticResiduesAndNonResidues(n_big, order)
	}
	if len(ans) == 0 {
		file.WriteString("\n" + "Нет квадратичного вычета и невычета")
		return nil
	} else {
		return ans
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

func test_2_conditions(a, b, one_big, p *big.Int) (*big.Int, *big.Int, *big.Int, bool) {
	two_big, _ := new(big.Int).SetString("2", 10)
	n_1 := new(big.Int).Mul(a, two_big)
	n_2 := new(big.Int).Mul(b, two_big)
	var t_arr [2]big.Int = [2]big.Int{*n_1, *n_2}
	n_sum := new(big.Int)
	n_sub := new(big.Int)
	p_plus_one := new(big.Int).Add(p, one_big)

	k := 10

	for _, t_el := range t_arr {
		n_sub.Sub(p_plus_one, &t_el)
		r, mult, flag := if_true(n_sub, two_big, k)
		if flag {
			return n_sub, r, mult, true
		}
		n_sum.Add(p_plus_one, &t_el)
		r, mult, flag = if_true(n_sum, two_big, k)
		if flag {
			return n_sum, r, mult, true
		} else {
			return nil, nil, nil, false
		}
	}
	return nil, nil, nil, false
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
	a = new(big.Int)
	b = new(big.Int)
	one := big.NewInt(1)

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
		return true
	} else {
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
