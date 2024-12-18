import sympy
import random
import math
from collections import namedtuple

primes = set()


def add_ellip_points(x, y, x1, y1, A, p):
    if x == None and y == None:
        return x1, y1

    if x1 == None and y1 == None:
        return x, y
    x = x % p
    x1 = x1 % p
    y = y % p
    y1 = y1 % p
    alpha = None
    if x != x1 or y != y1:
        if x == x1:
            return None, None
        alpha = ((y1 - y) * sympy.mod_inverse(x1 - x, p)) % p
    else:
        if y == 0:
            return None, None
        alpha = ((3 * x * x + A) * sympy.mod_inverse(2 * y, p)) % p

    x3 = (alpha ** 2 - x - x1) % p
    y3 = (alpha * (x - x3) - y) % p

    return (x3, y3)


def scalarMultiply(other, x, y, p, A):
    if other == 0:
        return (None, None)
    negate = False
    if other < 0:
        negate = True
        other = -other
    init = False
    powerOf2 = (x, y)
    while True:
        if other % 2 == 1:
            if init:
                res = add_ellip_points(*res, *powerOf2, A, p)
            else:
                res = powerOf2
                init = True
        other //=2
        if other == 0:
            break
        powerOf2 = add_ellip_points(*powerOf2, *powerOf2, A, p)
    if negate:
        res = (res[0], -res[1] % p)
    return res

def qudro_sqrt(p: int) -> int:
    p_float = float(p * 2)
    k = int(math.ceil(math.sqrt(math.sqrt(p_float))))
    return k


def gen_point(a, p):
  
  x = random.randint(1, int(p))
  y = random.randint(1, int(p))
  while (y*y) % p != (x*x*x + a * x) % p:
    x = random.randint(1, int(p))
    y = random.randint(1, int(p))
  return x, y
   

def negate_point(y, p):
  if y != None:
    return (-y) % p


Pair = namedtuple("Pair", ["x", "y", "k"])


def gen_base_q(qx, qy, r, a, p):
    tuple_k = []
    tuple_k.append(Pair(qx, qy, 0))

    for i in range(r):
        if i + 1 == 0:
            continue
        big_i = i + 1
        kx, ky = scalarMultiply(big_i,qx, qy, p, a)
        tuple_k.append(Pair(kx, ky, big_i))
        tuple_k.append(Pair(kx, negate_point(ky, p), -big_i))
    
    return tuple_k


def gen_base_r_plus_p(p_big, r_big, k, a, p):
    tuple_r_p = []
    tuple_r_p.append(Pair(r_big[0],r_big[1], 0))
    for i in range(1, k+1):
        kpx, kpy = scalarMultiply(i,  p_big[0],  p_big[1], p, a)
        sumx, sumy = add_ellip_points(r_big[0],r_big[1], kpx, kpy, a,p)
        tuple_r_p.append(Pair(sumx, sumy, i))
        sumx, sumy = add_ellip_points(r_big[0],r_big[1], kpx, negate_point(kpy, p), a, p)
        tuple_r_p.append(Pair(sumx, sumy, - i))
    return tuple_r_p


def cheсk_comparison(res_1, res_2):
  maybe_ed = []
  flag = False
  for r_add_p in res_2:
      for kq in res_1:
          if r_add_p.x == kq.x:
              maybe_ed.append([r_add_p.k, kq.k])
              flag = True
  return flag, maybe_ed


def findinfg( k, r_big, x, y, kq, a, p):
   print("\n")
   kt = k*2
   res = []
   for d in range(-kt, kt):
      mult = scalarMultiply(d, kq[0], kq[1], p, a)
      summ = add_ellip_points(mult[0], mult[1], r_big[0], r_big[1], a, p)
      for e in range(-kt, kt):
        mult_e_q = scalarMultiply(e, x, y, p, a)
        # print(d, e, summ, mult_e_q)
        if summ[0] == mult_e_q[0] and summ[1]==mult_e_q[1]:
           if abs(p+1+d*k-e - p+1) <= 2* math.sqrt(p):
                # print("N =", p+1+d*k-e)
                print("d, e =", d, e)
                res.append(p+1+d*k-e)
   return res

def main():
  p = int(input("p = "))
  a = int(input("A = "))
  
  x, y  = gen_point(a, p)

  print("\nСгенерирована точка ", x, y, "\n" )

  k = qudro_sqrt(p)

  flag = False

  while not flag:
    res_1 = gen_base_q(x, y, k, a, p)
    kq = res_1[-2]
    res_1 = sorted(gen_base_q(x, y, k, a, p), key=lambda pair: getattr(pair, 'x'))

    for el in res_1:
        print(el.k,"Q =", el.x, el.y, "\n")
   
    p_big = scalarMultiply(2*k+1 ,x, y, p, a)
    r_big = scalarMultiply(p + 1 ,x, y, p, a)

    print("\nСуммы R + iP ")
    res_2 = gen_base_r_plus_p(p_big, r_big, k, a, p)
    for el in res_2:
        print("R +", el.k,"P =", el.x, el.y, "\n")
    
    flag, mayb_de = cheсk_comparison(res_1, res_2)
    if not flag: 
      return
    else:
      print("Совпадения есть")
      for el in mayb_de:
        print("R +", el[0]," P = ", el[1], "Q")

    res = findinfg(k, r_big, x, y, kq, a, p)
    print("N =",  res[0])
    # print("N =",  res)
main()

