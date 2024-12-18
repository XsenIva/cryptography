package main

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"sort"
	"strconv"
)

func qudro_sqrt(p int) int {
	p_float := float64(p * 2)
	k := int(math.Floor(math.Sqrt(math.Sqrt(p_float)))) + 1
	return k
}

func genPointOnCurve(a, p int) (int, int) {
	for {
		x := rand.Intn(p)

		right := (x*x*x + a*x) % p
		for y := 0; y < p; y++ {
			if (y*y)%p == right {
				return x, y
			}
		}
	}
}

func isInfinity(x, y *big.Int) bool {
	return x == nil && y == nil
}

func addPoints(x1, y1, x2, y2, a, p *big.Int) (*big.Int, *big.Int, error) {
	if isInfinity(x1, y1) {
		return x2, y2, nil
	}

	if isInfinity(x2, y2) {
		return x1, y1, nil
	}

	if x1.Cmp(x2) == 0 && y1.Cmp(y2) == 0 {
		return doublePoint(x1, y1, a, p)
	}

	lambda := new(big.Int).Sub(y2, y1)
	xDiff := new(big.Int).Sub(x2, x1)
	xDiff.ModInverse(xDiff, p)
	lambda.Mul(lambda, xDiff).Mod(lambda, p)

	x3 := new(big.Int).Mul(lambda, lambda)
	x3.Sub(x3, x1).Sub(x3, x2).Mod(x3, p)

	y3 := new(big.Int).Sub(x1, x3)
	y3.Mul(lambda, y3).Sub(y3, y1).Mod(y3, p)

	return x3, y3, nil
}

func doublePoint(x, y, a, p *big.Int) (*big.Int, *big.Int, error) {
	if y.Cmp(big.NewInt(0)) == 0 {
		return nil, nil, errors.New("удвоение точки на бесконечности")
	}

	lambda := new(big.Int).Mul(x, x)
	lambda.Mul(lambda, big.NewInt(3)).Add(lambda, a)
	twoY := new(big.Int).Mul(y, big.NewInt(2))

	twoY.ModInverse(twoY, p)

	lambda.Mul(lambda, twoY).Mod(lambda, p)

	x3 := new(big.Int).Mul(lambda, lambda)
	x3.Sub(x3, new(big.Int).Mul(x, big.NewInt(2))).Mod(x3, p)

	y3 := new(big.Int).Sub(x, x3)
	y3.Mul(lambda, y3).Sub(y3, y).Mod(y3, p)

	return x3, y3, nil
}

func scalarMultiply(x, y, k, a, p *big.Int) (*big.Int, *big.Int) {
	if k.Cmp(big.NewInt(0)) == 0 {
		return big.NewInt(0), big.NewInt(0)
	}
	rx, ry := (*big.Int)(nil), (*big.Int)(nil)

	for i := k.BitLen() - 1; i >= 0; i-- {
		if !isInfinity(rx, ry) {
			rx, ry, _ = doublePoint(rx, ry, a, p)
		}

		if k.Bit(i) == 1 {
			if isInfinity(rx, ry) {
				rx, ry = x, y
			} else {
				rx, ry, _ = addPoints(rx, ry, x, y, a, p)
			}
		}
	}

	return rx, ry
}

type Pair struct {
	X *big.Int
	Y *big.Int
	k big.Int
}

func alike(tuple_k, tuple_r_plus_p []Pair, k int) (bool, [][]*big.Int) {
	var m [][]*big.Int
	flag := false
	for _, el := range tuple_k {
		for _, elp := range tuple_r_plus_p {
			if el.X.Cmp(elp.X) == 0 {
				fmt.Println("\nЕсть совпадениe:")
				flag = true
				fmt.Println(&el.k, "Q = ", el.X, " ", el.Y)
				fmt.Println("R +", &elp.k, "p = ", elp.X, " ", elp.Y)
				fmt.Println(el.X, "=", elp.X)
				fmt.Println("e =", &el.k)
				e_prob := el.k
				k_big := big.NewInt(int64(k))
				two_i_k := new(big.Int).Mul(new(big.Int).Mul(&elp.k, big.NewInt(2)), k_big)
				d_nt_true := new(big.Int).Div(new(big.Int).Abs(new(big.Int).Add(two_i_k, &elp.k)), k_big)
				big.NewInt(int64(k))
				fmt.Println("d =", d_nt_true)
				pair := []*big.Int{d_nt_true, &e_prob}
				m = append(m, pair)
			}
		}
	}
	return flag, m
}

func genPoint_two(qx, qy, r int, a, p *big.Int) (*big.Int, *big.Int) {
	big_qx, _ := new(big.Int).SetString(strconv.Itoa(qx), 10)
	big_qy, _ := new(big.Int).SetString(strconv.Itoa(qy), 10)
	k, _ := new(big.Int).SetString(strconv.Itoa(rand.Intn(r/2)), 10)
	x, y := scalarMultiply(big_qx, big_qy, k, a, p)
	return x, y
}

