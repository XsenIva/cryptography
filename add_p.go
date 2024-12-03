package main

import (
	"errors"
	"math/big"
)

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

func scalarMultiply(x, y, k, a, p *big.Int) (*big.Int, *big.Int, bool) {
	rx, ry := (*big.Int)(nil), (*big.Int)(nil)
	var flag bool = false
	for i := k.BitLen() - 1; i >= 0; i-- {
		if !isInfinity(rx, ry) {
			rx, ry, _ = doublePoint(rx, ry, a, p)
		}

		if k.Bit(i) == 1 {
			if isInfinity(rx, ry) {
				// file.WriteString("\n" + "Точка уходит в бесконечность")
				flag = true
				rx, ry = x, y
			} else {
				rx, ry, _ = addPoints(rx, ry, x, y, a, p)
				// file.WriteString("\n" + "Точка достигла " + rx.String() + ", " + ry.String())
			}
		}
	}

	return rx, ry, flag
}
