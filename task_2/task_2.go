// package main

// import (
// 	"fmt"
// 	"math/big"
// )

// type Curve struct {
// 	A *big.Int // Коэффициент a
// 	B *big.Int // Коэффициент b
// 	P *big.Int // Простое число (модуль)
// }

// // Проверяет, является ли число квадратичным вычетом по модулю p
// func isQuadraticResidue(value, p *big.Int) bool {
// 	// Проверяем (value ^ ((p-1)/2)) % p == 1
// 	exp := new(big.Int).Sub(p, big.NewInt(1))
// 	exp.Div(exp, big.NewInt(2))
// 	result := new(big.Int).Exp(value, exp, p)
// 	return result.Cmp(big.NewInt(1)) == 0
// }

// // Подсчитывает количество точек на кривой y^2 = x^3 + ax + b (mod p)
// func countPoints(curve Curve) int {
// 	count := 1 // Учитываем точку на бесконечности

// 	for x := big.NewInt(0); x.Cmp(curve.P) < 0; x.Add(x, big.NewInt(1)) {
// 		// Вычисляем правую часть уравнения: x^3 + ax + b (mod p)
// 		xCubed := new(big.Int).Exp(x, big.NewInt(3), curve.P)
// 		ax := new(big.Int).Mul(curve.A, x)
// 		rightSide := new(big.Int).Add(xCubed, ax)
// 		rightSide.Add(rightSide, curve.B)
// 		rightSide.Mod(rightSide, curve.P)

// 		// Проверяем, является ли правое значение квадратичным вычетом
// 		if isQuadraticResidue(rightSide, curve.P) {
// 			count += 2 // Если квадратичный вычет, то 2 точки (y и -y)
// 		}
// 	}

// 	return count
// }

// // Основная функция
// func main() {
// 	// Пример: кривая y^2 = x^3 + 2x + 3 (mod 97)
// 	curve := Curve{
// 		A: big.NewInt(6),
// 		B: big.NewInt(3),
// 		P: big.NewInt(23),
// 	}
// //  28
// 	// Подсчитываем количество точек
// 	numPoints := countPoints(curve)
// 	fmt.Printf("Число точек на кривой: %d\n", numPoints+1)
// }

package main

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"sort"
	"strconv"
)

func qudro_sqrt(p int) int {
	p_float := float64(p * 2)
	k := int(math.Ceil(math.Sqrt(math.Sqrt(p_float))))
	return k
}

func gen_points(a, b int) (int, int) {
	var x, y int
	for y = 0; y <= 1000; y++ {
		for x = 0; x <= 1000; x++ {
			if (y * y) == (x*x*x + a*x + b) {
				return x, y
			}
		}
	}
	return 0, 0
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
	// fmt.Println("add", lambda)

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

func main() {
	var a, b int
	var i int
	a, b = 5, 3
	p := 23
	x, y := gen_points(a, b)
	if x == 0 && y == 0 {
		fmt.Println("не выходит подобрать точку")
		return
	}
	quadro_p := qudro_sqrt(p)
	fmt.Println(x, y)
	fmt.Println(quadro_p)
	big_x, _ := new(big.Int).SetString(strconv.Itoa(x), 10)
	big_y, _ := new(big.Int).SetString(strconv.Itoa(y), 10)
	big_a, _ := new(big.Int).SetString(strconv.Itoa(a), 10)
	big_p, _ := new(big.Int).SetString(strconv.Itoa(p), 10)

	var tuple_k []Pair
	tuple_k = append(tuple_k, Pair{big_x, big_y, *big.NewInt(0)})
	// map_k[1] = Pair{big_x, big_y, big.NewInt(1)}
	for i = 1; i < quadro_p; i++ {
		big_i, _ := new(big.Int).SetString(strconv.Itoa(i+1), 10)
		kx, xy := scalarMultiply(big_x, big_y, big_i, big_a, big_p)
		fmt.Println(i+1, kx, xy)
		tuple_k = append(tuple_k, Pair{kx, xy, *big.NewInt(int64(i + 1))})
	}
	big_quadro, _ := new(big.Int).SetString(strconv.Itoa(quadro_p), 10)

	sort.Slice(tuple_k, func(i, j int) bool {
		return tuple_k[i].X.Cmp(tuple_k[j].X) < 0
	})

	fmt.Println(tuple_k)
	p_caps := (new(big.Int).Add(new(big.Int).Mul(big_quadro, big.NewInt(2)), big.NewInt(1)))
	r_caps := (new(big.Int).Add(big_p, big.NewInt(1)))
	p_x, p_y := scalarMultiply(big_x, big_y, p_caps, big_a, big_p)
	r_x, r_y := scalarMultiply(big_x, big_y, r_caps, big_a, big_p)
	r_p_x, r_p_y, _ := addPoints(r_x, r_y, p_x, p_y, big_a, big_p)

	var tuple_r_plus_p []Pair
	tuple_r_plus_p = append(tuple_r_plus_p, Pair{r_x, r_y, *big.NewInt(0)})
	tuple_r_plus_p = append(tuple_r_plus_p, Pair{r_p_x, r_p_y, *big.NewInt(1)})
	for i = 2; i <= quadro_p; i++ {
		kp_x, kp_y := scalarMultiply(p_x, p_y, big.NewInt(int64(i)), big_a, big_p)
		sum_x, sum_y, _ := addPoints(r_x, r_y, kp_x, kp_y, big_a, big_p)
		sub_x, sub_y, _ := addPoints(r_x, r_y, kp_x, kp_y.Neg(kp_y), big_a, big_p)
		tuple_r_plus_p = append(tuple_r_plus_p, Pair{sum_x, sum_y, *big.NewInt(int64(i))})
		tuple_r_plus_p = append(tuple_r_plus_p, Pair{sub_x, sub_y, *big.NewInt(int64(i))})
	}

	sort.Slice(tuple_r_plus_p, func(i, j int) bool {
		return tuple_r_plus_p[i].X.Cmp(tuple_r_plus_p[j].X) < 0
	})
	fmt.Println(tuple_r_plus_p)
	for _, el := range tuple_k {
		for _, elp := range tuple_r_plus_p {
			if (el.X == elp.X) && (el.Y == elp.Y) {

			}
		}
	}
}