func main() {

	var a int
	var i int
	p := 569
	a = 367

	x, y := genPointOnCurve(a, p)

	big_a, _ := new(big.Int).SetString(strconv.Itoa(a), 10)
	big_p, _ := new(big.Int).SetString(strconv.Itoa(p), 10)

	big_x, _ := new(big.Int).SetString(strconv.Itoa(x), 10)
	big_y, _ := new(big.Int).SetString(strconv.Itoa(y), 10)

	if x == -1 && y == -1 {
		fmt.Println("не выходит подобрать точку")
		return
	}

	d_p := qudro_sqrt(p)

	fmt.Println("Сгенерированная точка ", big_x, big_y)
	fmt.Println("k = ", d_p)

	var tuple_k []Pair
	tuple_k = append(tuple_k, Pair{big_x, big_y, *big.NewInt(0)})
	for i = 0; i < d_p; i++ {
		if i+1 == 0 {
			continue
		}
		big_i, _ := new(big.Int).SetString(strconv.Itoa(i+1), 10)
		kx, xy := scalarMultiply(big_x, big_y, big_i, big_a, big_p)
		tuple_k = append(tuple_k, Pair{kx, xy, *big.NewInt(int64(i + 1))})
		tuple_k = append(tuple_k, Pair{kx, NegatePoint(xy, big_p), *big.NewInt(int64(-(i + 1)))})
	}

	kq := tuple_k[len(tuple_k)-2]
	fmt.Println(kq.X, kq.Y, kq.k)

	big_quadro, _ := new(big.Int).SetString(strconv.Itoa(d_p), 10)

	sort.Slice(tuple_k, func(i, j int) bool {
		return tuple_k[i].X.Cmp(tuple_k[j].X) < 0
	})

	p_caps := (new(big.Int).Add(new(big.Int).Mul(big_quadro, big.NewInt(2)), big.NewInt(1)))
	r_caps := (new(big.Int).Add(big_p, big.NewInt(1)))
	p_x, p_y := scalarMultiply(big_x, big_y, p_caps, big_a, big_p)
	r_x, r_y := scalarMultiply(big_x, big_y, r_caps, big_a, big_p)

	var tuple_r_plus_p []Pair
	tuple_r_plus_p = append(tuple_r_plus_p, Pair{r_x, r_y, *big.NewInt(0)})

	for i = 1; i <= d_p; i++ {
		kp_x, kp_y := scalarMultiply(p_x, p_y, big.NewInt(int64(i)), big_a, big_p)
		sum_x, sum_y, _ := addPoints(r_x, r_y, kp_x, kp_y, big_a, big_p)
		sub_x, sub_y, _ := addPoints(r_x, r_y, kp_x, NegatePoint(kp_y, big_p), big_a, big_p)
		tuple_r_plus_p = append(tuple_r_plus_p, Pair{sum_x, sum_y, *big.NewInt(int64(i))})
		tuple_r_plus_p = append(tuple_r_plus_p, Pair{sub_x, sub_y, *big.NewInt(int64(-i))})
	}

	for _, el := range tuple_k {
		fmt.Println(&el.k, "Q = ", el.X, " ", el.Y)
	}

	fmt.Println("\n")
	for _, el := range tuple_r_plus_p {
		fmt.Println("R+", &el.k, "p = ", el.X, " ", el.Y)
	}

	fl, aboba := alike(tuple_k, tuple_r_plus_p, d_p)

	if !fl {
		return
	}

	for _, el := range aboba {
		d, e := el[0], el[1]
		n_big := f_n(d, e, big_p, d_p)
		fmt.Println("При коэффициентах d, e = ", d, " ", e)
		fmt.Println("Порядок равен = ", n_big, "\n")
	}

}
func f_n(d, e, big_p *big.Int, d_p int) *big.Int {
	p_plus_one := new(big.Int).Add(big_p, big.NewInt(1))
	dk_sub_e := new(big.Int).Sub(new(big.Int).Mul(d, big.NewInt(int64(d_p))), e)
	return new(big.Int).Add(p_plus_one, dk_sub_e)
}

func gen_coff(d_p int, x, y *big.Int, kq Pair, big_a, big_p, r_x, r_y *big.Int) []*big.Int {
	fmt.Println("\n")
	fmt.Println("R ", r_x, r_y)
	fmt.Println("Q ", x, y)
	fmt.Println("kQ ", kq.X, kq.Y)
	var list []*big.Int
	return list
}

func NegatePoint(y, p *big.Int) (negY *big.Int) {
	negY = new(big.Int).Neg(y)
	negY.Mod(negY, p)

	return negY
}
